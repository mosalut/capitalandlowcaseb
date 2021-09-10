package main

import (
	"io"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

const (
//	PORT = ":8888"
	PORT = ":8887"
//	WEBPORT = ":80"
	WEBPORT = ":10000"
)

type event_T struct {
	ID      int
	Message string
}

func runHTTP() {
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/signin", signIn)
	r.POST("/checksignin", checkSignIn)
	r.GET("/sse", sseHandler)
	r.GET("/testapi", testAPI)
	r.GET("/cirulations", getCirulations)
	r.Run(PORT)
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

	token := &token_T {account: accountP, timestamp: time.Now().Unix(), networking: c.ClientIP() + WEBPORT}

	data := []byte(accountP)
	data = append(data, uint64ToBytes(uint64(token.timestamp))...)
	token.hash = md5.Sum(data)

	key := hex.EncodeToString(token.hash[:])

	tokens[key] = token

	c.Redirect(http.StatusMovedPermanently, "http://47.98.204.151" + WEBPORT + "/signinsuccess.html?key=" + key + "&account=" + accountP)
}

func checkSignInOK(c *gin.Context, account, key string) bool {
	if !validation(account, key, c.ClientIP() + WEBPORT) {
		c.JSON(http.StatusOK, gin.H {
			"success": false,
			"message": "not auth",
		})
		return false
	}

	return true
}

func checkSignIn(c *gin.Context) {
	account := c.PostForm("account")
	key := c.PostForm("key")
	if checkSignInOK(c, account, key) {
		c.JSON(http.StatusOK, gin.H {
			"success": true,
			"message": "ok",
		})
	}
}

func getCirulations (c *gin.Context) {
	account := c.Query("account")
	key := c.Query("key")
	if checkSignInOK(c, account, key) {
		cirulations, err := getCirulationData()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(cirulations)
		c.JSON(http.StatusOK, gin.H {
			"success": true,
			"message": "ok",
			"data": cirulations,
		})
	}
}

func testAPI(c *gin.Context) {
	return
}

func sseHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Cache-Control", "no-cache")
//	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	for _, token := range tokens {
		if token.networking == c.ClientIP() + WEBPORT {
			c.Stream(func(w io.Writer) bool {
				getCapitalB(c)
				getLowcaseB(c)
				getPutIn(c)
				getSettled(c)
				getRewarded(c)
				getFilPrice(c)

				time.Sleep(time.Second)
				c.Writer.(http.Flusher).Flush()

				return true
			})

			break
		}
	}
}

func getCapitalB(c *gin.Context) {
	c.SSEvent("capitalb", gin.H {
		"success": true,
		"message": "ok",
		"data": gin.H {
			"locked": 100,
			"unlocked": 700,
		},
	})
}

func getLowcaseB(c *gin.Context) {
	c.SSEvent("lowcaseb", gin.H {
		"success": true,
		"message": "ok",
		"data": 100,
	})
}

func getPutIn(c *gin.Context) {
	c.SSEvent("totalputin", gin.H {
		"success": true,
		"message": "ok",
		"data": 10000,
	})
}

func getSettled(c *gin.Context) {
	c.SSEvent("settled", gin.H {
		"success": true,
		"message": "ok",
		"data": 1000,
	})
}

func getRewarded(c *gin.Context) {
	c.SSEvent("rewarded", gin.H {
		"success": true,
		"message": "ok",
		"data": 1000,
	})
}

func getFilPrice(c *gin.Context) {
	c.SSEvent("filprice", gin.H {
		"success": true,
		"message": "ok",
		"data": 1000.23,
	})
}
