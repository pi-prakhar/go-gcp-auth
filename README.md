Desired structure

project-root/
├── api/
│   ├── proto/
│   │   ├── service1/
│   │   │   ├── v1/
│   │   │   │   └── service1.proto
│   │   │   └── service1.proto
│   │   ├── service2/
│   │   │   ├── v1/
│   │   │   │   └── service2.proto
│   │   │   └── service2.proto
│   │   └── common/
│   │       └── common.proto
│   ├── genproto/
│   │   ├── service1/
│   │   │   └── v1/
│   │   │       ├── service1.pb.go
│   │   │       └── service1_grpc.pb.go
│   │   └── service2/
│   │       └── v1/
│   │           ├── service2.pb.go
│   │           └── service2_grpc.pb.go
│   └── swagger/
├── cmd/
│   ├── service1/
│   │   └── main.go
│   ├── service2/
│   │   └── main.go
│   └── ...
├── internal/
│   ├── service1/
│   │   ├── handler/
│   │   │   └── grpc.go
│   │   ├── repository/
│   │   └── usecase/
│   ├── service2/
│   │   ├── handler/
│   │   │   └── grpc.go
│   │   ├── repository/
│   │   └── usecase/
│   └── ...
├── pkg/
│   ├── common/
│   ├── database/
│   ├── logger/
│   └── grpc/
│       ├── client/
│       │   └── client.go
│       └── server/
│           └── server.go
├── configs/
│   ├── service1.yaml
│   └── service2.yaml
├── deployments/
│   ├── docker/
│   │   ├── service1.Dockerfile
│   │   └── service2.Dockerfile
│   ├── kubernetes/
│   │   ├── service1-deployment.yaml
│   │   └── service2-deployment.yaml
│   ├── swarm/
│   │   └── docker-compose.yaml
│   └── docker-compose.yaml
├── scripts/
│   ├── ci/
│   │   └── build.sh
│   ├── deploy/
│   │   └── deploy.sh
│   └── protoc/
│       └── generate.sh
├── .gitignore
├── .gitlab-ci.yml
├── go.mod
├── buf.yaml
└── README.md