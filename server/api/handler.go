package api

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// Struktur Data Request dari User
type RequestData struct {
	URL string `json:"url"`
}

// Struktur Data Job yang akan disimpan ke Redis
type JobData struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Status    string `json:"status"` // pending, processing, completed, failed
	CreatedAt string `json:"created_at"`
	// Field hasil download (nanti diisi worker)
	Filename string `json:"filename,omitempty"`
	Title    string `json:"title,omitempty"`
}

func StartServer(rdb *redis.Client) {
	app := fiber.New()

	// Setup CORS agar Vue bisa akses
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// --- ROUTE 1: TERIMA ORDER (POST) ---
	app.Post("/api/queue", func(c *fiber.Ctx) error {
		// 1. Parsing Body JSON
		var req RequestData
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "JSON tidak valid"})
		}
		if req.URL == "" {
			return c.Status(400).JSON(fiber.Map{"error": "URL tidak boleh kosong"})
		}

		// 2. Buat Tiket Antrian Baru
		jobID := uuid.New().String()
		newJob := JobData{
			ID:        jobID,
			URL:       req.URL,
			Status:    "pending",
			CreatedAt: time.Now().Format(time.RFC3339),
		}

		// 3. Simpan Data Job ke Redis (Hash Map) -> Key: "job:{id}"
		// Kita simpan sebagai JSON string biar gampang
		jobJSON, _ := json.Marshal(newJob)
		err := rdb.Set(ctx, "job:"+jobID, jobJSON, 2*time.Hour).Err() // Expire 2 jam
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal simpan ke Redis"})
		}

		// 4. Masukkan ID ke dalam Antrian (List) -> Key: "queue:jobs"
		// Worker nanti akan 'mengintip' key ini
		err = rdb.LPush(ctx, "queue:jobs", jobID).Err()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal masuk antrian"})
		}

		// 5. Respon Cepat ke User
		return c.JSON(fiber.Map{
			"message": "Berhasil masuk antrian",
			"job_id":  jobID,
			"status":  "pending",
		})
	})

	// --- ROUTE 2: CEK STATUS (GET) ---
	// Dipakai Frontend untuk Polling (Nanya terus-menerus)
	app.Get("/api/status/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Ambil data dari Redis
		val, err := rdb.Get(ctx, "job:"+id).Result()
		if err == redis.Nil {
			return c.Status(404).JSON(fiber.Map{"status": "not_found"})
		} else if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Redis error"})
		}

		// Unmarshal string JSON dari Redis kembali ke Struct
		var job JobData
		json.Unmarshal([]byte(val), &job)

		return c.JSON(job)
	})

	// --- ROUTE 3: DOWNLOAD FILE (GET) ---
	// Nanti user diarahkan ke sini kalau status sudah 'completed'
	app.Get("/api/download/:filename", func(c *fiber.Ctx) error {
		filename := c.Params("filename")
		// Pastikan path-nya benar sesuai volume Docker nanti
		path := "./downloads/" + filename

		c.Set("Content-Disposition", "attachment; filename="+filename)
		return c.SendFile(path)
	})

	// Jalankan Server
	app.Listen(":3000")
}
