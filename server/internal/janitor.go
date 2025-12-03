package internal

import (
	"fastunduh-backend/config"
	"log"
	"os"
	"time"
)

func StartJanitor() {
	log.Println("🧹 Janitor: Siap bersih-bersih folder /downloads setiap 1 menit...")
	
	// Cek setiap 1 menit
	ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		files, err := os.ReadDir("./downloads")
		if err != nil {
			log.Println("Janitor Error: Gagal baca folder downloads")
			continue
		}

		for _, file := range files {
			info, err := file.Info()
			if err != nil {
				continue
			}

			// Rumus: Umur File = Waktu Sekarang - Waktu Modifikasi
			umurFile := time.Since(info.ModTime())

			if umurFile > config.FileTTL {
				err := os.Remove("./downloads/" + file.Name())
				if err != nil {
					log.Printf("❌ Janitor: Gagal hapus %s\n", file.Name())
				} else {
					log.Printf("🔥 Janitor: Membakar sampah %s (Umur: %.0f menit)\n", file.Name(), umurFile.Minutes())
				}
			}
		}
	}
}
