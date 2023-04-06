package travel

import (
	"fmt"
	cmresp "gitee.com/liwh1227/common/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Metro 地铁出行数据
func Metro(ctx *gin.Context) {
	req, err := ctx.GetRawData()
	if err != nil {
		fmt.Println(err)

		return
	}

	if len(req) == 0 {
		err = errors.New("request parameter is nil")
		fmt.Println(err)
		cmresp.FailResponse(ctx, err)
		return
	}
	//
	//err = carbon_integral.HandleMetroTravelInfo(req)
	//if err != nil {
	//	cmlogger.Error(err)
	//	cmresp.FailResponse(ctx, err)
	//	return
	//}

	cmresp.SuccessResponse(ctx, nil)
}

func InitAdminUserWallet(ctx *gin.Context) {
	_, err := ctx.GetRawData()
	if err != nil {
		fmt.Println(err)
		cmresp.FailResponse(ctx, err)
		return
	}

	//err = carbon_integral.InitAdminUserWallet(req)
	//if err != nil {
	//	log.SystemLog().Error(err)
	//	cmresp.FailResponse(ctx, err)
	//	return
	//}

	cmresp.SuccessResponse(ctx, nil)
}
