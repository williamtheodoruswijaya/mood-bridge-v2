package main

import (
	"mood-bridge-v2/server/infrastructure/db"
	"mood-bridge-v2/server/internal/api"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Pertama-tama kita akan inisialisasi koneksi ke database PostgreSQL
	database := db.NewDbConnection()
	defer database.Close()

	// Setelah itu kita akan migrasi database-nya
	db.Migrate(database, "up")

	// Setup redis client untuk caching
	redisUrl := os.Getenv("REDIS_URL")
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic("Failed to parse Redis URL: " + err.Error())
	}
	rdb := redis.NewClient(opt)

	defer rdb.Close()

	// Terakhir kita jalankan server-nya
	router := api.SetupRoutes(database, rdb)
	router.Run(":8080")
}
