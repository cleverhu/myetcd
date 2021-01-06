package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"myetcd/services"
	"myetcd/util"
)

func main() {

	c := util.NewClient()
	c.LoadService()
	for i := 0; i < 10; i++ {
		endpoint := c.GetService("prodService", "get", services.ProdEncodeFunc)
		res, err := endpoint(context.Background(), services.ProdRequest{Id: 111})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res)
		<-time.After(100 * time.Millisecond)
	}

}
