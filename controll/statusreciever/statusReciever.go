package statusreciever

import (
	"encoding/json"
	"strconv"
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


func (s *server) ShareStatus(ctx context.Context, imageStatus *proto.ImageStatus) (*emptypb.Empty, error) {
	for _, container := range imageStatus.RunningImages {
		imageStatus.CpuUsage = SumUsage(imageStatus.CpuUsage, container.CpuTotal)
		imageStatus.RamUsage += container.RamTotal
	}

	data, err := json.Marshal(imageStatus)
	if err != nil {
		fmt.Println("Erro ao parsear a resposta:", err)
		return nil, err
	}

	if err := infra.Redis.Set(ctx, "machine-status:"+imageStatus.MachineId, string(data), 0).Err(); err != nil {
		fmt.Println("Erro ao salvar no Redis:", err)
		return nil, err
	}

	fmt.Println(time.Now().Unix(), " | Nova requisição ", imageStatus)
	return &emptypb.Empty{}, nil
}
func NewServer () *grpc.Server{
	grpcServer := grpc.NewServer()
	proto.RegisterStatusReceiverServiceServer(grpcServer, &server{})
	return grpcServer
}

func SumUsage(machine string, container string) string {

	if machine == "" {
		machine = "0"
	}
	if container == "" {
		container = "0"
	}
	
	machineUsage, err := strconv.ParseFloat(machine, 64)
	if err != nil {
		fmt.Println("Erro parseando dados da máquina:", err)
	}

	containerUsage, err := strconv.ParseFloat(container, 64)
	if err != nil {
		fmt.Println("Erro parseando dados do container:", err)
	}

	return fmt.Sprintf("%.2f", machineUsage+containerUsage)
}