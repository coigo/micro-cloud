package main

import (
	"context"
	"fmt"
	"os"

	"time"

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

	for {
		fmt.Println("weee")

		hostname, err := os.Hostname()
		if err != nil {
			fmt.Errorf("host err %v" , err)
			
		}
		
		memCurr, _ := os.ReadFile("/sys/fs/cgroup/memory.current")
		memMax, _ := os.ReadFile("/sys/fs/cgroup/memory.max")
		cpuMax, _ := os.ReadFile("/sys/fs/cgroup/cpu.max")
		cpuCurr, _ := os.ReadFile("/sys/fs/cgroup/cpu.stat")
		
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		
		fmt.Println("Fazendo requisição")
		c.ShareStatus(ctx, &proto.ImageStatus{
			MachineId:     hostname,
			CpuUsage:      string(cpuCurr),
			RamUsage:      string(memCurr),
			RunningImages: []*proto.RunningImage{},
			CpuTotal:      string(cpuMax),
			RamTotal:      string(memMax),
		})
		
		time.Sleep(5 * time.Second)
	}

}
