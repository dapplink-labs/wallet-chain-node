package rpc

import (
	"github.com/SavourDao/savour-core/config"
	"github.com/SavourDao/savour-core/rpc/savourrpc/go-savourrpc/wallet"
	"github.com/SavourDao/savour-core/rpc/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type RpcServer struct {
	Log  logrus.StdLogger
	Conf *config.Config
}

func NewRpcServer(config *config.Config) *RpcServer {
	return &RpcServer{
		Log:  logrus.New(),
		Conf: config,
	}
}

func (this *RpcServer) StartRpcSever() error {
	lis, err := net.Listen("tcp", this.Conf.RpcServer.RpcUrl)
	if err != nil {
		this.Log.Fatal("rpc Listen failed, err:", err.Error())
		return err
	}
	this.Log.Println("start rpc server")
	s := grpc.NewServer()
	wallet.RegisterWalletServiceServer(s, &service.WalletRpcServer{})
	err = s.Serve(lis)
	if err != nil {
		this.Log.Fatal("rpc serve failed,err:", err)
		return err
	}
	this.Log.Println("rpc server started listening on", this.Conf.RpcServer.RpcUrl)
	select {}
}
