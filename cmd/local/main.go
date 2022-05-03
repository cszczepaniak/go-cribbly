package main

import (
	"net/http"

	"github.com/cszczepaniak/go-cribbly/cmd/common"
)

func main() {
	common.Start(http.ListenAndServe)
}
