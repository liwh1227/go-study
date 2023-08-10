package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	cmresp "github.com/liwh1227/go-common/response"
	"github.com/pkg/errors"
	"qtx/rabbitmq"
)

// 用户收集气泡
func UpdateBubbleStatus(ctx *gin.Context) {
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

	err = rabbitmq.NewChannel().Publish(rabbitmq.BubbleExchange, rabbitmq.BubbleKey, req)
	if err != nil {
		fmt.Println(err)
		cmresp.FailResponse(ctx, err)
		return
	}

	cmresp.SuccessResponse(ctx, nil)
}

// 更新用户的兑换信息状态
func UpdateExchangeStatus(ctx *gin.Context) {
	req, err := ctx.GetRawData()
	if err != nil {
		fmt.Println(err)
		cmresp.FailResponse(ctx, err)
		return
	}
	fmt.Printf("%#v\n", req)

	err = rabbitmq.NewChannel().Publish(rabbitmq.ExchangeExchange, rabbitmq.ExchangeKey, req)
	if err != nil {
		fmt.Println(err)
		cmresp.FailResponse(ctx, err)
		return
	}

	cmresp.SuccessResponse(ctx, nil)
}
