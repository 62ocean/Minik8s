package replicaset

import "time"

const (
	// 每隔5s检查一次pod的replica数量是否符合要求
	checkInterval = 5
)

func createReplicaset(replicasetConfig ReplicaSet) error {
	go func() {
		for range time.Tick(time.Second * checkInterval) {

		}
	}()
	return nil
}

func deleteReplicaset(replicasetConfig ReplicaSet) error {

}

func syncReplicaNum(replicasetConfig ReplicaSet) error {

}

func morePods() error {

}

func fewerPods() error {

}
