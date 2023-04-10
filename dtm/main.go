package main

import (
	"fmt"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/go-resty/resty/v2"
	"github.com/lithammer/shortuuid/v3"
)

const (
	DefaultHTTPServer = "http://localhost:36789/api/dtmsvr"
)

func main() {
	gid := shortuuid.New()
	dtmcli.XaGlobalTransaction(DefaultHTTPServer, gid, testXaFunc)
}

func testXaFunc(xa *dtmcli.Xa) (*resty.Response, error) {
	fmt.Println("hello world")

	return nil, nil

	//return nil, errors.New("handle test xa func...")
}
