package commands

import (
	"context"

	proto "github.com/coigo/image/proto/command_dispatcher"
)

type server struct {
	proto.UnimplementedCommandDispatcherServiceServer
}

func (s *server) CreateCommand(context.Context, *proto.CreateCommandRequest) (*proto.CommandResponse, error) {
	
}
func (s *server) DownCommand(context.Context, *proto.DownCommandRequest) (*proto.CommandResponse, error) {

}

func NewServer (ctx context.Context) {

}