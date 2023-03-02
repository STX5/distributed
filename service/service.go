package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, serviceName, host, port string,
	registerHandersFunc func()) (context.Context, error) {
	registerHandersFunc() // 注册handler
	ctx = startService(ctx, serviceName, host, port)

	return ctx, nil
}

func startService(ctx context.Context, serviceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	var srv http.Server
	srv.Addr = ":" + port

	go func() {
		log.Println(srv.ListenAndServe()) //出错才会返回，返回error，然后执行下一行的cancel()
		cancel()
	}()

	go func() {
		fmt.Printf("%v started, Press any key to exit. \n", serviceName)
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
