# panic问题复盘

## Q1: 并发初始化sdk client导致panic

### 背景

该问题发生在某次生产环境中更新服务导致，只有服务启动时发生了panic，后面再没有出现过，由于查看服务的容器日志，将发生panic的地方进行了截取如下：

```log
{"log":"fatal error: concurrent map writes\n","stream":"stderr","time":"2023-03-16T08:56:45.00349472Z"}
{"log":"\n","stream":"stderr","time":"2023-03-16T08:56:45.006692321Z"}
{"log":"goroutine 302 [running]:\n","stream":"stderr","time":"2023-03-16T08:56:45.006706629Z"}
{"log":"chainmaker.org/chainmaker/common/v2/log.newLogger(0xc0001698e0, {0x203000?})\n","stream":"stderr","time":"2023-03-16T08:56:45.006713412Z"}
{"log":"\u0009/go/pkg/mod/chainmaker.org/chainmaker/common/v2@v2.2.1/log/log.go:159 +0x1b3\n","stream":"stderr","time":"2023-03-16T08:56:45.006717114Z"}
{"log":"chainmaker.org/chainmaker/common/v2/log.InitSugarLogger(0xc0001698e0)\n","stream":"stderr","time":"2023-03-16T08:56:45.006721139Z"}
{"log":"\u0009/go/pkg/mod/chainmaker.org/chainmaker/common/v2@v2.2.1/log/log.go:141 +0x8b\n","stream":"stderr","time":"2023-03-16T08:56:45.006752304Z"}

// 服务代码逻辑
{"log":"gateway/lib/chainmaker/client.getDefaultLogger()\n","stream":"stderr","time":"2023-03-16T08:56:45.006765689Z"}
{"log":"\u0009/gateway/lib/chainmaker/client/sdk.go:42 +0x11d\n","stream":"stderr","time":"2023-03-16T08:56:45.006771836Z"}
{"log":"gateway/lib/chainmaker/client.getSdkClient({0x13076c9, 0x6}, {0xc00060e780, 0xe0ea79?})\n","stream":"stderr","time":"2023-03-16T08:56:45.006785395Z"}
{"log":"\u0009/gateway/lib/chainmaker/client/sdk.go:82 +0xacb\n","stream":"stderr","time":"2023-03-16T08:56:45.006801366Z"}
{"log":"gateway/lib/chainmaker/client.QueryTxInfosByTxIds({0x13076c9, 0x6}, {0xc00060e780?, 0x0?}, {0xc000356280, 0x6, 0x0?})\n","stream":"stderr","time":"2023-03-16T08:56:45.006835003Z"}
{"log":"\u0009/gateway/lib/chainmaker/client/sdk.go:230 +0x88\n","stream":"stderr","time":"2023-03-16T08:56:45.00685314Z"}
{"log":"gateway/service/carbon-integral.HandleExchangeInfoUpChainTransaction()\n","stream":"stderr","time":"2023-03-16T08:56:45.006865409Z"}
{"log":"\u0009/gateway/service/carbon-integral/exchange_info.go:241 +0xea\n","stream":"stderr","time":"2023-03-16T08:56:45.006870451Z"}
{"log":"gateway/service/monitor.HandleExchangeWaitingResultTransaction.func1()\n","stream":"stderr","time":"2023-03-16T08:56:45.006890629Z"}
{"log":"\u0009/gateway/service/monitor/exchange.go:24 +0x34\n","stream":"stderr","time":"2023-03-16T08:56:45.006899764Z"}

// 系统方法
{"log":"github.com/robfig/cron/v3.FuncJob.Run(0xc0004977d0?)\n","stream":"stderr","time":"2023-03-16T08:56:45.006905911Z"}
{"log":"\u0009/go/pkg/mod/github.com/robfig/cron/v3@v3.0.1/cron.go:136 +0x1a\n","stream":"stderr","time":"2023-03-16T08:56:45.006917682Z"}
{"log":"github.com/robfig/cron/v3.(*Cron).startJob.func1()\n","stream":"stderr","time":"2023-03-16T08:56:45.006930399Z"}
{"log":"\u0009/go/pkg/mod/github.com/robfig/cron/v3@v3.0.1/cron.go:312 +0x6a\n","stream":"stderr","time":"2023-03-16T08:56:45.006935191Z"}
{"log":"created by github.com/robfig/cron/v3.(*Cron).startJob\n","stream":"stderr","time":"2023-03-16T08:56:45.006949821Z"}
{"log":"\u0009/go/pkg/mod/github.com/robfig/cron/v3@v3.0.1/cron.go:310 +0xad\n","stream":"stderr","time":"2023-03-16T08:56:45.006973944Z"}
```

产生的根本原因是：并发写map导致。根据panic日志，发现是在 `newLogger` 过程中出现的，该方法是 `chainmaker` 的官方包中的方法，详细代码如下：

chainmaker.org/chainmaker/common/v2@v2.2.1/log/log.go:159 

```go
func newLogger(logConfig *LogConfig, level zap.AtomicLevel) *zap.Logger {
	var (
		hook io.Writer
		ok   bool
		err  error
	)

	_, ok = hookMap[logConfig.LogPath]
	if !ok {
		hook, err = getHook(logConfig.LogPath, logConfig.MaxAge, logConfig.RotationTime, logConfig.RotationSize)
		if err != nil {
			log.Fatalf("new logger get hook failed, %s", err)
		}
		// line 159
		hookMap[logConfig.LogPath] = struct{}{}
	} else {
		hook, err = getHook(logConfig.LogPath, logConfig.MaxAge, 0, logConfig.RotationSize)
		if err != nil {
			log.Fatalf("new logger get hook failed, %s", err)
		}
	}
}
```

其实同go的报错信息一致，就是存在并发调用map的情况，这是根本原因。经过排查，发现业务系统中调用这部分逻辑确实存在并发的情况。业务层中，存在多个定时任务进行逻辑处理的地方，定时任务我们使用的是 `github.com/robfig/cron/v3` 这个包，这里会存在以下两种情况引起并发写：

1. 我们为了并行开启监听任务，使用了go协程分别启动不同的监听任务；
2. cron包中进行每次进行start job时都会启用一个新的协程来进行；

定时任务开启后，会存在 `newSdkClient` 的情况，该方法会调用 `newLogger` ，这就会导致多协程并发写map，从而引发系统的panic。

为了解决上述问题，其实需要在初始化sdk client时加上锁，防止多个协程同时进行sdk client的初始化。