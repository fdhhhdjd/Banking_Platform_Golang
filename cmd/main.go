package main

import (
	"sync"

	"github.com/fdhhhdjd/Banking_Platform_Golang/server"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		server.Server()
	}()

	go func() {
		defer wg.Done()
		server.StartGateWayGRPCServer()
	}()

	go func() {
		defer wg.Done()
		server.StartGRPCServer()
	}()

	wg.Wait()

}
