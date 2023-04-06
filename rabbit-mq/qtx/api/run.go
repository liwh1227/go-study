package api

import (
	"context"
	"fmt"
	"net/http"
	"qtx/api/router"
	"qtx/common"
	"time"
)

// Run 启动api服务
func Run(ctx context.Context, errChan chan error) {
	// panic 自动恢复
	defer func() {
		if err := recover(); err != nil {
			common.Log.Error(err)
			return
		}
	}()

	start := time.Now()

	engine := router.Init()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "loclhost", 10000),
		Handler: engine,
	}

	common.Log.Info("http: successfully initialized %v", time.Since(start))

	// 启动http server
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				errChan <- err
				common.Log.Error("http: web server shutdown complete")
			} else {
				fmt.Println("http: web server closed unexpect: ", err)
			}
		}
	}()

	// 关闭http server
	<-ctx.Done()
	err := server.Close()
	if err != nil {
		fmt.Println("http: web server shutdown failed: ", err)
	}
}
