package wallet

import (
	"github.com/gin-gonic/gin"
	cmlog "github.com/liwh1227/go-common/logger"
	cmresp "github.com/liwh1227/go-common/response"
	"github.com/pkg/errors"
)

var (
	cmlogger = cmlog.GetLoggerByService("API", "gateway")

	ErrRequestParamIsNil = errors.New("request parameter is nil")
)

// 注册钱包
func RegisterWallet(ctx *gin.Context) {
	req, err := ctx.GetRawData()
	if err != nil {
		cmlogger.Error(err)
		cmresp.FailResponse(ctx, err)
		return
	}

	if len(req) == 0 {
		cmlogger.Error(ErrRequestParamIsNil)
		cmresp.FailResponse(ctx, err)
		return
	}

	// add method
	//resp, err := wallet.Register(req)
	//if err != nil {
	//	cmlogger.Error(err)
	//	cmresp.FailResponse(ctx, err)
	//	return
	//}

	cmresp.SuccessResponse(ctx, nil)
}
