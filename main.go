package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/webook/internal/repository"
	"github.com/webook/internal/repository/dao"
	"github.com/webook/internal/service"
	"github.com/webook/internal/web"
	"github.com/webook/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	var db *gorm.DB
	var err error
	maxRetries := 10
	retryInterval := 5 * time.Second
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(retryInterval)
		if i == maxRetries-1 {
			panic(err)
		}
	}

	if err := dao.InitTables(db); err != nil {
		panic(err)
	}

	return db
}

func initUserHandler(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	uh := web.NewUserHandler(us)
	uh.RegisterRoutes(server)

}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		// AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "webook.com")
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}), func(ctx *gin.Context) {
		println("my middleware")
	})

	login := &middleware.LoginMiddlewareBuilder{}

	store := cookie.NewStore([]byte("secret"))

	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
	return server
}

func main() {
	initDB()
	server := initWebServer()
	initUserHandler(initDB(), server)
	server.Run(":8080")
}
