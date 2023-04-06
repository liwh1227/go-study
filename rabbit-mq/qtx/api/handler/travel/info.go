package travel

import (
	"gateway/lib/log"
	carbon_integral "gateway/service/carbon-integral"

	cmresp "gitee.com/liwh1227/common/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GetExchangeByType(ctx *gin.Context) {
	req, err := ctx.GetRawData()
	if err != nil {
		cmlogger.Error(err)
		cmresp.FailResponse(ctx, err)
		return
	}

	if len(req) == 0 {
		err = errors.New("request parameter is nil")
		log.RequestLogger.Error(err)
		cmresp.FailResponse(ctx, err)
		return
	}

	resp, err := carbon_integral.GetCarbonIntegralInfoByExchangeId(req)
	if err != nil {
		log.SystemLog().Error(err)
		cmresp.FailResponse(ctx, err)
		return
	}
	cmresp.SuccessResponse(ctx, resp)
}
