package pod

import (
	"fmt"
	object2 "k8s/object"
	"k8s/pkg/apiserver/flannel"
	"k8s/pkg/kubelet/cache"
	"log"
)

var ipCnt = 2

func CreatePod(podConfig *object2.Pod) ([]cache.ContainerMeta, error) {
	// 分配podip
	localNodeNetWork := flannel.GetLocalNodeNetwork()
	fmt.Println(localNodeNetWork.SubnetPrefix)
	//subnetPrefix: x.x.x
	subnet := fmt.Sprintf("%s.%d", localNodeNetWork.SubnetPrefix, ipCnt)
	ipCnt++
	podConfig.IP = subnet

	// 拉取镜像
	var images []string
	for _, configItem := range podConfig.Spec.Containers {
		images = append(images, configItem.Image)
	}

	err := PullImages(images)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// 创建emptyDir数据卷（pod中的各个容器共享）
	_, err = createVolumes(podConfig.Spec.Volumes)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// 创建pod中的容器并运行
	var containerMeta []cache.ContainerMeta
	containerMeta, err = CreateContainers(podConfig.Spec.Containers, podConfig.Metadata.Name)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// 打印容器信息
	for _, it := range containerMeta {
		log.Println(it.Name, " id:", it.ContainerID)
	}

	return containerMeta, nil
}

func StartPod(containers []cache.ContainerMeta) {
	// 开启容器
	for _, it := range containers {
		StartContainer(it.ContainerID)
	}
}

func ClosePod(containers []cache.ContainerMeta) {
	// 关闭容器
	for _, it := range containers {
		StopContainer(it.ContainerID)
	}
}

func RemovePod(podConfig *cache.PodCache) {
	log.Println("remove pod now")
	containerMeta := podConfig.ContainerMeta
	// 关闭容器
	for _, it := range containerMeta {
		log.Println("stop container " + it.Name)
		StopContainer(it.ContainerID)
	}
	// 删除容器
	for _, it := range containerMeta {
		log.Println("remove container " + it.Name)
		RemoveContainer(it.ContainerID)
	}
}

// SyncPod 返回的bool值若为true表示pod需要更新重启了
func SyncPod(podConfig *cache.PodCache) (update bool, err error) {
	for _, container := range podConfig.ContainerMeta {
		if SyncLocalContainer(container) == false {
			// container目前不存在了，我们选择把pod都关了重新起个pod
			return true, nil
		}
	}
	podConfig.PodStorage.RunningMetrics = GetStatusOfPod(podConfig)
	return false, nil
}

// GetStatusOfPod 获取pod状态
func GetStatusOfPod(cache *cache.PodCache) object2.RunningMetrics {
	var totalCpuUse, totalCpuLimit, totalMemUse, totalMemLimit uint64
	for _, container := range cache.ContainerMeta {
		a, b, c, d, _ := GetContainerStatus(container.ContainerID, container.Limit)
		//fmt.Printf("cpu use: %d\n", a)
		//fmt.Printf("cpu limit: %d\n", b)
		//fmt.Printf("mem use: %d\n", c)
		//fmt.Printf("mem limit: %d\n", d)

		totalCpuUse += a
		totalCpuLimit += b
		totalMemUse += c
		totalMemLimit += d
	}
	return object2.RunningMetrics{
		CPUUtil: float64(totalCpuUse) / float64(totalCpuLimit),
		MemUtil: float64(totalMemUse) / float64(totalMemLimit),
	}
}
