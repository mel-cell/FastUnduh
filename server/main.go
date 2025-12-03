package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// --- Configuration ---
const (
	MaxWorkers     = 5               // Concurrency Control: Cuma 5 kuli yang kerja barengan
	QueueKey       = "queue:downloads"
	JobPrefix      = "job:"
	ResultPrefix   = "result:"
	FileTTL        = 15 * time.Minute // Janitor Rule
)

var ctx = context.Background()
var rdb *redis.Client


type DownloadRequest struct {
	URL string `json:"url"`
}

// Data Status Job (Untuk tracking)
type JobData struct {
	ID        string `redis:"id"`
	URL       string `redis:"url"`
	Status    string `redis:"status"` // pending, processing, completed, failed
	Message   string `redis:"message"`
	CreatedAt string `redis:"created_at"`
}

// Data Result (Hasil download)
type ResultData struct {
	Filename    string `redis:"filename"`
	Filepath    string `redis:"filepath"`
	ContentType string `redis:"content_type"`
}

func main() {
	// 1. Initialize Redis (The Queue & Storage)
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Pastikan Redis jalan di port ini
		DB:   0,
	})

	// 2. Start The Workers (Kuli Panggul)
	// Kita nyalakan 5 worker yang akan standby nungguin tugas dari Redis
	for i := 1; i <= MaxWorkers; i++ {
		go startWorker(i)
	}

	// 3. Start The Janitor (Satpam Kebersihan)
	go startJanitor()

	// 4. Start The Receiver (API Gateway)
	app := fiber.New()
	app.Use(cors.New())

	// Endpoint: Menerima Pesanan
	app.Post("/download", handleDownloadRequest)
	
	// Endpoint: Cek Status Pesanan
	app.Get("/status/:jobId", handleCheckStatus)

	log.Println("🚀 FastUnduh Backend System Started on :3000")
	log.Fatal(app.Listen(":3000"))
}

// --- The Receiver (API Gateway) ---
func handleDownloadRequest(c *fiber.Ctx) error {
	var req DownloadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validasi URL sederhana
	if req.URL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "URL cannot be empty"})
	}

	// Generate Tiket Antrian
	jobID := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	// Simpan Info Tiket ke Redis (Status: PENDING)
	err := rdb.HSet(ctx, JobPrefix+jobID, map[string]interface{}{
		"id":         jobID,
		"url":        req.URL,
		"status":     "pending",
		"message":    "Menunggu antrian...",
		"created_at": now,
	}).Err()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data ke Redis"})
	}

	// Masukkan Tiket ke Antrian (Push ke Redis List)
	// Worker akan mengambil dari sisi kiri (LPUSH -> BRPOP atau RPUSH -> BLPOP)
	// Kita pakai RPUSH (masuk belakang), Worker BLPOP (ambil depan)
	err = rdb.RPush(ctx, QueueKey, jobID).Err()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal masuk antrian"})
	}

	// Respon Cepat ke User (Deferred)
	return c.JSON(fiber.Map{
		"job_id":  jobID,
		"message": "Pesanan diterima. Mohon tunggu.",
		"status":  "pending",
	})
}

func handleCheckStatus(c *fiber.Ctx) error {
	jobID := c.Params("jobId")

	// Ambil data job
	jobData, err := rdb.HGetAll(ctx, JobPrefix+jobID).Result()
	if err != nil || len(jobData) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tiket tidak ditemukan atau sudah kadaluarsa"})
	}

	// Jika completed, ambil juga data result-nya
	response := fiber.Map{
		"job": jobData,
	}

	if jobData["status"] == "completed" {
		resultData, _ := rdb.HGetAll(ctx, ResultPrefix+jobID).Result()
		response["result"] = resultData
	}

	return c.JSON(response)
}

// --- The Worker (Engine) ---
func startWorker(id int) {
	log.Printf("👷 Worker #%d siap bekerja!", id)

	for {
		// BLPOP: Blocking Pop. Worker akan "tidur" kalau antrian kosong.
		// Begitu ada data, dia langsung bangun. (Timeout 0 = tunggu selamanya)
		result, err := rdb.BLPop(ctx, 0*time.Second, QueueKey).Result()
		if err != nil {
			log.Printf("Worker #%d error Redis: %v", id, err)
			time.Sleep(1 * time.Second) // Istirahat bentar kalau error
			continue
		}

		// result[0] adalah nama key (queue:downloads), result[1] adalah value (jobID)
		jobID := result[1]
		log.Printf("👷 Worker #%d memproses tiket: %s", id, jobID)

		processJob(jobID)
	}
}

func processJob(jobID string) {
	// 1. Update Status -> PROCESSING
	rdb.HSet(ctx, JobPrefix+jobID, "status", "processing", "message", "Sedang mendownload...")

	// Ambil URL
	url, _ := rdb.HGet(ctx, JobPrefix+jobID, "url").Result()

	// 2. Panggil Python (yt-dlp)
	// Output template: ./downloads/<jobID>.%(ext)s
	outputTemplate := fmt.Sprintf("./downloads/%s.%%(ext)s", jobID)
	
	// Command: python3 downloader.py <url> <output_template>
	cmd := exec.Command("python3", "downloader.py", url, outputTemplate)
	
	// Kita tangkap outputnya kalau mau debug, tapi yang penting exit code-nya
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("❌ Job %s Gagal: %v\nOutput: %s", jobID, err, string(output))
		rdb.HSet(ctx, JobPrefix+jobID, "status", "failed", "message", "Gagal mendownload video. Cek linknya.")
		return
	}

	// 3. Sukses! Cari nama filenya.
	// Karena yt-dlp bisa save jadi .mp4, .mkv, atau .webm, kita harus cari file yang depannya jobID
	finalFilename := findDownloadedFile(jobID)
	if finalFilename == "" {
		rdb.HSet(ctx, JobPrefix+jobID, "status", "failed", "message", "File tidak ditemukan setelah download.")
		return
	}

	// 4. Simpan Result & Update Status -> COMPLETED
	rdb.HSet(ctx, ResultPrefix+jobID, map[string]interface{}{
		"filename":    finalFilename,
		"filepath":    "./downloads/" + finalFilename,
		"content_type": "video/mp4", // Simplifikasi, aslinya bisa detect mime type
	})
	
	// Set TTL untuk Result (Data hilang dari Redis setelah 15 menit)
	rdb.Expire(ctx, ResultPrefix+jobID, FileTTL)
	rdb.Expire(ctx, JobPrefix+jobID, FileTTL) // Job info juga dihapus

	rdb.HSet(ctx, JobPrefix+jobID, "status", "completed", "message", "Download selesai!")
	log.Printf("✅ Job %s Selesai! File: %s", jobID, finalFilename)
}

// Helper untuk mencari file hasil download
func findDownloadedFile(jobID string) string {
	files, err := os.ReadDir("./downloads")
	if err != nil {
		return ""
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), jobID) {
			return file.Name()
		}
	}
	return ""
}

// --- The Janitor (Pembersih) ---
func startJanitor() {
	log.Println("🧹 Janitor siap bersih-bersih setiap 1 menit.")
	ticker := time.NewTicker(1 * time.Minute)
	
	for range ticker.C {
		files, err := os.ReadDir("./downloads")
		if err != nil {
			continue
		}

		for _, file := range files {
			info, err := file.Info()
			if err != nil {
				continue
			}

			// Hapus file fisik jika umurnya > 15 menit
			if time.Since(info.ModTime()) > FileTTL {
				os.Remove("./downloads/" + file.Name())
				log.Printf("🧹 Janitor membakar sampah: %s", file.Name())
			}
		}
	}
}
