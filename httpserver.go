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

type connection_I interface {
	disconnect(key string)
}

type conn_T struct {
	*token_T
	cfToFCh chan string
	lowcaseBCh chan string
	lossCh chan string
	drawnFilCh chan string
	filNodesCh chan map[string]cacheFilNode_T

	cirulationCh chan []float64
	worthDepositCh chan []float64
	filDrawnsCh chan []float64
	cfilDrawnsCh chan []float64

	pingCh chan byte
}

func (conn *conn_T)disconnect(key string) {
	_, ok := conns[key]
	if ok {
		close(conn.cfToFCh)
		close(conn.lowcaseBCh)
		close(conn.lossCh)
		close(conn.drawnFilCh)
		close(conn.filNodesCh)

		close(conn.cirulationCh)
		close(conn.worthDepositCh)
		close(conn.filDrawnsCh)
		close(conn.cfilDrawnsCh)

		close(conn.pingCh)
		delete(conns, key)
	}
}

type conn2_T struct {
	cirulationCh chan []float64
	worthDepositCh chan []float64

	pingCh chan byte
}

func (conn *conn2_T)disconnect(key string) {
	_, ok := conns2[key]
	if ok {
		close(conn.cirulationCh)
		close(conn.worthDepositCh)

		close(conn.pingCh)
		delete(conns, key)
	}
}

var conns2 map[string]connection_I

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
	r.GET("/sse2", sseHandler2)
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
		make(chan []float64),
		make(chan []float64),
		make(chan []float64),
		make(chan []float64),
		make(chan byte),
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
//	account := c.Query("account")
//	key := c.Query("key")
//	if checkSignInOK(c, account, key) {
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
//	}
}

func getWorthDeposits (c *gin.Context) {
//	account := c.Query("account")
//	key := c.Query("key")
//	if checkSignInOK(c, account, key) {
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
//	}
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
		conns[key].disconnect(key)

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


	for _, cc := range conns {
		conn := cc.(*conn_T)
		if conn.networking == c.ClientIP() + WEBPORT {
			c.Stream(func(w io.Writer) bool {
				select {
				case data := <-conn.cfToFCh:
					pushCfilToFil(c, data) // CfilToFil
				case data := <-conn.lowcaseBCh:
					pushLowcaseB(c, data) // 可流通量b
				case data := <-conn.lossCh:
					pushLoss(c, data) // 损耗值
				case data := <-conn.drawnFilCh:
					pushDrawnFil(c, data) // 累计已提取FIL
				case data := <-conn.filNodesCh:
					pushFilNodes(c, data)
				case data := <-conn.cirulationCh:
					pushCirulations(c, data)
				case data := <-conn.worthDepositCh:
					pushWorthDeposits(c, data)
				case data := <-conn.filDrawnsCh:
					pushFilDrawns(c, data)
				case data := <-conn.cfilDrawnsCh:
					pushCfilDrawns(c, data)
				case <-conn.pingCh:
					pong(c)
				}

				c.Writer.(http.Flusher).Flush()
				return true
			})
		}
	}
}

func sseHandler2(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Cache-Control", "no-cache")
//	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	key := c.ClientIP()
	cc, ok := conns2[key]
	if ok {
		cc.disconnect(key)
	}

	conn := &conn2_T {
		make(chan []float64),
		make(chan []float64),
		make(chan byte),
	}
	conns2[key] = conn

	c.Stream(func(w io.Writer) bool {
		select {
		case data := <-conn.cirulationCh:
			pushCirulations(c, data) // CfilToFil
		case data := <-conn.worthDepositCh:
			pushWorthDeposits(c, data) // 可流通量b
		case <-conn.pingCh:
			pong(c)
		}

		c.Writer.(http.Flusher).Flush()
		return true
	})
}

// 年化收益率
func pushApyRate(c *gin.Context, data string) {
	c.SSEvent("apyrate", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// CfilToFil
func pushCfilToFil(c *gin.Context, data string) {
	c.SSEvent("cfiltofil", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// 可流通量b
func pushLowcaseB(c *gin.Context, data string) {
	c.SSEvent("lowcaseb", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// 损耗值
func pushLoss(c *gin.Context, data string) {
	c.SSEvent("loss", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// 累计已提取FIL
func pushDrawnFil(c *gin.Context, data string) {
	c.SSEvent("drawnfil", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// B锁仓量投资FIL节点
func pushFilNodes(c *gin.Context, data map[string]cacheFilNode_T) {
	c.SSEvent("filNodes", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

func pushCirulations(c *gin.Context, data []float64) {
	c.SSEvent("cirulations", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

func pushWorthDeposits(c *gin.Context, data []float64) {
	c.SSEvent("worthdeposits", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

func pushFilDrawns(c *gin.Context, data []float64) {
	c.SSEvent("fildrawns", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

func pushCfilDrawns(c *gin.Context, data []float64) {
	c.SSEvent("cfildrawns", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

func startPing(conns map[string]connection_I) {
	for _, c := range conns {
		go ping(c)
	}
}

func ping(c connection_I) {
	for {
		switch c.(type) {
		case *conn_T:
			c.(*conn_T).pingCh <- byte(0)
		case *conn2_T:
			c.(*conn2_T).pingCh <- byte(0)
		}
		modTime := time.Now().Unix() % 10
		time.Sleep(time.Second * time.Duration(TIME_CFILTOFIL - modTime))
	}
}

func pong(c *gin.Context) {
	c.SSEvent("pong", gin.H {
		"success": true,
		"message": "ok",
	})
}
