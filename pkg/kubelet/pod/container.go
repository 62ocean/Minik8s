package pod

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	volume2 "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/go-connections/nat"
	"io"
	"k8s/object"
	"k8s/pkg/kubelet/cache"
	"log"
	"os"
	"path/filepath"
	"runtime"
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
func SyncLocalContainer(containers map[string]*cache.ContainerMeta) bool {
	curList, err := Client.ContainerList(Ctx, types.ContainerListOptions{All: true})
	log.Print("SYNC container: ")
	log.Println(curList)
	totalNum := 0
	if err != nil {
		log.Println(err.Error())
		return false
	}
	for _, c := range curList {
		meta := containers[c.Names[0][1:]]
		if meta != nil {
			totalNum++
			if strings.Contains(c.Status, "Exited") {
				log.Println("Container " + c.Names[0] + " is exited, try to restart it now")
				StartContainer(meta.ContainerID)
			}
		}
	}
	if totalNum < len(containers) {
		log.Println("totalNum: " + strconv.Itoa(totalNum))
		log.Println("len of containers in pod: " + strconv.Itoa(len(containers)))
		return false
	}
	return true
}

// CreateContainers 创建容器们
func CreateContainers(containerConfigs []object.Container, podName string) (map[string]*cache.ContainerMeta, error) {
	result := make(map[string]*cache.ContainerMeta)
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
	result["pause_"+podName] = &cache.ContainerMeta{Name: "pause_" + podName, ContainerID: pauseID, InitialName: "pause"}

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
			resourceConfig.CPUQuota = parseCPU(config.Resources.Limits.Cpu)
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
			//ExposedPorts: nat.PortSet{},
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
		result[config.Name+"_"+podName] = &cache.ContainerMeta{
			Name:        config.Name + "_" + podName,
			ContainerID: resp.ID,
			InitialName: config.Name,
			Limit:       config.Resources.Limits,
		}

		// copy to container if needed
		if config.CopyFile != "" {
			err := copy(Ctx, config.CopyFile, config.CopyDst, resp.ID)
			if err != nil {
				log.Println(err.Error())
				return result, err
			}
			log.Println("Copy file successfully")
		}
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

func GetContainerStatus(containerID string, resources object.ContainerResources) (uint64, uint64, uint64, uint64, error) {
	containerStats, err := Client.ContainerStats(Ctx, containerID, false)
	// 这个container被意外删除了会抛出错误
	if err != nil {
		log.Println("get status of container: " + err.Error())
		return 0, 0, 0, 0, err
	}
	// 解析容器的状态数据
	var stats types.StatsJSON
	if err := json.NewDecoder(containerStats.Body).Decode(&stats); err != nil {
		fmt.Println(err.Error())
		return 0, 0, 0, 0, err
	}

	// 获取资源利用率
	// cpu
	cpuUsage := stats.CPUStats.CPUUsage.TotalUsage
	cpuSys := stats.CPUStats.SystemUsage
	if cpuSys == 0 {
		return 0, 0, 0, 0, nil
	}
	cpuUsage = cpuUsage * 1e5 / cpuSys
	var cpuLimit uint64
	if resources.Cpu != "" {
		cpuLimit = uint64(parseCPU(resources.Cpu))
	} else {
		cpuLimit = 1e5
	}

	// memory
	memoryUsage := stats.MemoryStats.Usage
	memoryLimit := stats.MemoryStats.Limit
	if resources.Memory != "" {
		memoryLimit = uint64(parseMemory(resources.Memory))
	}
	return cpuUsage, cpuLimit, memoryUsage, memoryLimit, nil
}

func getContainerIP(id string) string {
	containerInfo, err := Client.ContainerInspect(Ctx, id)
	if err != nil {
		fmt.Println(err.Error())
	}
	ipAddress := containerInfo.NetworkSettings.IPAddress
	log.Printf("IP ADDRESS: %s\n", ipAddress)
	return ipAddress
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

// 将文件复制到docker容器
func copy(ctx context.Context, file string, dest string, container string) error {
	srcPath := file
	dstPath := dest
	// Prepare destination copy info by stat-ing the container path.
	dstInfo := archive.CopyInfo{Path: dstPath}
	dstStat, err := Client.ContainerStatPath(ctx, container, dstPath)

	// If the destination is a symbolic link, we should evaluate it.
	if err == nil && dstStat.Mode&os.ModeSymlink != 0 {
		linkTarget := dstStat.LinkTarget
		if !system.IsAbs(linkTarget) {
			// Join with the parent directory.
			dstParent, _ := archive.SplitPathDirEntry(dstPath)
			linkTarget = filepath.Join(dstParent, linkTarget)
		}

		dstInfo.Path = linkTarget
		dstStat, err = Client.ContainerStatPath(ctx, container, linkTarget)
	}

	if err == nil {
		dstInfo.Exists, dstInfo.IsDir = true, dstStat.Mode.IsDir()
	}

	var (
		content         io.Reader
		resolvedDstPath string
	)

	// Prepare source copy info.
	srcInfo, err := archive.CopyInfoSourcePath(srcPath, true)
	if err != nil {
		return err
	}

	srcArchive, err := archive.TarResource(srcInfo)
	if err != nil {
		return err
	}
	defer srcArchive.Close()

	dstDir, preparedArchive, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
	if err != nil {
		return err
	}
	defer preparedArchive.Close()

	resolvedDstPath = dstDir
	content = preparedArchive

	options := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}
	return Client.CopyToContainer(ctx, container, resolvedDstPath, content, options)
}

/*-----------------------Tools------------------------*/
func parseCPU(cpu string) int64 {
	length := len(cpu)
	result := 0.0
	if cpu[length-1] == 'm' {
		// 这里的m指的k8s中的微核，转换为docker要求的纳秒返回
		result, _ = strconv.ParseFloat(cpu[:length-1], 32)
		result *= 1e2
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

// 计算CPU利用率
func calculateCPUPercentage(cpuUsage uint64, sysCPUUsage uint64, limit int64) float64 {
	cpuUsageDelta := float64(cpuUsage)
	sysCPUUsageDelta := float64(sysCPUUsage)
	//cpuUsagePercentage := cpuUsageDelta / sysCPUUsageDelta
	cpuNum := runtime.NumCPU()
	if limit >= 1 {
		return cpuUsageDelta / (sysCPUUsageDelta / float64(cpuNum) * float64(limit))
	} else {
		exactLimit := float64(limit) / 1e5
		return cpuUsageDelta / sysCPUUsageDelta / float64(cpuNum) * exactLimit
	}
}

// 计算内存利用率
func calculateMemoryPercentage(memoryUsage uint64, memoryLimit uint64) float64 {
	memoryUsageDelta := float64(memoryUsage)
	memoryLimitDelta := float64(memoryLimit)
	memoryUsagePercentage := (memoryUsageDelta / memoryLimitDelta) * 100

	return memoryUsagePercentage
}
