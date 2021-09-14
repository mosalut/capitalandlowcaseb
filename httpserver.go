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
	r.GET("/worthdeposits", getWorthDeposits)
	r.GET("/drawns", getDrawns)
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

func getWorthDeposits (c *gin.Context) {
	account := c.Query("account")
	key := c.Query("key")
	if checkSignInOK(c, account, key) {
		worthDeposits, err := getWorthDepositData()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(worthDeposits)
		c.JSON(http.StatusOK, gin.H {
			"success": true,
			"message": "ok",
			"data": worthDeposits,
		})
	}
}

func getDrawns (c *gin.Context) {
	account := c.Query("account")
	key := c.Query("key")
	if checkSignInOK(c, account, key) {
		drawns, err := getDrawnData()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(drawns)
		c.JSON(http.StatusOK, gin.H {
			"success": true,
			"message": "ok",
			"data": drawns,
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
				getApyRate(c) // 年化收益率
				getCfilToFil(c) // CfilToFil
				getCapitalB(c) // 可流通量b
				getLowcaseB(c) // 锁仓量B
				getLoss(c) // 损耗值
				getLockedFilNode(c) // B锁仓量投资FIL节点
				getDrawnFil(c) // 累计已提取FIL
			//	getDrawnCfil(c) // 已提取CFIL
			//	getRewardedFaci(c) // 已奖励Faci
			//	getFaciTotal(c) // Faci总发行量

				time.Sleep(time.Second)
				c.Writer.(http.Flusher).Flush()

				return true
			})

			break
		}
	}
}

// 年化收益率
func getApyRate(c *gin.Context) {
	c.SSEvent("apyrate", gin.H {
		"success": true,
		"message": "ok",
		"data": 800,
	})
}

// CfilToFil
func getCfilToFil(c *gin.Context) {
	c.SSEvent("cfiltofil", gin.H {
		"success": true,
		"message": "ok",
		"data": 1.2,
	})
}

// 可流通量b
func getCapitalB(c *gin.Context) {
	c.SSEvent("capitalb", gin.H {
		"success": true,
		"message": "ok",
		"data": 800,
	})
}

// 锁仓量B
func getLowcaseB(c *gin.Context) {
	c.SSEvent("lowcaseb", gin.H {
		"success": true,
		"message": "ok",
		"data": 100,
	})
}

// 损耗值
func getLoss(c *gin.Context) {
	c.SSEvent("loss", gin.H {
		"success": true,
		"message": "ok",
		"data": 0.321,
	})
}

// B锁仓量投资FIL节点
func getLockedFilNode(c *gin.Context) {
	c.SSEvent("lockedfilnode", gin.H {
		"success": true,
		"message": "ok",
		"data": 0.321,
	})
}

// 累计已提取FIL
func getDrawnFil(c *gin.Context) {
	c.SSEvent("drawnfil", gin.H {
		"success": true,
		"message": "ok",
		"data": 1000,
	})
}

/*
// 已提取CFIL
func getDrawnCfil(c *gin.Context) {
	c.SSEvent("drawncfil", gin.H {
		"success": true,
		"message": "ok",
		"data": 1000,
	})
}

// 已奖励Faci
func getRewardedFaci(c *gin.Context) {
	c.SSEvent("rewardedfaci", gin.H {
		"success": true,
		"message": "ok",
		"data": 1000,
	})
}

// Faci总发行量
func getFaciTotal(c *gin.Context) {
	c.SSEvent("facitotal", gin.H {
		"success": true,
		"message": "ok",
		"data": 1000.23,
	})
}
*/
