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
	readConf()
	ol.Println(config)

	log = newLog()
	setHttpLog()


	go listenRequests()
	go startPing(conns)
	go startPing(conns2)
	runHTTP()
}
