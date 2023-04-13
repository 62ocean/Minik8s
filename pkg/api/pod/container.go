package pod

import "C"
import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"time"
)

var Client = newClient()
var Ctx = context.Background()

func newClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	defer cli.Close()
	if err != nil {
		panic(err)
	}
	return cli
}

// 列出镜像
func listImage() {
	images, err := Client.ImageList(Ctx, types.ImageListOptions{})
	log(err)

	for _, image := range images {
		fmt.Println(image)
	}
}

// 创建容器
func createContainer(*container.Config) string {
	exports := make(nat.PortSet, 10)
	port, err := nat.NewPort("tcp", "80")
	log(err)
	exports[port] = struct{}{}
	config := &container.Config{Image: "nginx", ExposedPorts: exports}

	portBind := nat.PortBinding{HostPort: "80"}
	portMap := make(nat.PortMap, 0)
	tmp := make([]nat.PortBinding, 0, 1)
	tmp = append(tmp, portBind)
	portMap[port] = tmp
	hostConfig := &container.HostConfig{PortBindings: portMap}
	// networkingConfig := &network.NetworkingConfig{}
	containerName := "hel"
	body, err := Client.ContainerCreate(Ctx, config, hostConfig, nil, containerName)
	log(err)
	fmt.Printf("ID: %s\n", body.ID)
	return body.ID
}

// 启动容器
func startContainer(containerID string) {
	err := Client.ContainerStart(Ctx, containerID, types.ContainerStartOptions{})
	log(err)
	if err == nil {
		fmt.Println("容器", containerID, "启动成功")
	}
}

// 停止容器
func stopContainer(containerID string) {
	timeout := time.Second * 10
	err := Client.ContainerStop(Ctx, containerID, &timeout)
	if err != nil {
		log(err)
	} else {
		fmt.Printf("容器%s已经被停止\n", containerID)
	}
}

// 删除容器
func removeContainer(containerID string) (string, error) {
	err := Client.ContainerRemove(Ctx, containerID, types.ContainerRemoveOptions{})
	log(err)
	return containerID, err
}

func log(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
}
