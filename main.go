//go:generate protoc --go_out=plugins=grpc:. proto/chainnode.proto

package main

import (
	"flag"
	"github.com/SavourDao/savour-core/config"
	wallet2 "github.com/SavourDao/savour-core/rpc/savourrpc/go-savourrpc/wallet"
	"github.com/SavourDao/savour-core/walletdispatcher"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"net"
)

func main() {
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	dispatcher, err := walletdispatcher.New(conf)
	if err != nil {
		log.Error("Setup dispatcher failed", "err", err)
		panic(err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(dispatcher.Interceptor))
	defer grpcServer.GracefulStop()

	wallet2.RegisterWalletServiceServer(grpcServer, dispatcher)

	listen, err := net.Listen("tcp", ":"+conf.Server.Port)
	if err != nil {
		log.Error("net listen failed", "err", err)
		panic(err)
	}

	log.Info("savour dao start success", "port", conf.Server.Port)

	if err := grpcServer.Serve(listen); err != nil {
		log.Error("grpc server serve failed", "err", err)
		panic(err)
	}
}
