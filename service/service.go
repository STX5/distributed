package service

import (
	"context"
	"distributed/registry"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, host, port string, reg registry.Registration,
	registerHandersFunc func()) (context.Context, error) {
	registerHandersFunc() // 注册handler

	ctx = startService(ctx, reg.ServiceName, host, port) // 启动后，进行注册
	err := registry.RegisterService(reg)                 //注册服务
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName,
	host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	var srv http.Server
	srv.Addr = ":" + port

	go func() {
		log.Println(srv.ListenAndServe()) //出错才会返回，返回error，然后执行下一行的cancel()
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err) //不需要return，只需要记录。因为需要继续执行cancel()
		}
		cancel()
	}()

	go func() {
		fmt.Printf("%v started, Press any key to exit. \n", serviceName)
		var s string
		fmt.Scanln(&s)
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err) //不需要return，只需要记录。因为需要继续执行cancel()
		}
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
