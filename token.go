package main

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

const overtime = 900

type token_T struct {
	hash [16]byte
	account string
	timestamp int64
}

func validation(account, key string) bool {
	for k, t := range tokens {
		if k == key {
			if account != t.account {
				return false
			}

			if time.Now().Unix() - t.timestamp > 900 {
				delete(tokens, key)
				return false
			}

			data := []byte(t.account)
			data = append(data, uint64ToBytes(uint64(t.timestamp))...)
			hash := md5.Sum(data)

			if hex.EncodeToString(t.hash[:]) != hex.EncodeToString(hash[:]) {
				return false
			}

			return true
		}
	}

	return false
}
