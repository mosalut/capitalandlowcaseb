package main

import (
	"errors"
)

func bytesToUint64(b []byte) (uint64, error) {
	if len(b) != 8 {
		return 0, errors.New("bytes length must be 8")
	}

	return uint64(uint8(b[0])) << 56 | uint64(uint8(b[1])) << 48 | uint64(uint8(b[2])) << 40 | uint64(uint8(b[3])) << 32 | uint64(uint8(b[4])) << 24 | uint64(uint8(b[5])) << 16 | uint64(uint8(b[6])) << 8 | uint64(uint8(b[7])), nil

}

func uint64ToBytes(i uint64) []byte {
	return []byte{uint8(i >> 56), uint8(i << 8 >> 56), uint8(i << 16 >> 56), uint8(i << 24 >> 56), uint8(i << 32 >> 56), uint8(i << 40 >> 56), uint8(i << 48 >> 56), uint8(i)}
}
