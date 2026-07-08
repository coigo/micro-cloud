package main

import (
	// "fmt"
	// "github.com/coigo/micro-cloud/commandservice"
	"context"
	"fmt"
	"net"
	"sync"

	proto "github.com/coigo/micro-cloud/proto/status_receiver"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	// "google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedStatusReceiverServiceServer
}

func (s *server) ShareStatus (ctx context.Context, imageStatus *proto.ImageStatus ) (*emptypb.Empty, error)  {
	fmt.Println("aqui ->>", imageStatus)
	return &emptypb.Empty{}, nil
}

func main () {
	fmt.Println("funcionando?")

	lis, err := net.Listen("tcp",":50051")
	if (err != nil) {
		fmt.Errorf("Erro ouvindo a porta 50051")
	}
	fmt.Println("funcionando?")

	grpcServer := grpc.NewServer()
	proto.RegisterStatusReceiverServiceServer(grpcServer, &server{})
	fmt.Println("funcionando?")

	var wg sync.WaitGroup

	wg.Add(1)

	
	go func () {
		defer wg.Done()
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Errorf("Erro ->>.", err)
		}
	}()
	
	
	fmt.Println("funcionando?")
	wg.Wait()
	// dockerId := commandservice.UpCommand()
	// fmt.Printf("Container %v criado.\n", dockerId)
	// commandservice.DownCommand(dockerId)
}