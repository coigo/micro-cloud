package main

import (
	"fmt"
	"github.com/coigo/micro-cloud/commandservice"
	proto "github.com/coigo/micro-cloud/controll/proto/status_receiver"
	"google.golang.org/grpc"
)

type server struct {
	proto.s *server
}

func main () {
	// dockerId := commandservice.UpCommand()
	// fmt.Printf("Container %v criado.\n", dockerId)
	// commandservice.DownCommand(dockerId)
}