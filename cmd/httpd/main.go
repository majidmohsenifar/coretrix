package main

import (
	"context"
	"coretrix/internal/di"
)

func main() {
	addr := "0.0.0.0:8000"
	container := di.NewContainer()
	httpServer := container.GetHttpServer()
	indexProductCommand := container.GetIndexProductCommand()
	var flags = []string{}
	go indexProductCommand.Run(context.Background(), flags)
	err := httpServer.ListenAndServe(addr)
	if err != nil {
		panic("can not run http server because of err" + err.Error())
	}
}
