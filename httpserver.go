package main

import (
	"io"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"math/big"
	"time"
	"fmt"
	"log"

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

type conn_T struct {
	*token_T
	cfToFCh chan string
	lowcaseBCh chan string
	lossCh chan string
	drawnFilCh chan string
	filNodesCh chan map[string]cacheFilNode_T
}

type gettingCode_T byte

func (gC *gettingCode_T) wait(ip string) {
	time.Sleep(time.Second * 55)
	delete(gettingCodes, ip)
}

var gettingCodes = make(map[string]*gettingCode_T)

func runHTTP() {
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/signin", signIn)
	r.POST("/checksignin", checkSignIn)
	r.GET("/sse", sseHandler)
	r.GET("/testapi", testAPI)
	r.GET("/cirulations", getCirulations)
	r.GET("/worthdeposits", getWorthDeposits)
	r.GET("/fildrawns", getFilDrawns)
	r.GET("/cfildrawns", getCfilDrawns)
	r.GET("/signout", signOut)
	r.GET("/", initData)
	r.POST("/code", getCode)
	r.Run(PORT)
}

func getCode(c *gin.Context) {
	return
	_, ok := gettingCodes[c.ClientIP()]
	if ok {
		fmt.Println("code waiting")
		return
	}

	accountP := c.PostForm("account")

	code, err := rand.Int(rand.Reader, big.NewInt(0x1000000))
	if err != nil {
		log.Println(err)
		return
	}
	codeS := hex.EncodeToString(code.Bytes())

	sms := &sms_T {
		accountP,
		"斯年云电子商务",
		"SMS_224640018",
		`{"code":"` + codeS + `"}`,
		"",
	}

	err = sms.send()
	if err != nil {
		c.JSON(http.StatusOK, gin.H {
			"success": false,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H {
		"success": true,
		"message": "ok",
	})

	smsM[accountP] = &smsStorage_T{codeS, time.Now().Unix()}
	go func() {
		time.Sleep(time.Second * 300)
		delete(smsM, accountP)
	}()

	gettingCodes[c.ClientIP()] = nil
	go gettingCodes[c.ClientIP()].wait(c.ClientIP())
}

func signIn(c *gin.Context) {
	accountP := c.PostForm("account")
	codeP := c.PostForm("code")
	fmt.Println(accountP)
	fmt.Println(codeP)

	/*
	accountM, ok := smsM[accountP]
	if !ok {
		c.String(http.StatusOK, "Invalid account")
		return
	}

	if accountM.code != codeP {
		c.String(http.StatusOK, "Invalid code")
		return
	}
	*/

	token := &token_T {account: accountP, timestamp: time.Now().Unix(), networking: c.ClientIP() + WEBPORT}

	data := []byte(accountP)
	data = append(data, uint64ToBytes(uint64(token.timestamp))...)
	token.hash = md5.Sum(data)

	key := hex.EncodeToString(token.hash[:])

	conn := &conn_T {
		token,
		make(chan string),
		make(chan string),
		make(chan string),
		make(chan string),
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
				"lowcaseb": cacheLowcaseB,
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

func getFilDrawns (c *gin.Context) {
	account := c.Query("account")
	key := c.Query("key")
	if checkSignInOK(c, account, key) {
		drawns, err := getFilDrawnsData()
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

func getCfilDrawns (c *gin.Context) {
	account := c.Query("account")
	key := c.Query("key")
	if checkSignInOK(c, account, key) {
		drawns, err := getCfilDrawnsData()
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
				case data := <-conn.lowcaseBCh:
					getLowcaseB(c, data) // 可流通量b
				case data := <-conn.lossCh:
					getLoss(c, data) // 损耗值
				case data := <-conn.drawnFilCh:
					getDrawnFil(c, data) // 累计已提取FIL
				case data := <-conn.filNodesCh:
					getFilNodes(c, data)
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

// 可流通量b
func getLowcaseB(c *gin.Context, data string) {
	c.SSEvent("lowcaseb", gin.H {
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
		close(conn.lowcaseBCh)
		close(conn.lossCh)
		close(conn.drawnFilCh)
		delete(conns, key)
	}
}
