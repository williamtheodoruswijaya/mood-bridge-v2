package main

import (
	"mood-bridge-v2/server/infrastructure/cache"
	"mood-bridge-v2/server/infrastructure/db"
	"mood-bridge-v2/server/internal/api"
)

func main() {
	// Pertama-tama kita akan inisialisasi koneksi ke database PostgreSQL
	database := db.NewDbConnection()
	defer database.Close()

	// Setelah itu kita akan migrasi database-nya
	db.Migrate(database, "up")

	// Setup redis client untuk caching
	rdb := cache.NewRedisClient()
	defer rdb.Close()

	// Terakhir kita jalankan server-nya
	router := api.SetupRoutes(database, rdb)
	router.Run(":8080")
}
