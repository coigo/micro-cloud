package main

import (
	"context"
	"fmt"
	// "strconv"
	"time"

	proto "github.com/coigo/image/proto/status_receiver"
	// "github.com/shirou/gopsutil/v4/cpu"
	// "github.com/shirou/gopsutil/v4/mem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main () {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("Erro ->> %v", err)
	}
	defer conn.Close()

	c := proto.NewStatusReceiverServiceClient(conn)
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// cpuUsage, err := cpu.Percent(time.Second, false)
	// if (err != nil) {
	// 	fmt.Errorf("err %v" , err)
	// }
	// memUsage, err := mem.VirtualMemory()
	

	c.ShareStatus(ctx, &proto.ImageStatus{
		MachineId: "123",
		CpuUsage: "a123",//strconv.FormatFloat(cpuUsage[0], 'f', -1, 64),
		RamUsage: "s",//strconv.FormatFloat(memUsage.UsedPercent[0], 'f', -1, 64),
		RunningImages: []*proto.RunningImage {},
	})
	
}