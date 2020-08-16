package client

import (
	"context"
	"fmt"
	ds "github.com/TheComputerDan/heimdall/proto/dockerService"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	// "net/http"
	"io"
)

func Start() {
	conn, err := grpc.Dial("localhost:8069", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := ds.NewDockerAgentClient(conn)

	g := gin.Default()

	g.GET("/containers", func(ctx *gin.Context) {

		req := &ds.ContainersRequest{}
		stream, err := client.Containers(context.Background(), req)
		if err != nil {
			panic(err)
		}

		for {
			containers, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			fmt.Println(containers)
		}
	})
}
