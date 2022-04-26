module connection-microservice

go 1.18

replace github.com/XWS-BSEP-TIM1-2022/dislinkt/util => ./../../util

require (
	github.com/XWS-BSEP-TIM1-2022/dislinkt/util v0.0.0-20220419090605-7ed74d3dfc18
	github.com/neo4j/neo4j-go-driver/v4 v4.4.2
	github.com/opentracing/opentracing-go v1.2.0
	go.mongodb.org/mongo-driver v1.9.0
	google.golang.org/grpc v1.46.0
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/net v0.0.0-20220421235706-1d1ef9303861 // indirect
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220422154200-b37d22cd5731 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
