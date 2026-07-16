package commands

import (
	"context"
	"fmt"

	"github.com/coigo/image/docker"
	proto "github.com/coigo/image/proto/command_dispatcher"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedCommandDispatcherServiceServer
}

func (s *server) CreateCommand(ctx context.Context, req *proto.CreateCommandRequest) (*proto.CommandResponse, error) {
	fmt.Println("Criando container")

	contaienrId, err := docker.CreateContainer(ctx)
	if err != nil {
		fmt.Println("Erro criando o container", err)
		return nil, err
	}
	response := &proto.CommandResponse {
		ContainerId: *contaienrId,
	}
	return response, err
}

func (s *server) DownCommand(ctx context.Context, req *proto.DownCommandRequest) (*proto.CommandResponse, error) {
	err := docker.DownContainer(ctx, req.ContainerId)
	if err != nil  {
		fmt.Println("Erro finalizando o container:", err)
		return nil, err
	}
	response := &proto.CommandResponse{
		ContainerId: req.ContainerId,
	}
	return response, err
}

func NewServer (ctx context.Context) *grpc.Server{
	grpcServer := grpc.NewServer()
	proto.RegisterCommandDispatcherServiceServer(grpcServer, &server{})
	return grpcServer
}