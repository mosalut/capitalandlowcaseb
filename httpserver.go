package main

import (
	"github.com/gin-gonic/gin"
)

func runHTTP() {
	r := gin.Default()
	r.Run(":8888")
}
