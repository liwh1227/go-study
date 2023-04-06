package router

import (
	"github.com/gin-gonic/gin"
)

type Option func(*gin.Engine)

var options = make([]Option, 0)

// Register 注册app的路由配置
func Register(opts ...Option) {
	options = append(options, opts...)
}

// Init 初始化
func Init() *gin.Engine {
	r := gin.Default()

	// 实例化http server
	for _, opt := range options {
		opt(r)
	}
	return r
}
