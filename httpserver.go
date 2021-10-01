package main

import (
	"io"
	"os"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"math/big"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func setHttpLog() {
	if config.mode {
		httpLog, err := os.Create(faciDir + "http.log")
		if err != nil {
			log.Fatal(err)
		}

		gin.DefaultWriter = io.MultiWriter(httpLog)
	}
}

type event_T struct {
	ID      int
	Message string
}

type connection_I interface {
	disconnect(key string)
}

type conn_T struct {
	*token_T
	cfToFCh chan float64
	lowcaseBCh chan float64
	capitalBCh chan float64
	lossCh chan float64
	drawnFilCh chan float64
	filNodesCh chan map[string]filNode_T

	lowcaseBsCh chan []curve_T
	capitalBsCh chan []curve_T
	filDrawnsCh chan []curve_T
	cfilDrawnsCh chan []curve_T

	pingCh chan byte
}

func (conn *conn_T)disconnect(key string) {
	_, ok := conns[key]
	if ok {
		close(conn.cfToFCh)
		close(conn.lowcaseBCh)
		close(conn.capitalBCh)
		close(conn.lossCh)
		close(conn.drawnFilCh)
		close(conn.filNodesCh)

		close(conn.lowcaseBsCh)
		close(conn.capitalBCh)
		close(conn.filDrawnsCh)
		close(conn.cfilDrawnsCh)

		close(conn.pingCh)
		delete(conns, key)
	}
}

type conn2_T struct {
	capitalBCh chan float64
	filDrawnsCh chan []curve_T
	cfilDrawnsCh chan []curve_T

	pingCh chan byte
}

func (conn *conn2_T)disconnect(key string) {
	_, ok := conns2[key]
	if ok {
		close(conn.capitalBCh)
		close(conn.filDrawnsCh)
		close(conn.cfilDrawnsCh)

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

var signIns [2]func(c *gin.Context)
var signIn func(c *gin.Context)

func runHTTP() {
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/signin", signIn)
	r.POST("/checksignin", checkSignIn)
	r.GET("/sse", sseHandler)
	r.GET("/sse2", sseHandler2)
	r.GET("/testapi", testAPI)
	r.GET("/capitalb", getCapitalB)
	r.GET("/curves", getCurves)
	r.GET("/signout", signOut)
	r.GET("/", initData)
	r.POST("/code", getCode)
	r.Run(":" + config.port)
}

func getCode(c *gin.Context) {
	_, ok := gettingCodes[c.ClientIP()]
	if ok {
		log.Info("code waiting")
		return
	}

	accountP := c.PostForm("account")

	code, err := rand.Int(rand.Reader, big.NewInt(0x1000000))
	if err != nil {
		log.Error(err)
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

func signInRelease(c *gin.Context) {
	accountP := c.PostForm("account")
	codeP := c.PostForm("code")
	log.Info("sign in:", accountP, codeP)

	accountM, ok := smsM[accountP]
	if !ok {
		c.String(http.StatusOK, "Invalid account")
		return
	}

	if accountM.code != codeP {
		c.String(http.StatusOK, "Invalid code")
		return
	}

	token := &token_T {account: accountP, timestamp: time.Now().Unix(), networking: c.ClientIP() + ":" + config.webPort}

	data := []byte(accountP)
	data = append(data, uint64ToBytes(uint64(token.timestamp))...)
	token.hash = md5.Sum(data)

	key := hex.EncodeToString(token.hash[:])

	conn := &conn_T {
		token,
		make(chan float64),
		make(chan float64),
		make(chan float64),
		make(chan float64),
		make(chan float64),
		make(chan map[string]filNode_T),
		make(chan []curve_T),
		make(chan []curve_T),
		make(chan []curve_T),
		make(chan []curve_T),
		make(chan byte),
	}
	conns[key] = conn

	c.Redirect(http.StatusMovedPermanently, "http://" + config.webHost + ":" + config.webPort + "/signinsuccess.html?key=" + key + "&account=" + accountP)
}

func signInDev(c *gin.Context) {
	accountP := c.PostForm("account")
	log.Info("sign in:", accountP)

	token := &token_T {account: accountP, timestamp: time.Now().Unix(), networking: c.ClientIP() + ":" + config.webPort}

	data := []byte(accountP)
	data = append(data, uint64ToBytes(uint64(token.timestamp))...)
	token.hash = md5.Sum(data)

	key := hex.EncodeToString(token.hash[:])

	conn := &conn_T {
		token,
		make(chan float64),
		make(chan float64),
		make(chan float64),
		make(chan float64),
		make(chan float64),
		make(chan map[string]filNode_T),
		make(chan []curve_T),
		make(chan []curve_T),
		make(chan []curve_T),
		make(chan []curve_T),
		make(chan byte),
	}
	conns[key] = conn

	c.Redirect(http.StatusMovedPermanently, "http://" + config.webHost + ":" + config.webPort + "/signinsuccess.html?key=" + key + "&account=" + accountP)
}

func checkSignInOK(c *gin.Context, account, key string) bool {
	if !validation(account, key, c.ClientIP() + ":" + config.webPort) {
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
				"apyrate": 7,
				"lowcaseb": cache.LowcaseB,
				"capitalb": cache.CapitalB,
				"cfiltofil": cache.CfToF,
				"loss": cache.Loss,
				"drawnfil": cache.DrawnFil,
				"filnodes": cache.FilNodes,
			},
		})
	}
}

func getCapitalB(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H {
		"success": true,
		"message": "ok",
		"data": gin.H {
			"capitalb": cache.CapitalB,
		},
	})
}

func getCurves (c *gin.Context) {
	account := c.Query("account")
	key := c.Query("key")
	if checkSignInOK(c, account, key) {
		lowcaseBs, err := getLowcaseBCurveData()
		if err != nil {
			log.Info(err)
			return
		}
		capitalBs, err := getCapitalBCurveData()
		if err != nil {
			log.Info(err)
			return
		}
		filDrawns, err := getDrawnFilCurveData()
		if err != nil {
			log.Info(err)
			return
		}
		cfToFs, err := getCfToFCurveData()
		if err != nil {
			log.Info(err)
			return
		}

		c.JSON(http.StatusOK, gin.H {
			"success": true,
			"message": "ok",
			"data": gin.H {
				"lowcasebs": lowcaseBs,
				"capitalbs": capitalBs,
				"fildrawns": filDrawns,
				"cftofs": cfToFs,
			},
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
		if conn.networking == c.ClientIP() + ":" + config.webPort {
			c.Stream(func(w io.Writer) bool {
				select {
				case data := <-conn.cfToFCh:
					pushCfilToFil(c, data) // CfilToFil
				case data := <-conn.lowcaseBCh:
					pushLowcaseB(c, data) // 流动余额b
				case data := <-conn.capitalBCh:
					pushCapitalB(c, data) // 质押余额B
				case data := <-conn.lossCh:
					pushLoss(c, data) // 损耗值
				case data := <-conn.drawnFilCh:
					pushDrawnFil(c, data) // 累计已提取FIL
				case data := <-conn.filNodesCh:
					pushFilNodes(c, data)
				case data := <-conn.lowcaseBsCh:
					pushLowcaseBs(c, data)
				case data := <-conn.capitalBsCh:
					pushCapitalBs(c, data)
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
		make(chan float64),
		make(chan []curve_T),
		make(chan []curve_T),
		make(chan byte),
	}
	conns2[key] = conn

	c.Stream(func(w io.Writer) bool {
		select {
		case data := <-conn.capitalBCh:
			pushCapitalB(c, data) // 质押余额B
			/*
		case data := <-conn.filDrawnsCh:
			pushFilDrawns(c, data) // 总提取数额FIL
		case data := <-conn.cfilDrawnsCh:
			pushCfilDrawns(c, data) // CFIL净值
			*/
		case <-conn.pingCh:
			pong(c)
		}

		c.Writer.(http.Flusher).Flush()
		return true
	})
}

// 年化收益率
func pushApyRate(c *gin.Context, data float64) {
	c.SSEvent("apyrate", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// CfilToFil
func pushCfilToFil(c *gin.Context, data float64) {
	c.SSEvent("cfiltofil", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// 流动余额b
func pushLowcaseB(c *gin.Context, data float64) {
	c.SSEvent("lowcaseb", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// 质押余额B
func pushCapitalB(c *gin.Context, data float64) {
	c.SSEvent("capitalb", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// 损耗值
func pushLoss(c *gin.Context, data float64) {
	c.SSEvent("loss", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// 累计已提取FIL
func pushDrawnFil(c *gin.Context, data float64) {
	c.SSEvent("drawnfil", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

// B锁仓量投资FIL节点
func pushFilNodes(c *gin.Context, data map[string]filNode_T) {
	c.SSEvent("filnodes", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

func pushLowcaseBs(c *gin.Context, data []curve_T) {
	c.SSEvent("lowcasebs", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

func pushCapitalBs(c *gin.Context, data []curve_T) {
	c.SSEvent("capitalbs", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

func pushFilDrawns(c *gin.Context, data []curve_T) {
	c.SSEvent("fildrawns", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

func pushCfilDrawns(c *gin.Context, data []curve_T) {
	c.SSEvent("cftofs", gin.H {
		"success": true,
		"message": "ok",
		"data": data,
	})
}

func startPing(conns map[string]connection_I) {
	for {
		for _, c := range conns {
			go ping(c)
		}
		modTime := time.Now().Unix() % config.period
		time.Sleep(time.Second * time.Duration(config.period - modTime))
	}
}

func ping(c connection_I) {
	defer recoverPanic()
	switch c.(type) {
	case *conn_T:
		c.(*conn_T).pingCh <- byte(0)
	case *conn2_T:
		c.(*conn2_T).pingCh <- byte(0)
	}
}

func pong(c *gin.Context) {
	c.SSEvent("pong", gin.H {
		"success": true,
		"message": "ok",
	})
}
