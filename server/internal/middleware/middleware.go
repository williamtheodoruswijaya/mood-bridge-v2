package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware buat menangani panic (error yang tidak terduga) di server
// Ini akan mengembalikan response JSON dengan status 500 Internal Server Error
func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		// rollback function yang akan dieksekusi
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Panic occurred:", r)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "error",
					"message": "Internal server error",
				})
				c.Abort()
			}
		}()

		// lanjutkan ke handler berikutnya
		// Jika ada panic di handler berikutnya, maka defer function di atas akan dieksekusi
		c.Next()
	}
}

// Intinya biar localhost:8080 bisa diakses sama localhost:3000
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}