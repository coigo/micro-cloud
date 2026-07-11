package docker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type CPUMonitor struct {
	lastUsage int64
	lastTime  time.Time
	numCPUs   float64
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
			NanoCPUs: 500_000_000,
		},
	}
	
	resp, err := Client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "teste")
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

func ListContainers (ctx context.Context) ([]container.Summary, error) {
	 return Client.ContainerList(ctx, container.ListOptions{})
}

func DownContainer (ctx context.Context, containerId string) (error) {
	return Client.ContainerKill(ctx, containerId, "SIGKILL")
}