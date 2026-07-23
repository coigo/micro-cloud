package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/coigo/micro-cloud/infra"
	cd "github.com/coigo/micro-cloud/proto/command_dispatcher"
	sr "github.com/coigo/micro-cloud/proto/status_receiver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NewContainer struct {
	ContainerSize ContainerSize
}

type ContainerSize int

type MachinePoints struct {
	MachineId string
	Points int32
}

const (
	SMALL ContainerSize = iota
	MEDIUM
	LARGE
)

func CreateContainer (ctx context.Context) {

	minMem := int64(128 * 1024 * 1024) // 512MB livre mínimo
	minCpu := 100_000_000.0  
	
	machines := infra.Redis.Scan(ctx, 0, "machine-status:*", 0).Iterator()
	bestOne := FindBestMachines(ctx, machines, minCpu, minMem)
	if len(bestOne) != 0 {
		address := bestOne[0].MachineId
		fmt.Println("address", address)
		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

		c := cd.NewCommandDispatcherServiceClient(conn)
		if err != nil {
			fmt.Println("Erro com o client:", err)
		}
	
		resp, err := c.CreateCommand(ctx, &cd.CreateCommandRequest{
			ContainerSize: cd.ContainerSize_SMALL,
		})

		if err != nil {
			fmt.Println("Erro na req:", err)
		}
		fmt.Printf("Container criado em %v com id", resp.ContainerId)
	}
	
}

func FindBestMachines(ctx context.Context, iter interface{ Next(context.Context) bool; Val() string }, containerCpu float64, containerMem int64) []*MachinePoints {
	minMem := int64(128 * 1024 * 1024) // 512MB livre mínimo
	minCpu := 100_000_000.0            // 0.1 core livre mínimo (em ns/s)

	var valid []*MachinePoints

	for iter.Next(ctx) {
		key := iter.Val()

		value, err := infra.Redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var machineStatus sr.ImageStatus
		if err := json.Unmarshal([]byte(value), &machineStatus); err != nil {
			fmt.Println("Erro lendo dados da máquina:", err)
			continue
		}

		cpuUsage := parseFloatOrZero(machineStatus.CpuUsage)
		cpuTotal := parseFloatOrZero(machineStatus.CpuTotal)
		ramTotal := parseIntOrZero(machineStatus.RamTotal)
		ramUsage := parseIntOrZero(machineStatus.RamUsage)

		cpuLivre := cpuTotal - cpuUsage
		ramLivre := ramTotal - ramUsage

		if cpuLivre < minCpu || ramLivre < minMem {
			continue
		}

		if cpuLivre < containerCpu || ramLivre < containerMem {
			continue
		}

		cabemPorCpu := int32(cpuLivre / containerCpu)
		cabemPorMem := int32(float64(ramLivre) / float64(containerMem))
		capacidade := cabemPorCpu
		if cabemPorMem < capacidade {
			capacidade = cabemPorMem
		}

		pontos := int32(1000) - capacidade
		if pontos < 0 {
			pontos = 0
		}

		valid = append(valid, &MachinePoints{
			MachineId: machineStatus.MachineId,
			Points:    pontos,
		})
	}

	return valid
}

func parseIntOrZero(valor string) int64 {
	valor = strings.TrimSpace(valor) // remove \n, \r, espaços extras
	if valor == "" {
		return 0
	}
	numero, err := strconv.ParseInt(valor, 10, 64)
	if err != nil {
		fmt.Println("Erro parseando int:", valor, "-", err)
		return 0
	}
	return numero
}

func parseFloatOrZero(valor string) float64 {
	valor = strings.TrimSpace(valor) // mesma correção aqui
	if valor == "" {
		return 0
	}
	numero, err := strconv.ParseFloat(valor, 64)
	if err != nil {
		fmt.Println("Erro parseando float:", valor, "-", err)
		return 0
	}
	return numero
}

	// eu tenho que buscar fechar uma maquina por inteiro? na minha visao isso parece o correto
	// pesos 
	// pegar o menor peso?
	// desconsidera o que nao consegue rodar
	// se for um peso muito diferente pega esse, senao, pega o menor
