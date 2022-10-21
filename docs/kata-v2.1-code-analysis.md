[TOC]
# Kata Containers v2.1 源码分析
[代码仓库](https://github.com/kata-containers/kata-containers)

kata 启动容器命令
```sh
/usr/local/bin/containerd-shim-kata-v2 
    -namespace k8s.io 
    -address /run/containerd/containerd.sock 
    -publish-binary /usr/bin/containerd 
    -id 41153ea62c0b0528e58eae5a345ce39e2e29c48b22a5664d2440313fde10f45c 
    -debug
```

## 分析
主函数入口
```go
func main() {
	shim.Run(types.DefaultKataRuntimeName, containerdshim.New, shimConfig)
}
```
调用 containerd/runtime/v2/shim
```go
func run(id string, initFunc Init, config Config) error {
	parseFlags() // 解析上述命令参数
	setRuntime() // 设置golang gc 和 GOMAXPROCS

	signals, err := setupSignals(config) // 设置接受的信号类型：unix.SIGTERM, unix.SIGINT, unix.SIGPIPE

	if !config.NoSubreaper { // TODO: 待查清楚，此处条件为 false，未执行
		if err := subreaper(); err != nil {
			return err
		}
	}

    // 设置事件通知
	publisher := &remoteEventsPublisher{
		address:              addressFlag,
		containerdBinaryPath: containerdBinaryFlag,
		noReaper:             config.NoReaper,
	}

    // 通过上下文传递解析的命令行参数
	ctx := namespaces.WithNamespace(context.Background(), namespaceFlag)
	ctx = context.WithValue(ctx, OptsKey{}, Opts{BundlePath: bundlePath, Debug: debugFlag})
	ctx = log.WithLogger(ctx, log.G(ctx).WithField("runtime", id))

    // 调用 containerdshim.New(),得到一个处理 GRPC 请求的 shim 进程
	service, err := initFunc(ctx, idFlag, publisher)

	switch action { // katashim 启动时，action 为 nil
	case "delete":
		logger := logrus.WithFields(logrus.Fields{
			"pid":       os.Getpid(),
			"namespace": namespaceFlag,
		})
		go handleSignals(logger, signals)
		response, err := service.Cleanup(ctx)
		if err != nil {
			return err
		}
		data, err := proto.Marshal(response)
		if err != nil {
			return err
		}
		if _, err := os.Stdout.Write(data); err != nil {
			return err
		}
		return nil
	case "start":
		address, err := service.StartShim(ctx, idFlag, containerdBinaryFlag, addressFlag)
		if err != nil {
			return err
		}
		if _, err := os.Stdout.WriteString(address); err != nil {
			return err
		}
		return nil
	default:
		if err := setLogger(ctx, idFlag); err != nil {
			return err
		}
        // 初始化 ShimClient
		client := NewShimClient(ctx, service, signals)
        // 调用 client.Serve()，启动 shim 服务
		return client.Serve()
	}
}
```
shim.Run() 执行完成后得到 client 结构为
```go
&Client{
		service: svc,
		context: ctx,
		signals: signals,
	}
```
其中 service 接口对应的结构体实现为 containerdshim.service 为
```go
// containerdshim.service
	s := &service{
		id:         id,
		pid:        uint32(os.Getpid()), // shim pid
		ctx:        ctx,
		containers: make(map[string]*container),    // 容器
		events:     make(chan interface{}, chSize), // 事件通道
		ec:         make(chan exit, bufferSize),    // 退出的容器进程
		cancel:     cancel,
	}
```
signals 接受的信号类型：unix.SIGTERM, unix.SIGINT, unix.SIGPIPE

接下来查看 client.Serve() 实现
```go
// 启动 shim 服务
func (s *Client) Serve() error {
	dump := make(chan os.Signal, 32)
	setupDumpStacks(dump) // 设置可接收信号 syscall.SIGUSR1

	path, err := os.Getwd() // 获取当前执行路径
	if err != nil {
		return err
	}
	server, err := newServer() // 新建 ttrpc server
	if err != nil {
		return errors.Wrap(err, "failed creating server")
	}

	logrus.Debug("registering ttrpc server")
    // 注册 ttrpc server
    // containerdshim.service 实现了接口 TaskService，所以可以注册到 ttrpc server
	shimapi.RegisterTaskService(server, s.service)

    // socketFlag 为空，服务将监听在 abstract socket 上
	if err := serve(s.context, server, socketFlag); err != nil {
		return err
	}
	logger := logrus.WithFields(logrus.Fields{
		"pid":       os.Getpid(),
		"path":      path,
		"namespace": namespaceFlag,
	})
    // 发送 SIGUSR1 可打印栈信息
	go func() {
		for range dump {
			dumpStacks(logger)
		}
	}()
    // 处理 client 信号
	return handleSignals(logger, s.signals)
}
```

## ttrpc 介绍