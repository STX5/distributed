package main

import (
	"context"
	"distributed/registry"
	"fmt"
	"log"
	"net/http"
)

func main() {
	registry.SetupRegistryService()
	http.Handle("/services", &registry.RegistryService{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = registry.ServerPort

	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			log.Println(srv.ListenAndServe()) //出错才会返回，然后执行cancel()
			cancel()
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("Registry service started, Press any key to exit.\n")
			var s string
			fmt.Scanln(&s)
			srv.Shutdown(ctx)
			cancel()
		}
	}()
	<-ctx.Done() // 接收到cancel()
	fmt.Println("Shutting down registry service")
}
