package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/bensmile/wekamakuta/api"
	db "github.com/bensmile/wekamakuta/db/sqlc"
	"github.com/bensmile/wekamakuta/db/util"
	_ "github.com/bensmile/wekamakuta/doc/statik"
	"github.com/bensmile/wekamakuta/gapi"
	"github.com/bensmile/wekamakuta/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: %w", err)
	}

	if err := server.Start(config.HttpServerAddress); err != nil {
		log.Fatal("cannot start the server:", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("could not create grpc server", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("cannot start grpc server :", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("could not create grpc server", err)
	}
	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOptions)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server); err != nil {
		log.Fatal("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()

	if err != nil {
		log.Fatal("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())

	if err := http.Serve(listener, mux); err != nil {
		log.Fatal("cannot start HTTP gateway server :", err)
	}
}
