package util

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	client   *clientv3.Client
	Services []*ServiceInfo
}

type ServiceInfo struct {
	ServiceID      string
	ServiceName    string
	ServiceHost    string
	ServicePort    int
	ServiceAddress string
}

func NewClient() *Client {
	config := clientv3.Config{
		Endpoints:   []string{"101.132.107.3:2379"},
		DialTimeout: 10 * time.Second,
		Context:     context.Background(),
	}
	client, err := clientv3.New(config)
	if err != nil {
		log.Fatal(err)
	}
	return &Client{client: client}
}

func (this *Client) LoadService() {
	kv := clientv3.NewKV(this.client)
	res, err := kv.Get(context.TODO(), "/services", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range res.Kvs {
		this.parseService(v.Key, v.Value)
		fmt.Println(string(v.Key), string(v.Value))
	}
}

func (this *Client) parseService(key, address []byte) {

	reg := regexp.MustCompile(`/services/([\w-]+)/(\w+)`)
	reg1 := regexp.MustCompile(`([.0-9]+):(\d+)`)
	if reg.Match(key) && reg1.Match(address) {
		match := reg.FindSubmatch(key)
		sid := match[1]
		sName := match[2]

		match1 := reg1.FindSubmatch(address)
		host := match1[1]
		port := match1[2]
		address := string(host) + ":" + string(port)
		portInt, _ := strconv.Atoi(string(port))
		this.Services = append(this.Services, &ServiceInfo{
			ServiceID:      string(sid),
			ServiceName:    string(sName),
			ServiceHost:    string(host),
			ServicePort:    portInt,
			ServiceAddress: address,
		})
	}
}

func (this *Client) GetService(sname string, method string, encodeFunc EncodeRequestFunc) Endpoint {
	services := make([]*ServiceInfo, 0)
	for _, service := range this.Services {
		if service.ServiceName == sname {
			services = append(services, service)
		}
	}

	if len(services) != 0 {
		rand.Seed(time.Now().UnixNano())
		service := services[rand.Intn(len(services))]

		fmt.Println(service)
		return func(ctx context.Context, requestParam interface{}) (response interface{}, err error) {
			httpClient := http.DefaultClient
			req, _ := http.NewRequest(strings.ToUpper(method), "http://"+service.ServiceAddress, nil)
			err = encodeFunc(context.Background(), req, requestParam)
			if err != nil {
				return nil, err
			}

			res, err := httpClient.Do(req)
			defer res.Body.Close()
			if err != nil {
				return nil, err
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			return string(body), nil
		}

	}

	return nil
}
