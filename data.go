package main

import (
	"sync"
	"crypto/rand"
	"math/big"
	"net/http"
	"encoding/json"
	"strconv"
	"time"
//	"fmt"
	"log"
)

const (
	TIME_CFILTOFIL = 60

	BSC_API_URL = "https://api.bscscan.com/api"
	BSC_API_KEY = "4C2328SSW63VFIWQZMYWD2NZERUUGN3VT1"
	BSC_WALLET_CAPITAL = "0x8A19846c7e057DBE6D77419BF1c864DAd0065d45"
	BSC_WALLET_LOWCASE = "0xFF6C223b9Dc11F247B16CED5e8a996DC4deA793E"

	PAGE_URL = "https://filfox.info/api/v1/address/"
	PAGE_URL_SUB_MINING = "/mining-stats?duration=24h"
)

type cacheB_T struct {
	CapitalB string `json:"capitalb"`
	LowcaseB string `json:"lowcaseb"`
}

type cacheFilNode_T struct {
	Balance string `json:"balance"`
	Pledge string `json:"pledge"`
	VestingFunds string `json:"vestingFunds"`
	SingleT float64 `json:"singlet"`
}

var cacheB cacheB_T
var cacheCfToF string
var cacheLoss string
var cacheDrawnFil string
var cacheFilNodes = make(map[string]cacheFilNode_T)

func listenRequests() {
	for {
		wg := &sync.WaitGroup{}
		wg.Add(5)
	//	go requestApyRate()
		go requestCfilToFil(wg)
		go requestB(wg)
		go requestLoss(wg)
		go requestDrawnFil(wg)
		go requestFilNodes(wg)
		wg.Wait()

		for _, c := range conns {
			c.cfToFCh <-cacheCfToF
			c.bCh <-cacheB
			c.lossCh <-cacheLoss
			c.drawnFilCh <-cacheDrawnFil
		//	c.apyRateCh <-cacheApyRate
			c.filNodesCh <-cacheFilNodes
		}
		time.Sleep(time.Second * TIME_CFILTOFIL)
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

// 24小时净存取
func getDrawnData() ([]float64, error) {
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

func requestCfilToFil(wg *sync.WaitGroup) {
	defer wg.Done()
	cacheCfToF = "1.2"
}

func requestB(wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(BSC_API_URL + "?module=account&action=balancemulti&address=" + BSC_WALLET_CAPITAL + "," + BSC_WALLET_LOWCASE + "&apikey=" + BSC_API_KEY)
	if err != nil {
		log.Println(err)
		return
	}
	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	capitalB := data["result"].([]interface{})[0].(map[string]interface{})["balance"].(string)
	lowcaseB := data["result"].([]interface{})[1].(map[string]interface{})["balance"].(string)

	cacheB.CapitalB = capitalB
	cacheB.LowcaseB = lowcaseB
}

func requestLoss(wg *sync.WaitGroup) {
	defer wg.Done()
	cacheLoss = "0.321"
}

func requestDrawnFil(wg *sync.WaitGroup) {
	defer wg.Done()
	cacheDrawnFil = "1000"
}

func requestFilNodes(wg *sync.WaitGroup) {
	defer wg.Done()
	filNodeKeys := []string{"f0715209"}
	for _, nodeKey := range filNodeKeys {
		resp, err := http.Get(PAGE_URL + nodeKey)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		data := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			log.Println(err)
			return
		}

		cacheFilNode := cacheFilNode_T{}
		cacheFilNode.Balance = data["miner"].(map[string]interface{})["availableBalance"].(string)
		cacheFilNode.Pledge = data["miner"].(map[string]interface{})["sectorPledgeBalance"].(string)
		cacheFilNode.VestingFunds = data["miner"].(map[string]interface{})["vestingFunds"].(string)

		active := data["miner"].(map[string]interface{})["sectors"].(map[string]interface{})["active"].(float64)

		respMining, err := http.Get(PAGE_URL + nodeKey + PAGE_URL_SUB_MINING)
		if err != nil {
			log.Println(err)
			return
		}
		defer respMining.Body.Close()

		params := make(map[string]interface{})
		err = json.NewDecoder(respMining.Body).Decode(&params)
		if err != nil {
			log.Println(err)
			return
		}

		totalRewards, err := strconv.ParseFloat(params["totalRewards"].(string)[: len(params["totalRewards"].(string)) - 14], 64)
		if err != nil {
			log.Println(err)
			return
		}
		totalRewards /= 10000

		cacheFilNode.SingleT = totalRewards / active * 16
		cacheFilNodes[nodeKey] = cacheFilNode
	}
}
