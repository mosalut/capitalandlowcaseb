package main

import (
	"bytes"
	"sync"
	"net/http"
	"encoding/json"
	"strconv"
	"math/big"
	"time"
)

const (
	BSC_API_URL = "https://api.bscscan.com/api"
	BSC_API_KEY = "4C2328SSW63VFIWQZMYWD2NZERUUGN3VT1"
//	BSC_WALLET_CAPITAL = "0x8A19846c7e057DBE6D77419BF1c864DAd0065d45"
	BSC_WALLET_LOWCASE = "0xFF6C223b9Dc11F247B16CED5e8a996DC4deA793E"

	PAGE_URL = "https://filfox.info/api/v1/address/"
	PAGE_URL_SUB_MINING = "/mining-stats?duration=24h"

	CONTRACT_URL = "https://data-seed-prebsc-1-s1.binance.org:8545/"
	CONTRACT_ADDRESS_FREEERC20 = "0xfcd60BC1c495cCb9A5539758CCc18766310653F6"
	CONTRACT_ADDRESS_VAULT = "0x8F7DEB527EE6C06cD34C54bA6d424C62c9634f61"
)

type filNode_T struct {
	Address string `json:"address"`
	Balance float64 `json:"balance"`
	WorkerBalance float64 `json:"workerbalance"`
	QualityAdjPower float64 `json:"qualityadjpower"`
	AvailableBalance float64 `json:"availableBalance"`
	Pledge float64 `json:"pledge"`
	VestingFunds float64 `json:"vestingFunds"`
	SingleT float64 `json:"singlet"`
}

type curve_T struct {
	CreateTime time.Time `json:"createtime"`
	Value float64 `json:"value"`
}

type cache_T struct {
	CfToF float64 `json:"cftof"`
	Loss float64 `json:"loss"`
	DrawnFil float64 `json:"fil"`
	FilNodes map[string]filNode_T `json:"filnodes"`
	LowcaseB float64 `json:"lowcaseb"`
	CapitalB float64 `json:"capitalb"`
}

var cache = cache_T {FilNodes: make(map[string]filNode_T)}

type contractParams_T struct {
	Method string `json:"method"`
	Sync int `json:"id"`
	JsonRPC string `json:"jsonrpc"`
	Data interface{} `json:"params"`
}

func listenRequests() {
	for {
		modTime := time.Now().Unix() % config.period
		time.Sleep(time.Second * time.Duration(config.period - modTime))

		createTime := time.Now()

		wg := &sync.WaitGroup{}
		wg.Add(4)
		go requestB(wg)
		go requestLoss(wg)
		go requestDrawnFil(wg)
		go requestFilNodes(wg)
		wg.Wait()

		id, err := insertHourData(createTime, cache.LowcaseB, cache.CapitalB, cache.DrawnFil)
		if err != nil {
			log.Error(err)
			continue
		}

		err = insertFilNodes(id, cache.FilNodes)
		if err != nil {
			log.Error(err)
			continue
		}

		for k, v := range cache.FilNodes {
			log.Info(k, v)
		}

		for _, c := range conns {
			cc := c.(*conn_T)
			go func() {
				defer recoverPanic()
				cc.lowcaseBCh <- &curve_T{createTime, cache.LowcaseB}
				cc.capitalBCh <- &curve_T{createTime, cache.CapitalB}
				cc.lossCh <- cache.Loss
				cc.drawnFilCh <- &curve_T{createTime, cache.DrawnFil}
				cc.filNodesCh <- cache.FilNodes
			}()
		}

		for _, c := range conns2 {
			cc := c.(*conn2_T)
			go func() {
				defer recoverPanic()
				cc.capitalBCh <- &curve_T{createTime, cache.CapitalB}
				cc.drawnFilCh <- &curve_T{createTime, cache.DrawnFil}
			}()
		}
	}

}

func listenRequests5min() {
	for {
		modTime := time.Now().Unix() % config.period5
		time.Sleep(time.Second * time.Duration(config.period5 - modTime))

		createTime := time.Now()

		requestCfilToFil()

		err := insert5MinsData(createTime, cache.CfToF)
		if err != nil {
			log.Error(err)
			continue
		}

		for _, c := range conns {
			cc := c.(*conn_T)
			go func() {
				defer recoverPanic()
				cc.cfToFCh <- &curve_T{createTime, cache.CfToF}
			}()
		}

		for _, c := range conns2 {
			cc := c.(*conn2_T)
			go func() {
				defer recoverPanic()
				cc.cfToFCh <- &curve_T{createTime, cache.CfToF}
			}()
		}
	}
}

// CFIL净值/5分钟
func getCfilDrawnsData() []curve_T {
	values := fibonache()

	curve := make([]curve_T, 288, 288)
	for i, _ := range curve {
		curve[i].Value = values[i]
		curve[i].CreateTime = time.Unix(time.Now().Unix() / config.period * config.period, 0)
	}

	return curve
}

func requestCfilToFil() {
	pData := []interface{}{map[string]string{"to": CONTRACT_ADDRESS_VAULT, "data": "0x4661cb8a"}, "latest"}
	params := &contractParams_T{"eth_call", 1, "2.0", pData}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(params)
	resp, err := http.Post(CONTRACT_URL, "application/json", buffer)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Error(err)
		return
	}

//	log.Info(data)

	cFtoF := big.NewFloat(0)
	B1e18 := big.NewFloat(1e18)
	_, ok := cFtoF.SetString(data["result"].(string))
	if !ok {
		log.Error("can not turn to *big.Float")
		return
	}

	cFtoF.Quo(cFtoF, B1e18)

//	log.Info(cFtoF)

	cache.CfToF, _ = cFtoF.Float64()
}

func requestB(wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(BSC_API_URL + "?module=account&action=balance&address=" + BSC_WALLET_LOWCASE + "&apikey=" + BSC_API_KEY)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Error(err)
		return
	}

	lowcaseB, err := strconv.ParseFloat(data["result"].(string), 64)
	if err != nil {
		log.Error(err)
		return
	}
	cache.LowcaseB = lowcaseB / 1e18
}

func requestLoss(wg *sync.WaitGroup) {
	defer wg.Done()
	cache.Loss = 0.321
}

func requestDrawnFil(wg *sync.WaitGroup) {
	defer wg.Done()
	cache.DrawnFil = 1000
}

func requestFilNodes(wg *sync.WaitGroup) {
	defer wg.Done()
	cache.CapitalB = 0

	for _, nodeKey := range config.nodes {
		resp, err := http.Get(PAGE_URL + nodeKey)
		if err != nil {
			log.Error(err)
			return
		}
		defer resp.Body.Close()

		data := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			log.Error(err)
			return
		}

		filNode := filNode_T{}
		filNode.Address = data["miner"].(map[string]interface{})["owner"].(map[string]interface{})["address"].(string)

		balance, err := strconv.ParseFloat(data["balance"].(string), 64)
		if err != nil {
			log.Error(err)
			return
		}
		filNode.Balance = balance / 1e18

		availableBalance, err := strconv.ParseFloat(data["miner"].(map[string]interface{})["availableBalance"].(string), 64)
		if err != nil {
			log.Error(err)
			return
		}
		filNode.AvailableBalance = availableBalance / 1e18

		pledge, err := strconv.ParseFloat(data["miner"].(map[string]interface{})["sectorPledgeBalance"].(string), 64)
		if err != nil {
			log.Error(err)
			return
		}
		filNode.Pledge = pledge / 1e18

		vestingFunds, err := strconv.ParseFloat(data["miner"].(map[string]interface{})["vestingFunds"].(string), 64)
		if err != nil {
			log.Error(err)
			return
		}
		filNode.VestingFunds = vestingFunds / 1e18

		workerBalance, err := strconv.ParseFloat(data["miner"].(map[string]interface{})["worker"].(map[string]interface{})["balance"].(string), 64)
		if err != nil {
			log.Error(err)
			return
		}
		filNode.WorkerBalance = workerBalance / 1e18

		active := data["miner"].(map[string]interface{})["sectors"].(map[string]interface{})["active"].(float64)

		respMining, err := http.Get(PAGE_URL + nodeKey + PAGE_URL_SUB_MINING)
		if err != nil {
			log.Error(err)
			return
		}
		defer respMining.Body.Close()

		params := make(map[string]interface{})
		err = json.NewDecoder(respMining.Body).Decode(&params)
		if err != nil {
			log.Error(err)
			return
		}

		qualityAdjPower, err := strconv.ParseFloat(data["miner"].(map[string]interface{})["qualityAdjPower"].(string), 64)
		if err != nil {
			log.Error(err)
			return
		}
		filNode.QualityAdjPower = qualityAdjPower / 1125899906842624

		if active == 0 {
			filNode.SingleT = 0
		} else {
			totalRewards, err := strconv.ParseFloat(params["totalRewards"].(string), 64)
			if err != nil {
				log.Error(err)
				return
			}
			totalRewards /= 1e18

			filNode.SingleT = totalRewards / active * 16
		}

		cache.CapitalB += balance + workerBalance
		cache.FilNodes[nodeKey] = filNode
	}

	cache.CapitalB /= 1e18
}

func fibonache() []float64 {
	balances := make([]float64, 288, 288)
	var value float64
	for i := 0; i < len(balances); i++ {
		value += float64(i)
		balances[i] = value
	}
	return balances
}

func recoverPanic() {
	err := recover()
	if err != nil {
		log.Error(err)
	}
}
