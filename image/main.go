package main

import (
	"context"
	"fmt"
	"os"

	"time"

	proto "github.com/coigo/image/proto/status_receiver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/moby/moby/client"
)

func main() {

	ctx := context.Background()
	
	conn, err := grpc.NewClient("controll:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("Erro ->> %v", err)
	}
	defer conn.Close()


	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("containers", containers)
	c := proto.NewStatusReceiverServiceClient(conn)

	for {
		fmt.Println("weee")

		hostname, err := os.Hostname()
		if err != nil {
			fmt.Errorf("host err %v" , err)
			
		}
		
		memMax, _ := os.ReadFile("/sys/fs/cgroup/memory.max")
		cpuMax, _ := os.ReadFile("/sys/fs/cgroup/cpu.max")
		
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		
		fmt.Println("Fazendo requisição")
		c.ShareStatus(ctx, &proto.ImageStatus{
			MachineId:     hostname,
			CpuUsage:      "",
			RamUsage:      "",
			RunningImages: []*proto.RunningImage{},
			CpuTotal:      string(cpuMax),
			RamTotal:      string(memMax),
		})
		
		time.Sleep(5 * time.Second)
	}

}
