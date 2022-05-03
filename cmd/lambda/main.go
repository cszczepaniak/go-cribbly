package main

import (
	"github.com/apex/gateway"

	"github.com/cszczepaniak/go-cribbly/cmd/common"
)

func main() {
	common.Start(gateway.ListenAndServe)
}
