package main

import (
	"crypto/rand"
	"math/big"
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
