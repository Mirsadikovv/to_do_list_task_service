package client

import (
	"go_schedule_service/config"

	"log"

	gr "go_schedule_service/genproto/group_service"
	st "go_schedule_service/genproto/student_service"

	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientI interface {
	Group() gr.GroupServiceClient
	Student() st.StudentServiceClient
}

type GrpcClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

func New(cfg config.Config) (*GrpcClient, error) {

	connUser, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("user service dial host: %v port :%v err:%v",
			cfg.UserServiceHost, cfg.UserServicePort, err)
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"group_service":   gr.NewGroupServiceClient(connUser),
			"student_service": st.NewStudentServiceClient(connUser),
		},
	}, nil
}

func (g *GrpcClient) Group() gr.GroupServiceClient {
	client, ok := g.connections["group_service"].(gr.GroupServiceClient)
	if !ok {
		log.Println("failed to assert type for group_service")
		return nil
	}
	return client
}

func (g *GrpcClient) Student() st.StudentServiceClient {
	client, ok := g.connections["student_service"].(st.StudentServiceClient)
	if !ok {
		log.Println("failed to assert type for student_service")
		return nil
	}
	return client
}
