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

var conns map[string]*conn_T

func validation(account, key, networking string) bool {
	for k, c := range conns {
		if k == key {
			if account != c.account {
				return false
			}

			if networking != c.networking {
				return false
			}

			now := time.Now().Unix()
			fmt.Println(now, c.timestamp)
			if now - c.timestamp > 900 || k != hex.EncodeToString(c.hash[:]) {
				disconnect(k)
				return false
			}

			data := []byte(c.account)
			data = append(data, uint64ToBytes(uint64(c.timestamp))...)
			hash := md5.Sum(data)

			if hex.EncodeToString(c.hash[:]) != hex.EncodeToString(hash[:]) {
				return false
			}

			return true
		}
	}

	return false
}
