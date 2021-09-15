package main

import (
	"crypto/rand"
	"math/big"
	"net/http"
//	"io/ioutil"
	"encoding/json"
	"time"
	"fmt"
	"log"
)

const (
	TIME_CFILTOFIL = 15

	BSC_API_URL = "https://api.bscscan.com/api"
	BSC_API_KEY = "4C2328SSW63VFIWQZMYWD2NZERUUGN3VT1"
	BSC_WALLET_CAPITAL = "0x8A19846c7e057DBE6D77419BF1c864DAd0065d45"
	BSC_WALLET_LOWCASE = "0xFF6C223b9Dc11F247B16CED5e8a996DC4deA793E"
)

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

func requestCfilToFil() {
	cfilToFilCh <- "1.2"
}

func requestCapitalB() {
	for {
		resp, err := http.Get(BSC_API_URL + "?module=account&action=balance&address=" + BSC_WALLET_CAPITAL + "&apikey=" + BSC_API_KEY)
		if err != nil {
			log.Println(err)
			continue
		}

		data := make(map[string]string)

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(data)
		capitalBCh <- data["result"]
		time.Sleep(time.Second * TIME_CFILTOFIL)
	}
}

func requestLowcaseB() {
	for {
		resp, err := http.Get(BSC_API_URL + "?module=account&action=balance&address=" + BSC_WALLET_LOWCASE + "&apikey=" + BSC_API_KEY)
		if err != nil {
			continue
		}

		data := make(map[string]string)

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(data)
		lowcaseBCh <- data["result"]
		time.Sleep(time.Second * TIME_CFILTOFIL)
	}
}

func requestLoss() {
	lossCh <- "0.321"
}

func requestDrawnFil() {
	drawnFilCh <- "1000"
}
