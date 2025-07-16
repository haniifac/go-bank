package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/haniifac/simplebank/api"
	db "github.com/haniifac/simplebank/db/sqlc"
	"github.com/haniifac/simplebank/gapi"
	"github.com/haniifac/simplebank/pb"
	"github.com/haniifac/simplebank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	runGrpcServer(config, store)
	// runGinServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	grpcServer := grpc.NewServer()

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create gRPC server: ", err)
	}

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create gRPC listener: ", err)
	}

	log.Printf("gRPC server is running on %s", config.GRPCServerAddress)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server: ", err)
	}
}
