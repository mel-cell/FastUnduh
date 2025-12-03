package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// Kita definisikan ulang strukturnya biar worker mandiri
type JobData struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Status   string `json:"status"`
	Filename string `json:"filename"` // Nanti diisi worker
	Title    string `json:"title"`    // Nanti diisi worker
}

// Struct untuk menangkap output JSON dari yt-dlp
type YtDlpOutput struct {
	Title    string `json:"title"`
	Filename string `json:"_filename"` // yt-dlp kadang kasih nama file asli
}

func StartWorker(workerID int, rdb *redis.Client) {
	fmt.Printf("👷 Worker-%d: Siap menunggu perintah...\n", workerID)

	for {
		// 1. BLPOP: Blocking Pop (Tunggu sampai ada data di queue)
		// Ini efisien! Kalau antrian kosong, Worker TIDUR (CPU 0%).
		// Tidak perlu "Polling" tiap detik.
		result, err := rdb.BLPop(ctx, 0*time.Second, "queue:jobs").Result()
		if err != nil {
			fmt.Println("Error antrian:", err)
			continue
		}

		// result[0] = nama key ("queue:jobs")
		// result[1] = value (Job ID yang kita butuhkan)
		jobID := result[1]

		// 2. Ambil Detail Job dari "Papan Pengumuman" Redis
		val, err := rdb.Get(ctx, "job:"+jobID).Result()
		if err != nil {
			fmt.Printf("Worker-%d: Job %s hilang dari radar\n", workerID, jobID)
			continue
		}

		var job JobData
		json.Unmarshal([]byte(val), &job)

		// 3. Update Status -> PROCESSING
		job.Status = "processing"
		updateRedis(rdb, job)
		// fmt.Printf("Worker-%d: Memasak Job %s...\n", workerID, jobID)

		// 4. EKSEKUSI YT-DLP (Inti dari segalanya)
		// Kita paksa output filenya jadi: ./downloads/JOB_ID.mp4
		// Biar nanti API gampang nyarinya.
		outputTarget := fmt.Sprintf("./downloads/%s.mp4", jobID)

		cmd := exec.Command("yt-dlp",
			"--print-json", // Minta output JSON biar kita bisa ambil Judul
			"--no-warnings",
			"--format", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best", // Paksa MP4
			"-o", outputTarget, // Rename file jadi ID unik
			job.URL,
		)

		// Jalankan perintah dan tangkap outputnya
		outputBytes, err := cmd.CombinedOutput()

		if err != nil {
			// GAGAL :(
			fmt.Printf("Worker-%d: Gagal download %s. Error: %v\n", workerID, job.URL, err)
			job.Status = "failed"
			updateRedis(rdb, job)
			continue
		}

		// 5. Parsing Output yt-dlp untuk ambil Judul Asli
		var ytdlpMeta YtDlpOutput
		// yt-dlp print json di baris terakhir, kadang kecampur log lain
		// Tapi biasanya unmarshal cukup pintar cari json valid
		json.Unmarshal(outputBytes, &ytdlpMeta)

		// 6. SUKSES! Update Status -> COMPLETED
		job.Status = "completed"
		job.Filename = jobID + ".mp4" // Nama file di server
		job.Title = ytdlpMeta.Title
		if job.Title == "" {
			job.Title = "Video Downloaded" // Fallback kalau gagal parsing title
		}

		updateRedis(rdb, job)
		fmt.Printf("✅ Worker-%d: Selesai! %s -> %s\n", workerID, job.URL, job.Filename)
	}
}

// Fungsi bantu update Redis biar codingan rapi
func updateRedis(rdb *redis.Client, job JobData) {
	jobJSON, _ := json.Marshal(job)
	// Keep TTL (tetap expire dalam sisa waktu, atau reset ke 2 jam)
	rdb.Set(ctx, "job:"+job.ID, jobJSON, 2*time.Hour)
}
