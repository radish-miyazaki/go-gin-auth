package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/radish-miyazaki/go-auth/db"
	"github.com/radish-miyazaki/go-auth/routes"
	"time"
)

func main() {
	// DBの接続開始
	db.Connect()

	// Routerの生成
	r := gin.Default()

	// CORSの設定
	r.Use(cors.New(cors.Config{
		// アクセスを許可したいアクセス元
		AllowOrigins: []string{
			"http://localhost:8000",
			"http://localhost:3000",
		},
		// アクセスを許可したいHTTPメソッド
		AllowMethods: []string{
			"DELETE",
			"PUT",
			"POST",
			"GET",
			"OPTIONS",
		},
		// 許可したいHTTPリクエストヘッダ
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		// cookieなどの情報を必要とするかどうか
		AllowCredentials: true,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 24 * time.Hour,
	}))

	routes.Setup(r)

	// サーバー起動
	if err := r.Run(":8080"); err != nil {
		panic("the server couldn't be started")
	}
}
