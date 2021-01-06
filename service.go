package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"myetcd/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	name := flag.String("name", "", "输入服务名字")
	port := flag.Int("p", 0, "输入端口号")
	flag.Parse()
	fmt.Println(*name, *port)
	if *name == "" || *port == 0 {
		log.Fatal("服务信息输入错误")
	}

	r := mux.NewRouter()

	r.HandleFunc(`/prod/{id:\d+}`, func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		res := "get prod by id:" + vars["id"]
		_, _ = writer.Write([]byte(res))
	})


	s := util.NewService()
	service := util.ServiceInfo{
		ServiceID:      uuid.New().String(),
		ServiceName:    *name,
		ServiceAddress: "127.0.0.1",
		ServicePort:    *port,
	}


	errChan := make(chan error)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", service.ServicePort),
		Handler: r}
	go func() {
		err := s.RegService(service)
		if err != nil {
			errChan <- err
			return
		}
		err = httpServer.ListenAndServe()
		if err != nil {
			errChan <- err
			return
		}
	}()

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sig)
	}()
	getErr := <-errChan
	log.Println("发生异常,服务正在停止")
	err := s.UnRegService(service.ServiceID)
	if err != nil {
		log.Fatal(err)
	}
	err = httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(getErr)
}
