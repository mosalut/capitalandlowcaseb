package main

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func runHTTP() {
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/signin", signIn)
	r.GET("/getcapitalb", getCapitalB)
	r.GET("/getlowcaseb", getLowcaseB)
	r.GET("/testapi", testAPI)
	r.Run(":8888")
}

func signIn(c *gin.Context) {
	accountP := c.PostForm("account")
	passwordP := c.PostForm("password")
	fmt.Println(accountP)
	fmt.Println(passwordP)

	if accountP != account {
		c.String(http.StatusOK, "Invalid&nbsp;account")
		return

		if passwordP != password {
			c.String(http.StatusOK, "Invalid&nbsp;password")
			return
		}
	}

	token := &token_T {account: accountP, timestamp: time.Now().Unix()}

	data := []byte(accountP)
	data = append(data, uint64ToBytes(uint64(token.timestamp))...)
	token.hash = md5.Sum(data)

	key := hex.EncodeToString(token.hash[:])

	tokens[key] = token

	c.Redirect(http.StatusMovedPermanently, "http://47.98.204.151/signinsuccess.html?key=" + key + "&account=" + accountP)
}

func getCapitalB(c *gin.Context) {
	account := c.Query("account")
	key := c.Query("key")

	if !validation(account, key) {
		c.JSON(http.StatusOK, gin.H {
			"success": false,
			"message": "not auth",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"success": true,
		"message": "ok",
		"data": gin.H {
			"locked": 500,
			"unlocked": 400,
		},
	})
}

func getLowcaseB(c *gin.Context) {
	account := c.Query("account")
	key := c.Query("key")

	if !validation(account, key) {
		c.JSON(http.StatusOK, gin.H {
			"success": false,
			"message": "not auth",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"success": true,
		"message": "ok",
		"data": 100,
	})
}

func testAPI(c *gin.Context) {
	return
}
