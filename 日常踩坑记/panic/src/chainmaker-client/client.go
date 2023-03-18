package chainmaker_client

import (
	"chainmaker.org/chainmaker/common/v2/log"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"go.uber.org/zap"
	"sync"
	"time"
)

// 获取默认的日志级别
func getDefaultLogger() *zap.SugaredLogger {
	logConfig := log.LogConfig{
		Module:       "[SDK]",
		LogPath:      "./sdk.log",
		LogLevel:     log.GetLogLevel("WARN"),
		MaxAge:       10,
		JsonFormat:   false,
		ShowLine:     true,
		LogInConsole: false,
	}

	logger, _ := log.InitSugarLogger(&logConfig)
	return logger
}

// 这种方式会导致panic
// fatal error: concurrent map writes
func NewClient() (*sdk.ChainClient, error) {
	time.Sleep(1 * time.Second)
	return sdk.NewChainClient(
		sdk.WithConfPath("./chainmaker-client/sdk.yml"),
		sdk.WithChainClientLogger(getDefaultLogger()),
	)
}

var (
	mutex sync.Mutex
)

// 加锁防止并发写
func NewClientWithMutex() (*sdk.ChainClient, error) {
	time.Sleep(1500 * time.Millisecond)
	mutex.Lock()
	defer mutex.Unlock()
	return sdk.NewChainClient(
		sdk.WithConfPath("./chainmaker-client/sdk.yml"),
		sdk.WithChainClientLogger(getDefaultLogger()),
	)
}
