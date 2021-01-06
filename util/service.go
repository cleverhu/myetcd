package util

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

type Service struct {
	client *clientv3.Client
}

func NewService() *Service {
	config := clientv3.Config{
		Endpoints:   []string{"101.132.107.3:2379"},
		DialTimeout: 10 * time.Second,
		Context:     context.Background(),
	}
	client, err := clientv3.New(config)
	if err != nil {
		log.Fatal(err)
	}
	return &Service{client: client}
}

func (this *Service) RegService(s ServiceInfo) error {
	ctx := context.Background()
	lease := clientv3.NewLease(this.client)
	leaseResponse, err := lease.Grant(ctx, 30)
	if err != nil {
		return err
	}

	kv := clientv3.NewKV(this.client)
	pref := "/services/"
	_, err = kv.Put(ctx, pref+s.ServiceID+"/"+s.ServiceName,
		s.ServiceAddress+fmt.Sprintf(":%d", s.ServicePort),
		clientv3.WithLease(leaseResponse.ID))
	if err != nil {
		return err
	}

	leaseChan, err := lease.KeepAlive(context.TODO(), leaseResponse.ID)
	if err != nil {
		return err
	}

	go leaseKeepLive(leaseChan)

	return err
}

func leaseKeepLive(l <-chan *clientv3.LeaseKeepAliveResponse) {
	for true {
		select {
		case ret := <-l:
			if ret != nil {
				fmt.Println(time.Now(), ":续租成功")
			}
		default:
			<-time.After(1 * time.Second)
		}
	}
}

func (this *Service) UnRegService(id string) error {
	ctx := context.Background()

	kv := clientv3.NewKV(this.client)
	pref := "/services/"
	_, err := kv.Delete(ctx, pref+id, clientv3.WithPrefix())

	return err
}
