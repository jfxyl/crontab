package job

import (
	"context"
	"crontab/internal/common"
	"crontab/internal/global"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net"
	"time"
)

func InitRegister() (err error) {
	var (
		regKey  string
		localIP string
	)
	if localIP, err = getLocalIP(); err != nil {
		return
	}
	regKey = common.JOB_WORKER_DIR + localIP
	go keepOnline(regKey)
	return
}

func getLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet
		isIpNet bool
	)
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	for _, addr = range addrs {
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}
	err = common.ERR_NO_LOCAL_IP_FOUND
	return
}

func keepOnline(regKey string) {
	var (
		err                    error
		ctx                    context.Context
		cancelFunc             context.CancelFunc
		leaseGrantResp         *clientv3.LeaseGrantResponse
		leaseKeepAliveResp     *clientv3.LeaseKeepAliveResponse
		leaseKeepAliveRespChan <-chan *clientv3.LeaseKeepAliveResponse
	)
	for {
		ctx, cancelFunc = context.WithCancel(context.Background())
		if leaseGrantResp, err = global.EtcdClient.Grant(context.Background(), 10); err != nil {
			goto RETRY
		}
		if leaseKeepAliveRespChan, err = global.EtcdClient.KeepAlive(ctx, leaseGrantResp.ID); err != nil {
			goto RETRY
		}
		if _, err = global.EtcdClient.Put(ctx, regKey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
			goto RETRY
		}
		for {
			select {
			case leaseKeepAliveResp = <-leaseKeepAliveRespChan:
				if leaseKeepAliveResp == nil {
					goto RETRY
				}
			}
		}

	RETRY:
		cancelFunc()
		time.Sleep(1 * time.Second)
	}
	return
}
