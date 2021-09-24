package main

import (
	"log"
)

var account string
var password string

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	account = "aaaaaa"
	password = "999999"

	conns = make(map[string]connection_I)
	conns2 = make(map[string]connection_I)
}

func main() {
	go listenRequests()
	go startPing(conns)
	go startPing(conns2)
	runHTTP()
}
