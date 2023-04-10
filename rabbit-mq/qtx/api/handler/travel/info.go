package travel

import (
	"github.com/gin-gonic/gin"
	cmresp "github.com/liwh1227/go-common/response"
	"github.com/pkg/errors"
	"qtx/common"
)

func GetExchangeByType(ctx *gin.Context) {
	req, err := ctx.GetRawData()
	if err != nil {
		common.Log.Error(err)
		cmresp.FailResponse(ctx, err)
		return
	}

	if len(req) == 0 {
		err = errors.New("request parameter is nil")
		common.Log.Error(err)
		cmresp.FailResponse(ctx, err)
		return
	}

	//resp, err := carbon_integral.GetCarbonIntegralInfoByExchangeId(req)
	//if err != nil {
	//	common.Log.Error(err)
	//	cmresp.FailResponse(ctx, err)
	//	return
	//}
	cmresp.SuccessResponse(ctx, nil)
}
