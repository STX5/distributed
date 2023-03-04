package main

import (
	"context"
	"distributed/grades"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddr := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registration{
		ServiceName: registry.GradingService,
		ServiceURL:  serviceAddr,
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		grades.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting Down Log Service...")
}
