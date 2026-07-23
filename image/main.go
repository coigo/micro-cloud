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

	ctx := context.Background()

	conn, err := grpc.NewClient("controll:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Erro ->> %v", err)
	}
	defer conn.Close()

	c := proto.NewStatusReceiverServiceClient(conn)

	docker.NewDockerClient()
	defer docker.CloseDockerClient()

	
	var wp sync.WaitGroup
	wp.Add(1)

	go func () {
		defer wp.Done()
		for {
			ControllerRequest(ctx, c)
			time.Sleep(time.Second * 5)
		}
	}()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
	    fmt.Println("Erro ouvindo: %v", err)
	}

	server := commands.NewServer(ctx)
	if err := server.Serve(lis); err != nil {
		fmt.Println("Erro servindo: %v", err)
	}
	wp.Wait()
	
	
}

func ControllerRequest (ctx context.Context, c proto.StatusReceiverServiceClient) {
			address:= os.Getenv("CONTAINER_ADDRESS")
			
			memMax, _ := os.ReadFile("/sys/fs/cgroup/memory.max")
			cpuMax, _ := os.ReadFile("/sys/fs/cgroup/cpu.max")
			containers, err := docker.ListContainers(ctx)
			if err != nil {
				fmt.Println("Erro buscando containers", err)
			}

			runningContainers := []*proto.RunningImage{} 
			for _, container := range(containers) {
			    runningContainers = append(runningContainers, &proto.RunningImage{
			        MachineId: container.ID,
			        CpuUsage:  "",
			        RamUsage:  "",
			        CpuTotal:  fmt.Sprintf("%d", container.NanoCPUs),
			        RamTotal:  fmt.Sprintf("%d", container.MemoryMax),
			    })
			}
			
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			
			fmt.Println("Fazendo requisição")
			c.ShareStatus(ctx, &proto.ImageStatus{
				MachineId:     address,
				CpuUsage:      "",
				RamUsage:      "",
				RunningImages: runningContainers,
				CpuTotal:      string(cpuMax),
				RamTotal:      string(memMax),
			})
}