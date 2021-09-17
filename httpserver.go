package main

import (
	"io"
	"crypto/md5"
	"encoding/hex"
//	"encoding/json"
	"net/http"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

const (
	PORT = ":8888"
//	PORT = ":8887"
	WEBPORT = ":80"
//	WEBPORT = ":10000"
)

type event_T struct {
	ID      int
	Message string
}

type conn_T struct {
	*token_T
	cfToFCh chan string
	bCh chan cacheB_T
	lossCh chan string
	drawnFilCh chan string
//	apyRateCh chan string 
	filNodesCh chan map[string]cacheFilNode_T
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
	r.GET("/signout", signOut)
	r.GET("/", initData)
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

	conn := &conn_T{
		token,
		make(chan string),
		make(chan cacheB_T),
		make(chan string),
		make(chan string),
	//	make(chan string),
		make(chan map[string]cacheFilNode_T),
	}
	conns[key] = conn

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

func initData(c *gin.Context) {
	account := c.Query("account")
	key := c.Query("key")
	if checkSignInOK(c, account, key) {
		c.JSON(http.StatusOK, gin.H {
			"success": true,
			"message": "ok",
			"data": gin.H {
				"apyrate": "7",
				"capitalb": cacheB.CapitalB,
				"lowcaseb": cacheB.LowcaseB,
				"cfiltofil": cacheCfToF,
				"loss": cacheLoss,
				"drawnfil": cacheDrawnFil,
				"filNodes": cacheFilNodes,
			},
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

func signOut (c *gin.Context) {
	account := c.Query("account")
	key := c.Query("key")
	if checkSignInOK(c, account, key) {
		disconnect(key)

		c.JSON(http.StatusOK, gin.H {
			"success": true,
			"message": "ok",
		})

		c.Request.Body.Close()
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


	for _, conn := range conns {
		if conn.networking == c.ClientIP() + WEBPORT {
			c.Stream(func(w io.Writer) bool {
				select {
				case data := <-conn.cfToFCh:
					getCfilToFil(c, data) // CfilToFil
				case data := <-conn.bCh:
					getB(c, data) // 可流通量b // 锁仓量B
				case data := <-conn.lossCh:
					getLoss(c, data) // 损耗值
				case data := <-conn.drawnFilCh:
					getDrawnFil(c, data) // 累计已提取FIL
				case data := <-conn.filNodesCh:
					getFilNodes(c, data)
					/*
				default:
					getApyRate(c, "7") // 年化收益率
					*/
				}

				c.Writer.(http.Flusher).Flush()
				return true
			})
		}
	}
}

// 年化收益率
func getApyRate(c *gin.Context, data string) {
	c.SSEvent("apyrate", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// CfilToFil
func getCfilToFil(c *gin.Context, data string) {
	c.SSEvent("cfiltofil", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// 可流通量b 锁仓量B
func getB(c *gin.Context, data cacheB_T) {
	c.SSEvent("bb", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// 损耗值
func getLoss(c *gin.Context, data string) {
	c.SSEvent("loss", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// 累计已提取FIL
func getDrawnFil(c *gin.Context, data string) {
	c.SSEvent("drawnfil", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
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

// B锁仓量投资FIL节点
func getFilNodes(c *gin.Context, data map[string]cacheFilNode_T) {
	c.SSEvent("filNodes", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

func disconnect(key string) {
	conn, ok := conns[key]
	if ok {
		close(conn.cfToFCh)
		close(conn.bCh)
		close(conn.lossCh)
		close(conn.drawnFilCh)
	//	close(conn.apyRateCh)
		delete(conns, key)
	}
}
