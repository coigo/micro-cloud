package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/coigo/micro-cloud/controller"
	"github.com/coigo/micro-cloud/infra"
	"github.com/coigo/micro-cloud/statusreciever"
	
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

	for {
		controller.CreateContainer(ctx)
		time.Sleep(10 * time.Second)
	}

}

// CALCULAR O PROCESSAMENTO DISPONIVEL
// LIMITAR OS RESULTADOS DE CONTAINERS POR PROCESSAMENTO DISPONIVEL
