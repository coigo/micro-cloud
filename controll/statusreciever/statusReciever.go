package statusreciever

import (
	"encoding/json"
	"time"

	"context"
	"fmt"

	"github.com/coigo/micro-cloud/infra"
	proto "github.com/coigo/micro-cloud/proto/status_receiver"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	proto.UnimplementedStatusReceiverServiceServer
}


func (s *server) ShareStatus (ctx context.Context, imageStatus *proto.ImageStatus ) (*emptypb.Empty, error)  {

	data, err := json.Marshal(imageStatus)
	if err != nil {
		fmt.Errorf("Erro ao parsear a resposta:", err)
	}

	infra.Redis.Set(ctx, "machine-status:"+imageStatus.MachineId, string(data), 0)
		iter := infra.Redis.Scan(ctx, 0,"machine-status:*", 0).Iterator()
		
		for iter.Next(ctx) {
	    key := iter.Val()
	
	    value, err := infra.Redis.Get(ctx, key).Result()
	    if err != nil {
	        continue
	    }
	
	    fmt.Println(key, value)

	}
	
	
	fmt.Println(time.Now().Unix(), " | Nova requisição ", imageStatus)
	return &emptypb.Empty{}, nil
}

func NewServer () *grpc.Server{
	grpcServer := grpc.NewServer()
	proto.RegisterStatusReceiverServiceServer(grpcServer, &server{})
	return grpcServer
}