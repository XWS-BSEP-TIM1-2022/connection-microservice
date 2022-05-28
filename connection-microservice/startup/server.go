package startup

import (
	"connection-microservice/application"
	"connection-microservice/infrastructure/api"
	"connection-microservice/infrastructure/persistance"
	"connection-microservice/model"
	"connection-microservice/startup/config"
	"fmt"
	connectionService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/connection"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/token"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	otgo "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type Server struct {
	config      *config.Config
	tracer      otgo.Tracer
	closer      io.Closer
	jwtManager  *token.JwtManager
	neo4jDriver neo4j.Driver
}

func NewServer(config *config.Config) *Server {
	tracer, closer := tracer.Init(config.ConnectionServiceName)
	otgo.SetGlobalTracer(tracer)
	jwtManager := token.NewJwtManagerDislinkt(config.ExpiresIn)
	return &Server{
		config:     config,
		tracer:     tracer,
		closer:     closer,
		jwtManager: jwtManager,
	}
}

func (server *Server) GetTracer() otgo.Tracer {
	return server.tracer
}

func (server *Server) GetCloser() io.Closer {
	return server.closer
}

func (server *Server) Start() {
	server.neo4jDriver = server.initNeo4jClient()
	connectionStore := server.initConnectionStore(server.neo4jDriver)
	blockStore := server.initBlockStore(server.neo4jDriver)
	initConnectionService := server.initConnectionService(connectionStore)
	blockService := server.initBlockService(blockStore, connectionStore)
	connectionHandler := server.initConnectionHandler(initConnectionService, blockService)

	server.startGrpcServer(connectionHandler)
}

func (server *Server) Stop() {
	log.Println("stopping server")

}

func (server *Server) initNeo4jClient() neo4j.Driver {
	driver, err := persistance.GetDriver(server.config.ConnectionDBURI, server.config.ConnectionDBUsername, server.config.ConnectionDBPassword)
	if err != nil {
		log.Fatal(err)
	}
	return driver
}

func (server *Server) startGrpcServer(connectionHandler *api.ConnectionHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	log.Println(fmt.Sprintf("started grpc server on localhost:%s", server.config.Port))
	connectionService.RegisterConnectionServiceServer(grpcServer, connectionHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initConnectionStore(driver neo4j.Driver) model.ConnectionStore {
	store := persistance.NewConnectionNeo4jStore(driver)
	return store
}

func (server *Server) initConnectionService(store model.ConnectionStore) *application.ConnectionService {
	return application.NewConnectionService(store, server.config)
}

func (server *Server) initConnectionHandler(connectionService *application.ConnectionService, blockService *application.BlockService) *api.ConnectionHandler {
	return api.NewConnectionHandler(connectionService, blockService)
}

func (server *Server) initBlockStore(driver neo4j.Driver) model.BlockStore {
	store := persistance.NewBlockNeo4jStore(driver)
	return store
}

func (server *Server) initBlockService(store model.BlockStore, connectionStore model.ConnectionStore) *application.BlockService {
	return application.NewBlockService(store, connectionStore, server.config)
}
