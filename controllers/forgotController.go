package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/radish-miyazaki/go-auth/db"
	"github.com/radish-miyazaki/go-auth/models"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"net/smtp"
)

func Forgot(c *gin.Context) {

	data := map[string]string{}

	// リクエストのBodyを変数に格納
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// パスワードリセット用のtoken生成
	token := RandStringRunes(12)

	// パスワードリセットモデルを生成
	passwordReset := models.PasswordReset{
		Email: data["email"],
		Token: token,
	}

	db.DB.Create(&passwordReset)

	// パスワードリセットのメール送信
	from := "admin@example.com"
	to := []string{
		data["email"],
	}
	url := "http://localhost:3000/reset/" + token
	message := []byte("Click <a href=\"" + url + "\">here</a> to reset your password")
	// use Mailhog for test
	if err := smtp.SendMail("0.0.0.0:1025", nil, from, to, message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func Reset(c *gin.Context) {
	data := map[string]string{}

	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// パスワードの確認チェック
	if data["password"] != data["password_confirm"] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password do not match!",
		})
		return
	}

	// パスワードを暗号化
	var passwordReset = models.PasswordReset{}

	// 一番直近のパスワードリセットモデルを取得し、トークンを確認
	if err := db.DB.Where("token = ?", data["token"]).Last(&passwordReset); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid token!",
		})
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	db.DB.Model(&models.User{}).Where("email = ?", passwordReset.Email).Update("password", password)

	// TODO: Password変更後に安全性を保つためDBからデータを消去する or Tokenを失効する

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// RandStringRunes ランダムな文字列を生成する
func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
