package docker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type ContainerResources struct {
	ID        string
	Name      string
	MemoryMax int64   // em bytes, 0 = sem limite
	NanoCPUs  int64   // CPUs * 1e9 (ex: 1.5 CPU = 1500000000)
	CPUQuota  int64   // microssegundos por período
	CPUPeriod int64   // duração do período em microssegundos
	CPUShares int64   // peso relativo (legado, cgroup v1)
}

var Client *client.Client

func NewDockerClient() {
	var err error
	Client, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Errorf("Erro instanciando docker: %v", err)
	}

}

func CloseDockerClient () {
	Client.Close()
}

func CreateContainer (ctx context.Context) (*string, error) {

	reader, err := Client.ImagePull(ctx, "nginx:latest", image.PullOptions{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	io.Copy(io.Discard, reader)

	containerConfig := &container.Config{
		Image: "nginx:latest",
		
	}

	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			Memory: 512 * 1024 * 1024,
			NanoCPUs: 100_000_000,
		},
	}

	containerUid := generateContainerUid()
	
	resp, err := Client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerUid)
	if err != nil {
		return &resp.ID, err
	}

	err = Client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return &resp.ID, err
	}
	fmt.Printf("Conteiner criado e rodando: %v\n", resp.ID)
	return &resp.ID, nil
}

func ListContainers (ctx context.Context) ([]*ContainerResources, error) {
	 	summaries, err := Client.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []*ContainerResources
	for _, s := range summaries {
		res, err := GetContainerResources(ctx, s.ID)
		if err != nil {
			continue
		}
		result = append(result, res)
	}

	return result, nil
}

func GetContainerResources(ctx context.Context, containerID string) (*ContainerResources, error) {
	inspect, err := Client.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}

	hc := inspect.HostConfig

	return &ContainerResources{
		ID:        inspect.ID,
		Name:      inspect.Name,
		MemoryMax: hc.Memory,
		NanoCPUs:  hc.NanoCPUs,
		CPUQuota:  hc.CPUQuota,
		CPUPeriod: hc.CPUPeriod,
		CPUShares: hc.CPUShares,
	}, nil
}

func DownContainer (ctx context.Context, containerId string) (error) {
	return Client.ContainerKill(ctx, containerId, "SIGKILL")
}

func generateContainerUid () string {
	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	soma := sha256.Sum256([]byte(ts))
	return hex.EncodeToString(soma[:])[:12]
	
}