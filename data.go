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

/*
type cacheB_T struct {
	CapitalB string `json:"capitalb"`
	LowcaseB string `json:"lowcaseb"`
}
*/

type cacheFilNode_T struct {
	Address string `json:"address"`
	Balance float64 `json:"balance"`
	WorkerBalance float64 `json:"workerbalance"`
	QualityAdjPower float64 `json:"qualityadjpower"`
	AvailableBalance float64 `json:"availableBalance"`
	Pledge float64 `json:"pledge"`
	VestingFunds float64 `json:"vestingFunds"`
	SingleT float64 `json:"singlet"`
}

// var cacheB cacheB_T
var cacheCfToF float64
var cacheLoss float64
var cacheDrawnFil float64
var cacheFilNodes = make(map[string]cacheFilNode_T)
var cacheLowcaseB float64
var cacheCapitalB float64

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

		id, err := insertHourData(cacheLowcaseB, cacheDrawnFil)
		if err != nil {
			log.Error(err)
			continue
		}

		err = insertFilNodes(id, cacheFilNodes)
		if err != nil {
			log.Error(err)
			continue
		}

		for k, v := range cacheFilNodes {
			log.Info(k, v)
		}

		for _, c := range conns {
			cc := c.(*conn_T)
			go func() {
				defer recoverPanic()
				cc.cfToFCh <- cacheCfToF
				cc.lowcaseBCh <- cacheLowcaseB
				cc.lossCh <- cacheLoss
				cc.drawnFilCh <- cacheDrawnFil
				cc.filNodesCh <- cacheFilNodes

				cirulations, err := getCirulationData()
				if err != nil {
					log.Error(err)
				}
				cc.cirulationCh <- cirulations

				worthDeposits, err := getWorthDepositData()
				if err != nil {
					log.Error(err)
				}
				cc.worthDepositCh <- worthDeposits

				filDrawns, err := getFilDrawnsData()
				if err != nil {
					log.Error(err)
				}
				cc.filDrawnsCh <- filDrawns

				cfilDrawns, err := getCfilDrawnsData()
				if err != nil {
					log.Error(err)
				}
				cc.cfilDrawnsCh <- cfilDrawns
			}()
		}

		for _, c := range conns2 {
			cc := c.(*conn2_T)
			go func() {
				defer recoverPanic()
				filDrawns, err := getFilDrawnsData()
				if err != nil {
					log.Error(err)
				}
				cc.filDrawnsCh <- filDrawns

				cfilDrawns, err := getCfilDrawnsData()
				if err != nil {
					log.Error(err)
				}
				cc.cfilDrawnsCh <- cfilDrawns
			}()
		}
	}

}

// 最新总资产/日
func getCirulationData() ([]float64, error) {
	balances := make([]float64, 24, 24)

	for i, _ := range balances {
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
		balances[i] = integer + decimal / 100000
	}

	return balances, nil
}

// 24小时净存取
func getWorthDepositData() ([]float64, error) {
	balances := make([]float64, 24, 24)

	for i, _ := range balances {
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
		balances[i] = integer + decimal / 100000
	}

	return balances, nil
}

// 24小时FIL提现
func getFilDrawnsData() ([]float64, error) {
	balances := make([]float64, 24, 24)

	for i, _ := range balances {
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
		balances[i] = integer + decimal / 100000
	}

	return balances, nil
}

// 24小时CFIL提现
func getCfilDrawnsData() ([]float64, error) {
	balances := fibonache()

	return balances, nil
}

func requestCfilToFil(wg *sync.WaitGroup) {
	defer wg.Done()
	cacheCfToF = 1.2
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
	cacheLowcaseB = lowcaseB / 1e18
}

func requestLoss(wg *sync.WaitGroup) {
	defer wg.Done()
	cacheLoss = 0.321
}

func requestDrawnFil(wg *sync.WaitGroup) {
	defer wg.Done()
	cacheDrawnFil = 1000
}

func requestFilNodes(wg *sync.WaitGroup) {
	defer wg.Done()
	cacheCapitalB = 0

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

		cacheFilNode := cacheFilNode_T{}
		cacheFilNode.Address = data["miner"].(map[string]interface{})["owner"].(map[string]interface{})["address"].(string)

		balance, err := strconv.ParseFloat(data["balance"].(string), 64)
		if err != nil {
			log.Error(err)
			return
		}
		cacheFilNode.Balance = balance / 1e18

		availableBalance, err := strconv.ParseFloat(data["miner"].(map[string]interface{})["availableBalance"].(string), 64)
		if err != nil {
			log.Error(err)
			return
		}
		cacheFilNode.AvailableBalance = availableBalance / 1e18

		pledge, err := strconv.ParseFloat(data["miner"].(map[string]interface{})["sectorPledgeBalance"].(string), 64)
		if err != nil {
			log.Error(err)
			return
		}
		cacheFilNode.Pledge = pledge / 1e18

		vestingFunds, err := strconv.ParseFloat(data["miner"].(map[string]interface{})["vestingFunds"].(string), 64)
		if err != nil {
			log.Error(err)
			return
		}
		cacheFilNode.VestingFunds = vestingFunds / 1e18

		workerBalance, err := strconv.ParseFloat(data["miner"].(map[string]interface{})["worker"].(map[string]interface{})["balance"].(string), 64)
		if err != nil {
			log.Error(err)
			return
		}
		cacheFilNode.WorkerBalance = workerBalance / 1e18

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
		cacheFilNode.QualityAdjPower = qualityAdjPower / 1125899906842624

		if active == 0 {
			cacheFilNode.SingleT = 0
		} else {
			totalRewards, err := strconv.ParseFloat(params["totalRewards"].(string), 64)
			if err != nil {
				log.Error(err)
				return
			}
			totalRewards /= 1e18

			cacheFilNode.SingleT = totalRewards / active * 16
		}

		cacheCapitalB += balance + workerBalance
		cacheFilNodes[nodeKey] = cacheFilNode
	}
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
