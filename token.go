package main

import (
	"crypto/md5"
	"encoding/hex"
	"time"
	"fmt"
)

const overtime = 900

type token_T struct {
	hash [16]byte
	account string
	timestamp int64
	networking string
}

var conns map[string]connection_I

func validation(account, key, networking string) bool {
	for k, c := range conns {
		conn := c.(*conn_T)
		if k == key {
			if account != conn.account {
				return false
			}

			if networking != conn.networking {
				return false
			}

			now := time.Now().Unix()
			fmt.Println(now, conn.timestamp)
			if now - conn.timestamp > 900 || k != hex.EncodeToString(conn.hash[:]) {
				c.disconnect(k)
				return false
			}

			data := []byte(conn.account)
			data = append(data, uint64ToBytes(uint64(conn.timestamp))...)
			hash := md5.Sum(data)

			if hex.EncodeToString(conn.hash[:]) != hex.EncodeToString(hash[:]) {
				return false
			}

			return true
		}
	}

	return false
}
