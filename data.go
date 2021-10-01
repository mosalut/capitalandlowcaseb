package main

import (
	"sync"
	"crypto/rand"
	"math/big"
	"net/http"
	"encoding/json"
	"strconv"
	"time"
)

const (
	BSC_API_URL = "https://api.bscscan.com/api"
	BSC_API_KEY = "4C2328SSW63VFIWQZMYWD2NZERUUGN3VT1"
//	BSC_WALLET_CAPITAL = "0x8A19846c7e057DBE6D77419BF1c864DAd0065d45"
	BSC_WALLET_LOWCASE = "0xFF6C223b9Dc11F247B16CED5e8a996DC4deA793E"

	PAGE_URL = "https://filfox.info/api/v1/address/"
	PAGE_URL_SUB_MINING = "/mining-stats?duration=24h"
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

func listenRequests() {
	for {
		modTime := time.Now().Unix() % config.period
		time.Sleep(time.Second * time.Duration(config.period - modTime))

		wg := &sync.WaitGroup{}
		wg.Add(5)
		go requestCfilToFil(wg)
		go requestB(wg)
		go requestLoss(wg)
		go requestDrawnFil(wg)
		go requestFilNodes(wg)
		wg.Wait()

		id, err := insertHourData(cache.LowcaseB, cache.CapitalB, cache.DrawnFil)
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
				cc.cfToFCh <- cache.CfToF
				cc.lowcaseBCh <- cache.LowcaseB
				cc.capitalBCh <- cache.CapitalB
				cc.lossCh <- cache.Loss
				cc.drawnFilCh <- cache.DrawnFil
				cc.filNodesCh <- cache.FilNodes
			}()
		}

		for _, c := range conns2 {
			cc := c.(*conn2_T)
			go func() {
				defer recoverPanic()
				cc.capitalBCh <- cache.CapitalB
				/*
				filDrawns, err := getDrawnFilCurveData()
				if err != nil {
					log.Error(err)
				}

				cfilDrawns, err := getDrawnCfilCurveData()
				if err != nil {
					log.Error(err)
				}
				cc.filDrawnsCh <- filDrawns
				cc.cfilDrawnsCh <- cfilDrawns
				*/
			}()
		}
	}

}

// 质押余额/时
func getData24() ([]curve_T, error) {
	curve := make([]curve_T, 24, 24)

	for i, _ := range curve {
		max := big.NewInt(65536)
		integerI, err := rand.Int(rand.Reader, max)
		if err != nil {
			return nil, err
		}
		decimalI, err := rand.Int(rand.Reader, max)
		if err != nil {
			return nil, err
		}

		integerF := big.NewFloat(0)
		integerF.SetInt(integerI)
		decimalF := big.NewFloat(0)
		decimalF.SetInt(decimalI)

		integer, _ := integerF.Float64()
		decimal, _ := decimalF.Float64()
		curve[i].Value = integer + decimal / 100000
		curve[i].CreateTime = time.Unix(time.Now().Unix() / config.period * config.period, 0)
	}

	return curve, nil
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

func requestCfilToFil(wg *sync.WaitGroup) {
	defer wg.Done()
	cache.CfToF = 1.2
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
