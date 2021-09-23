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

	conns = make(map[string]*conn_T)
	conns2 = make(map[string]*conn2_T)
}

func main() {
	go listenRequests()
	runHTTP()
}
