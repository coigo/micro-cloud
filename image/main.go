package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	"time"

	"github.com/coigo/image/commands"
	"github.com/coigo/image/docker"
	proto "github.com/coigo/image/proto/status_receiver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("controll:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("Erro ->> %v", err)
	}
	defer conn.Close()

	c := proto.NewStatusReceiverServiceClient(conn)

	docker.NewDockerClient()
	defer docker.CloseDockerClient()

	

	var wp sync.WaitGroup
	wp.Add(1)
	go func () {
	
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
	}()
	wp.Wait()


	ctx := context.Background()

	lis, err := net.Listen("tcp", "50051")
	if err != nil {
		fmt.Errorf("Erro ouvindo: %v", err)
	}
	server := commands.NewServer(ctx)
	if err := server.Serve(lis); err != nil {
		fmt.Errorf("Erro servindo: %v", err)
	}
	
	
}
