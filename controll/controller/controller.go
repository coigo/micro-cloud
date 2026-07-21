package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/coigo/micro-cloud/infra"
	proto "github.com/coigo/micro-cloud/proto/status_receiver"
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

func OnUpdate() {

}

func FindBestMachines(ctx context.Context, iter interface{ Next(context.Context) bool; Val() string }, containerCpu float64, containerMem int64) []MachinePoints {
	minMem := int64(512 * 1024 * 1024) // 512MB livre mínimo
	minCpu := 100_000_000.0            // 0.1 core livre mínimo (em ns/s)

	var valid []MachinePoints

	for iter.Next(ctx) {
		key := iter.Val()

		value, err := infra.Redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var machineStatus proto.ImageStatus
		if err := json.Unmarshal([]byte(value), &machineStatus); err != nil {
			fmt.Println("Erro lendo dados da máquina:", err)
			continue
		}

		cpuTotal, err := strconv.ParseFloat(machineStatus.CpuTotal, 64)
		if err != nil {
			fmt.Println("Erro parseando CpuTotal:", err)
			continue
		}
		cpuUsage, err := strconv.ParseFloat(machineStatus.CpuUsage, 64)
		if err != nil {
			fmt.Println("Erro parseando CpuUsage:", err)
			continue
		}
		ramTotal, err := strconv.ParseInt(machineStatus.RamTotal, 10, 64)
		if err != nil {
			fmt.Println("Erro parseando RamTotal:", err)
			continue
		}
		
		ramUsage, err := strconv.ParseInt(machineStatus.RamUsage, 10, 64)
		if err != nil {
			fmt.Println("Erro parseando RamUsage:", err)
			continue
		}

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

		valid = append(valid, MachinePoints{
			MachineId: machineStatus.MachineId,
			Points:    pontos,
		})
	}

	return valid
}

	// eu tenho que buscar fechar uma maquina por inteiro? na minha visao isso parece o correto
	// pesos 
	// pegar o menor peso?
	// desconsidera o que nao consegue rodar
	// se for um peso muito diferente pega esse, senao, pega o menor
