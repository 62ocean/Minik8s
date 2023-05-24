package pod

import (
	"fmt"
	object2 "k8s/object"
	"log"
)

var ipCnt = 0

func CreatePod(podConfig *object2.Pod) error {
	// 分配podip
	//localNodeNetWork := flannel.GetLocalNodeNetwork()
	////subnetPrefix: x.x.x
	//subnet := fmt.Sprintf("%s.%d", localNodeNetWork.SubnetPrefix, ipCnt)
	//ipCnt++
	//podConfig.IP = subnet

	// 拉取镜像
	var images []string
	for _, configItem := range podConfig.Spec.Containers {
		images = append(images, configItem.Image)
	}

	err := PullImages(images)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 创建emptyDir数据卷（pod中的各个容器共享）
	_, err = createVolumes(podConfig.Spec.Volumes)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 创建pod中的容器并运行
	var containerMeta []object2.ContainerMeta
	containerMeta, err = CreateContainers(podConfig.Spec.Containers, podConfig.Metadata.Name)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	podConfig.Spec.ContainerMeta = containerMeta

	// 打印容器信息
	for _, it := range containerMeta {
		log.Println(it.Name, " id:", it.ContainerID)
	}

	return nil
}

func StartPod(podConfig *object2.Pod) {
	containerMeta := podConfig.Spec.ContainerMeta
	// 开启容器
	for _, it := range containerMeta {
		StartContainer(it.ContainerID)
	}
}

func ClosePod(podConfig *object2.Pod) {
	containerMeta := podConfig.Spec.ContainerMeta
	// 关闭容器
	for _, it := range containerMeta {
		StopContainer(it.ContainerID)
	}
}

func RemovePod(podConfig *object2.Pod) {
	log.Println("remove pod now")
	containerMeta := podConfig.Spec.ContainerMeta
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

// SyncPod 返回的bool值若为true表示pod状态需要更新了
func SyncPod(podConfig *object2.Pod) (update bool, err error) {
	for _, container := range podConfig.Spec.ContainerMeta {
		if SyncLocalContainer(container) == false {
			// container目前不存在了，我们选择把容器都关了重新起个pod
			RemovePod(podConfig)
			err = CreatePod(podConfig)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}
