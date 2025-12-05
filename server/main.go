package main

import (
	"fastunduh-backend/api"
	"fastunduh-backend/config"
	"fastunduh-backend/internal"
	"fastunduh-backend/worker"
	"log"
	"os"
)

func main() {
	// 0. Tambahkan CWD ke PATH agar bisa nemu ./yt-dlp
	cwd, _ := os.Getwd()
	os.Setenv("PATH", os.Getenv("PATH") + ":" + cwd)

	// 1. Buat folder downloads jika belum ada
	if err := os.MkdirAll("./downloads", 0755); err != nil {
		log.Fatal("Gagal membuat folder downloads:", err)
	}

	// 1. Init Config & Redis
	config.InitRedis()
	
	// Pastikan koneksi Redis berhasil
	if _, err := config.Rdb.Ping(config.Ctx).Result(); err != nil {
		log.Fatalf("Gagal connect ke Redis: %v. Pastikan Redis sudah jalan!", err)
	}

	// 2. Start Workers
	// Worker akan standby mendengarkan antrian dari Redis
	for i := 1; i <= config.MaxWorkers; i++ {
		go worker.StartWorker(i, config.Rdb)
	}

	// 3. Start Janitor
	// Membersihkan file lama secara berkala
	go internal.StartJanitor()

	// 4. Start API Server
	// Blocking call (akan menahan main process agar tidak exit)
	log.Println("🚀 FastUnduh Backend System Starting on :3000")
	api.StartServer(config.Rdb)
}
