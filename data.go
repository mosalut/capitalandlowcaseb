package main

import (
//	"crypto/rand"
)

// 最新流通性
/*
func getCirulationData() ([]float64, error) {
	balances make([]float64, 24, 24)

	for i, balances := range balances {
		max := big.NewInt(65536)
		intergerI, err := rand.Int(rand.Reader, max)
		if err != nil {
			return 0, err
		}
		decimalI, err := rand.Int(rand.Reader, max)
		if err != nil {
			return 0, err
		}

		intergerF := big.NewFloat(intergerI.String())
		decimalF := big.NewFloat(decimalI.String())

		balances[i] = intergerF.Float64() + decimalF.Float64() / 100000
	}

	fmt.Println(balances)

	return balances, nil
}
*/
