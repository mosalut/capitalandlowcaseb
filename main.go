package main

import (
	ol "log"
)

func init() {
	ol.SetFlags(ol.Lshortfile | ol.LstdFlags)

	setHomeDir()
	setFaciDir()
	setLogDir()
	readConf()

	log = newLog()
	setHttpLog()

	initDB()

	signIns[0] = signInDev
	signIns[1] = signInRelease
	signIn = signIns[config.sms]

	conns = make(map[string]connection_I)
	conns2 = make(map[string]connection_I)
}

func main() {
	go listenRequests()
	go startPing(conns)
	go startPing(conns2)
	runHTTP()
}
