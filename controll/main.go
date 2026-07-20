package main

import (
	// "fmt"
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/coigo/micro-cloud/infra"
	cd "github.com/coigo/micro-cloud/proto/command_dispatcher"
	"github.com/coigo/micro-cloud/statusreciever"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	

	conn, err := grpc.NewClient("dind1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("Erro com o client:", err)
	}

	c := cd.NewCommandDispatcherServiceClient(conn)
	resp, err := c.CreateCommand(ctx, &cd.CreateCommandRequest{
		ContainerSize: cd.ContainerSize_SMALL,
	})
	
	if err != nil {
		fmt.Errorf("Erro na req:", err)
	}
	
	// dockerId := commandservice.UpCommand()
	fmt.Printf("Container %v criado.\n", resp.ContainerId)
	wg.Wait()
	// commandservice.DownCommand(dockerId)
}

// CALCULAR O PROCESSAMENTO DISPONIVEL
// LIMITAR OS RESULTADOS DE CONTAINERS POR PROCESSAMENTO DISPONIVEL
