package main

import (
)

var account string
var password string
var tokens map[string]*token_T

func init() {
	account = "aaaaaa"
	password = "999999"

	tokens = make(map[string]*token_T)
}

func main() {
	runHTTP()
}
