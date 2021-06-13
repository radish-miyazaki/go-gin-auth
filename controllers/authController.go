package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/radish-miyazaki/go-auth/db"
	"github.com/radish-miyazaki/go-auth/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

type Claims struct {
	jwt.StandardClaims
}

func Register(c *gin.Context) {
	data := map[string]string{}

	// リクエストのBodyを変数に格納
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
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	u := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  password,
	}
	// DBに保存
	db.DB.Create(&u)

	// 作成したユーザとステータスを返す
	c.JSON(http.StatusCreated, u)
}

func Login(c *gin.Context) {
	data := map[string]string{}

	// リクエストのBodyを変数に格納
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var u models.User
	db.DB.Where("email = ?", data["email"]).First(&u)

	if u.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not Found!",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(data["password"])); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect Password!",
		})
		return
	}

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Issuer:    strconv.Itoa(int(u.ID)),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "An unexpected error has occurred",
		})
		return
	}
	c.SetCookie("jwt", token, time.Now().Add(time.Hour*24).Second(), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"jwt": token,
	})

}

func Logout(c *gin.Context) {
	// Cookieを初期化
	c.SetCookie("jwt", "", time.Now().Add(-time.Hour).Second(), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func User(c *gin.Context) {
	cookie, _ := c.Cookie("jwt")

	// tokenを作成
	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	// tokenの作成に失敗した、またはtokenが失効している場合
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthenticated!",
		})
		return
	}
	claims := token.Claims.(*Claims)

	var u models.User
	db.DB.Where("id = ?", claims.Issuer).First(&u)

	c.JSON(http.StatusOK, u)
}
