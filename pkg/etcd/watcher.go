package etcd

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

// 此文件包含了对etcd watch机制的封装函数
// 在handler中实现listener接口并添加监听回调函数即可实时对etcd状态变化做出响应

var (
	timeOut = time.Duration(3) * time.Second // 超时
)

// Listener 对外通知
type Listener interface {
	OnSet(kv mvccpb.KeyValue)
	OnCreate(kv mvccpb.KeyValue)
	OnModify(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue)
	OnDelete(kv mvccpb.KeyValue)
}

// EtcdWatcher ETCD key监视器
type EtcdWatcher struct {
	client       *clientv3.Client // etcd client
	waitGroup    sync.WaitGroup
	listener     Listener
	mutex        sync.Mutex
	closeHandler map[string]func()
}

// NewEtcdWatcher 构造
func NewEtcdWatcher(servers []string) (*EtcdWatcher, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   servers,
		DialTimeout: timeOut,
	})
	if err != nil {
		return nil, err
	}

	ew := &EtcdWatcher{
		client:       cli,
		closeHandler: make(map[string]func()),
	}

	return ew, nil
}

// AddWatch 添加监视
func (watcher *EtcdWatcher) AddWatch(key string, prefix bool, listener Listener) bool {
	watcher.mutex.Lock()
	defer watcher.mutex.Unlock()
	if _, ok := watcher.closeHandler[key]; ok {
		// 已有对该资源的监视
		return false
	}
	ctx, cancel := context.WithCancel(context.Background())
	watcher.closeHandler[key] = cancel

	watcher.waitGroup.Add(1)
	go watcher.watch(ctx, key, prefix, listener)

	return true
}

// RemoveWatch 删除监视
func (watcher *EtcdWatcher) RemoveWatch(key string) bool {
	watcher.mutex.Lock()
	defer watcher.mutex.Unlock()
	cancel, ok := watcher.closeHandler[key]
	if !ok {
		// 不存在对该资源的监视
		return false
	}
	cancel()
	delete(watcher.closeHandler, key)

	return true
}

// ClearWatch 清除所有监视
func (watcher *EtcdWatcher) ClearWatch() {
	watcher.mutex.Lock()
	defer watcher.mutex.Unlock()
	for k := range watcher.closeHandler {
		watcher.closeHandler[k]()
	}
	watcher.closeHandler = make(map[string]func())
}

// Close 关闭
func (watcher *EtcdWatcher) Close(wait bool) {
	watcher.ClearWatch()

	if wait {
		watcher.waitGroup.Wait()
	}

	watcher.client.Close()
	watcher.client = nil
}

// watch 监听的主要逻辑（协程）
func (watcher *EtcdWatcher) watch(ctx context.Context, key string, prefix bool, listener Listener) error {
	defer watcher.waitGroup.Done()

	ctx1, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	var getResp *clientv3.GetResponse
	var err error
	if prefix {
		getResp, err = watcher.client.Get(ctx1, key, clientv3.WithPrefix())
	} else {
		getResp, err = watcher.client.Get(ctx1, key)
	}
	if err != nil {
		return err
	}

	for _, ev := range getResp.Kvs {
		listener.OnSet(*ev)
	}

	var watchChan clientv3.WatchChan
	if prefix {
		watchChan = watcher.client.Watch(context.Background(), key, clientv3.WithPrefix(), clientv3.WithRev(getResp.Header.Revision+1), clientv3.WithPrevKV())
	} else {
		watchChan = watcher.client.Watch(context.Background(), key, clientv3.WithRev(getResp.Header.Revision+1), clientv3.WithPrevKV())
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case resp := <-watchChan:
			err := resp.Err()
			if err != nil {
				return err
			}
			for _, ev := range resp.Events {
				switch ev.Type {
				case mvccpb.PUT:
					if ev.Kv.Version == 1 {
						listener.OnCreate(*ev.Kv)
					} else {
						listener.OnModify(*ev.Kv, *ev.PrevKv)
					}
				case mvccpb.DELETE:
					listener.OnDelete(*ev.Kv)
				}
			}
		}
	}
}
