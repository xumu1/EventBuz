// @Author：sunhaolong.xumu@bytedance.com
// @Date：2022/8/19

package EventBuz

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

/*
	NetworkEventBuz 是 EventBuz 的网络升级版本，在原有功能的基础上加上了网络通信的能力
	Client 负责对网络上的其他 NetworkEventBuz 进行注册句柄、事件发送
	Server 负责接受网络上其他 NetworkEventBuz 的注册句柄的请求，将之注册在内部的 Bus 上，并且负责事件的触发。
*/
type NetworkEventBuz struct {
	*Client
	*Server
	bus     Bus
	address string // ip:port of the bus
	path    string // access path of the bus
	service *NetworkEventBuzService
}

func NewNetworkEventBuz(address, path string) *NetworkEventBuz {
	innerBus := NewEventBuz()
	return &NetworkEventBuz{
		Client:  NewClient(address, path, innerBus),
		Server:  NewServer(address, path, innerBus),
		bus:     innerBus,
		address: address,
		path:    path,
		service: &NetworkEventBuzService{},
	}
}

func (b NetworkEventBuz) Start() (err error) {
	clientService := b.Client.clientService
	serverService := b.Server.serverService
	rpcServer := rpc.NewServer()
	rpcServer.RegisterName("ClientService", clientService)
	rpcServer.RegisterName("ServerService", serverService)
	rpcServer.HandleHTTP(b.path, "/debug"+b.path)
	listen, err := net.Listen("tcp", b.address)
	if err != nil {
		err = fmt.Errorf("listen error : %v", err)
	}
	go http.Serve(listen, nil)
	return err
}

func (b NetworkEventBuz) Stop() {
}

// NetworkEventBuzService 用来进行实际业务操作的 serverService
type NetworkEventBuzService struct {
}
