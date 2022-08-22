// @Author：sunhaolong.xumu@bytedance.com
// @Date：2022/8/19

package EventBuz

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

// Server 负责接受网络上其他 NetworkEventBuz 的注册句柄的请求，将之注册在内部的 Bus 上
type Server struct {
	bus           Bus
	address       string
	path          string
	subscribers   map[string][]*SubscribeArg
	serverService *ServerService
}

type SubscribeArg struct {
	clientAddress string
	clientPath    string
	serviceMethod string // rpc 回调的方法
	topic         string
}

type ServerService struct {
	server  *Server
	started bool
}

func NewServer(address string, path string, eventbus Bus) *Server {
	server := &Server{
		bus:         eventbus,
		address:     address,
		path:        path,
		subscribers: map[string][]*SubscribeArg{},
	}
	server.serverService = &ServerService{server: server, started: false}
	return server
}

func (s *Server) rpcCallback(subscribeArg *SubscribeArg) EventHandler {
	return &EventHandlerImpl{
		eventHandlerFuc: func(formData string) error {
			client, err := rpc.DialHTTPPath("tcp", subscribeArg.clientAddress, subscribeArg.clientPath)
			defer client.Close()
			if err != nil {
				fmt.Errorf("dialing: %v", err)
			}
			clientArg := &ClientArg{
				Args:  formData,
				topic: subscribeArg.topic,
			}
			var reply bool
			err = client.Call(subscribeArg.serviceMethod, clientArg, &reply)
			if err != nil {
				fmt.Errorf("Call: %v", err)
			}
			return nil
		},
	}
}

// 查询是否有过订阅
func (s *Server) hasClientSubscribers(arg *SubscribeArg) bool {
	if topicSubscribers, ok := s.subscribers[arg.topic]; ok {
		for _, subscriber := range topicSubscribers {
			if *subscriber == *arg {
				return true
			}
		}
	}
	return false
}

func (s *Server) Start() (err error) {
	if s.serverService.started {
		return
	}
	rpcServer := rpc.NewServer()
	rpcServer.Register(s.serverService)
	rpcServer.HandleHTTP(s.path, "/debug"+s.path)
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	s.serverService.started = true
	go http.Serve(listener, nil)
	return nil
}

func (s *Server) Stop() {
	if s.serverService.started {
		s.serverService.started = false
	}
}

func (s ServerService) Register(arg *SubscribeArg, reply *bool) error {
	subscribers := s.server.subscribers
	if !s.server.hasClientSubscribers(arg) {
		callback := s.server.rpcCallback(arg)
		// 只能接收EventHandler
		s.server.bus.Subscribe(arg.topic, callback)
		var topicSubscribers []*SubscribeArg
		if _, ok := subscribers[arg.topic]; ok {
			topicSubscribers = subscribers[arg.topic]
			topicSubscribers = append(topicSubscribers, arg)
		} else {
			topicSubscribers = []*SubscribeArg{arg}
		}
		subscribers[arg.topic] = topicSubscribers
	}
	*reply = true
	return nil
}
