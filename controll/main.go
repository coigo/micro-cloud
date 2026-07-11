package main

import (
	// "fmt"
	"github.com/coigo/micro-cloud/commandservice"
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/coigo/micro-cloud/infra"
	"github.com/coigo/micro-cloud/statusreciever"
	// "google.golang.org/grpc"
)

func main () {

	ctx := context.Background()
	commandservice.UpCommand()
	infra.NewRedisConn(ctx)

	lis, err := net.Listen("tcp",":50051")
	if (err != nil) {
		fmt.Errorf("Erro ouvindo a porta 50051")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	
	go func () {
		defer wg.Done()
		server := statusreciever.NewServer()
		if err := server.Serve(lis); err != nil {
			fmt.Errorf("Erro ->>.", err)
		}
	}()
	
	wg.Wait()
	
	// dockerId := commandservice.UpCommand()
	fmt.Printf("Container %v criado.\n")
	// commandservice.DownCommand(dockerId)
}