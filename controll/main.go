package main

import (
	// "fmt"
	"context"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	cd "github.com/coigo/micro-cloud/proto/command_dispatcher"
	"github.com/coigo/micro-cloud/infra"
	"github.com/coigo/micro-cloud/statusreciever"
	// "google.golang.org/grpc"
)

func main () {

	ctx := context.Background()
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

	conn, err := grpc.NewClient("dind1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("Erro com o client:", err)
	}

	c := cd.NewCommandDispatcherServiceClient(conn)
	c.CreateCommand(ctx, &cd.CreateCommandRequest{
		ContainerSize: cd.ContainerSize_SMALL,
	})
	
	// dockerId := commandservice.UpCommand()
	fmt.Printf("Container %v criado.\n")
	// commandservice.DownCommand(dockerId)
}