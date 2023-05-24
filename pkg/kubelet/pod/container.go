package pod

import "C"
import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	volume2 "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"k8s/object"
	"log"
	"strconv"
	"strings"
)

var Client = newClient()
var Ctx = context.Background()

func newClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	defer cli.Close()
	if err != nil {
		panic(err)
	}
	return cli
}

/*-----------------------Image------------------------*/

// PullImages 拉取本地没有的镜像们，以及pause要用到的镜像
func PullImages(images []string) error {
	images = append(images, "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.6")
	for _, image := range images {
		existed, err := isImageExist(image)
		if err != nil {
			return err
		}
		if !existed {
			err1 := pullSingleImage(image)
			log.Printf("Succeesfully pull image %s\n", image)
			if err1 != nil {
				return err1
			}
		} else {
			log.Printf("Image %s is already existed\n", image)
		}
	}
	return nil
}

// ListImage 列出镜像
func ListImage() error {
	images, err := Client.ImageList(Ctx, types.ImageListOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	for _, image := range images {
		fmt.Println(image.RepoTags)
	}
	return nil
}

// 查看本地是否有该镜像
func isImageExist(name string) (bool, error) {
	curList, err := Client.ImageList(Ctx, types.ImageListOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	for _, image := range curList {
		for _, curName := range image.RepoTags {
			if curName == name {
				return true, nil
			}
			tmp := name + ":latest"
			if tmp == curName {
				return true, nil
			}
		}
	}
	return false, nil
}

// 通过网络拉取单个镜像
func pullSingleImage(image string) error {
	log.Printf("Prepare to pull image:%s\n", image)
	out, err := Client.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		fmt.Printf("PullSingleImage: Fail to pull image, err:%v\n", err)
		return err
	}
	defer out.Close()
	io.Copy(log.Writer(), out)
	return nil
}

/*-----------------------Volume------------------------*/

// 创建数据卷们
func createVolumes(volumesConfig []object.VolumeConfig) ([]volume2.Volume, error) {
	var result []volume2.Volume
	for _, config := range volumesConfig {
		existed, err := isVolumeExisted(config.Name)
		if err != nil {
			fmt.Println(err)
			return result, err
		}

		if existed {
			log.Printf("Volume %s already existed, no need to create\n", config.Name)
			continue
		} else {
			newVolume, err1 := Client.VolumeCreate(Ctx, volume2.CreateOptions{
				Name: config.Name,
			})
			if err1 != nil {
				fmt.Println(err1)
				return result, err1
			}
			log.Printf("Successfully create Volume %s\n", config.Name)
			result = append(result, newVolume)
		}
	}

	return result, nil
}

// 获取已创建的数据卷
func lisVolumes() ([]*volume2.Volume, error) {
	var ret []*volume2.Volume
	volumes, err := Client.VolumeList(context.Background(), filters.NewArgs())
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	for _, vol := range volumes.Volumes {
		ret = append(ret, vol)
	}
	return ret, nil
}

// 查看是否已有该卷
func isVolumeExisted(name string) (bool, error) {
	volumes, err := Client.VolumeList(context.Background(), filters.NewArgs())
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	for _, vol := range volumes.Volumes {
		if vol.Name == name {
			return true, nil
		}
	}
	return false, nil
}

/*----------------------Container------------------------*/

// ListContainer 列出镜像
func ListContainer() ([]types.Container, error) {
	containers, err := Client.ContainerList(Ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	for _, c := range containers {
		str := "Name: " + c.Names[0] + " Status: " + c.Status
		fmt.Println(str)
	}
	return containers, nil
}

// SyncLocalContainer 查看本地是否正在运行该容器
func SyncLocalContainer(container object.ContainerMeta) bool {
	curList, err := Client.ContainerList(Ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	for _, c := range curList {
		for _, curName := range c.Names {
			// ps. docker给出的容器名称前面都会加个"/"（虽然不知道是为啥）
			if curName == "/"+container.Name {
				if strings.Contains(c.Status, "Exited") {
					fmt.Println("Container " + container.Name + " is exited, try to restart it now")
					StartContainer(container.ContainerID)
				}
				return true
			}
		}
	}
	fmt.Println("Container " + container.Name + " is not existed, try to create it now")
	return false
}

// CreateContainers 创建容器们
func CreateContainers(containerConfigs []object.Container, podName string) ([]object.ContainerMeta, error) {
	var result []object.ContainerMeta
	var totalPort []int
	dupMap := make(map[int32]bool)

	//port, 各个容器共享同一个network命名空间即可通过localhost进行相互访问
	//用container网络模式跟pause容器（pause自己用的是默认的bridge模式）共享
	for _, config := range containerConfigs {
		for _, port := range config.Ports {
			if !dupMap[port.Port] {
				dupMap[port.Port] = true
				totalPort = append(totalPort, int(port.Port))
			}
		}
	}

	//create pause container
	pauseID, err3 := createPause(&totalPort, podName)
	if err3 != nil {
		fmt.Println(err3.Error())
		return nil, err3
	}
	log.Println("OnCreate pause container")
	result = append(result, object.ContainerMeta{Name: "pause_" + podName, ContainerID: pauseID})

	for _, config := range containerConfigs {
		// volume mount
		var mounts []mount.Mount
		mountType := mount.TypeBind
		if config.VolumeMounts != nil {
			for _, it := range config.VolumeMounts {
				existed, err1 := isVolumeExisted(it.Name)
				if err1 != nil {
					fmt.Println(err1)
					return nil, err1
				}
				if existed {
					// 若使用的是定义在pod中的emptyDir，则是挂载volume，否则就是挂载至宿主机指定目录
					mountType = mount.TypeVolume
				}
				mounts = append(mounts, mount.Mount{
					Type:     mountType,
					Source:   it.Name,
					Target:   it.MountPath,
					ReadOnly: it.ReadOnly,
				})
			}
		}
		log.Println("Add volume mounting config")

		// resource
		resourceConfig := container.Resources{}
		if config.Resources.Limits.Cpu != "" {
			resourceConfig.NanoCPUs = parseCPU(config.Resources.Limits.Cpu)
		}
		if config.Resources.Limits.Memory != "" {
			resourceConfig.Memory = parseMemory(config.Resources.Limits.Memory)
		}
		log.Println("Add resource config")

		// create container (可使用localhost通信)
		// k8s中pod内容器共享了net、ipc、uts namespace
		resp, err := Client.ContainerCreate(context.Background(), &container.Config{
			Image:      config.Image,
			Entrypoint: config.Command,
			Cmd:        config.Args,
		}, &container.HostConfig{
			NetworkMode: container.NetworkMode("container:" + pauseID),
			Mounts:      mounts,
			IpcMode:     container.IpcMode("container:" + pauseID),
			//UTSMode:     container.UTSMode("container:" + pauseID),
			Resources: resourceConfig,
		}, nil, nil, config.Name+"_"+podName)
		if err != nil {
			return nil, err
		}
		log.Printf("OnCreate container %s\n", resp.ID)

		// record container ID
		result = append(result, object.ContainerMeta{
			Name:        config.Name + "_" + podName,
			ContainerID: resp.ID,
		})
	}
	return result, nil
}

// StartContainer 启动容器
func StartContainer(containerID string) {
	err := Client.ContainerStart(Ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		log.Println("Container ", containerID, "starts successfully")
	}
}

// StopContainer 停止容器
func StopContainer(containerID string) {
	timeout := 10
	err := Client.ContainerStop(Ctx, containerID, container.StopOptions{
		Timeout: &timeout,
	})
	if err != nil {
		if client.IsErrNotFound(err) {
			log.Println("container " + containerID + " is not found, no need to stop it")
		} else {
			fmt.Println(err.Error())
		}
	} else {
		log.Printf("Container %s is stopped\n", containerID)
	}
}

// RemoveContainer 删除容器
func RemoveContainer(containerID string) {
	err := Client.ContainerRemove(Ctx, containerID, types.ContainerRemoveOptions{})
	if err != nil {
		if client.IsErrNotFound(err) {
			log.Println("container " + containerID + " is not found, no need to remove it")
		} else {
			fmt.Println(err.Error())
		}
	} else {
		log.Printf("Container %s is removed\n", containerID)
	}
}

// 创建pause容器用于管理网络
func createPause(ports *[]int, podName string) (string, error) {
	var exports nat.PortSet
	exports = make(nat.PortSet, len(*ports))
	for _, port := range *ports {
		// 默认使用TCP协议
		p, err := nat.NewPort("tcp", strconv.Itoa(port))
		if err != nil {
			return "", err
		}
		exports[p] = struct{}{}
	}

	resp, err := Client.ContainerCreate(context.Background(), &container.Config{
		Image:        "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.6",
		ExposedPorts: exports,
	}, &container.HostConfig{
		IpcMode: container.IpcMode("shareable"),
		//UTSMode: container.UTSMode("shareable"),
	}, nil, nil, "pause_"+podName)
	return resp.ID, err
}

/*-----------------------Tools------------------------*/
func parseCPU(cpu string) int64 {
	length := len(cpu)
	result := 0.0
	if cpu[length-1] == 'm' {
		result, _ = strconv.ParseFloat(cpu[:length-1], 32)
		result *= 1e3
	} else {
		result, _ = strconv.ParseFloat(cpu[:length], 32)
	}
	return int64(result)
}

func parseMemory(mem string) int64 {
	length := len(mem)
	result, _ := strconv.Atoi(mem[:length-1])
	mark := mem[length-1]
	if mark == 'K' || mark == 'k' {
		result *= 1024
	} else if mark == 'M' || mark == 'm' {
		result *= 1024 * 1024
	} else if mark == 'G' || mark == 'g' {
		result *= 1024 * 1024 * 1024
	}
	return int64(result)
}
