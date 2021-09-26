package main

import (
	ol "log"
)

func init() {
	ol.SetFlags(ol.Lshortfile | ol.LstdFlags)

	conns = make(map[string]connection_I)
	conns2 = make(map[string]connection_I)
}

func main() {
	setHomeDir()
	setFaciDir()
	setLogDir()
//	ol.Println(logDir)
	log = newLog()
	setHttpLog()

	go listenRequests()
	go startPing(conns)
	go startPing(conns2)
	runHTTP()
}
