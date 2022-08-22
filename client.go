// @Author：sunhaolong.xumu@bytedance.com
// @Date：2022/8/19

package EventBuz

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

// Client 负责对网络上的其他 NetworkEventBuz 进行注册句柄、事件发送
type Client struct {
	bus           Bus
	address       string
	path          string
	clientService *ClientService
}

type ClientArg struct {
	Args  string
	topic string
}

func NewClient(address, path string, bus Bus) *Client {
	client := &Client{
		bus:     bus,
		address: address,
		path:    path,
	}
	client.clientService = &ClientService{client: client}
	return client
}

func (c *Client) Start() (err error) {
	if c.clientService.started {
		return nil
	}
	server := rpc.NewServer()
	server.Register(c.clientService)
	server.HandleHTTP(c.path, "/debug"+c.path)
	listener, err := net.Listen("tcp", c.address)
	if err != nil {
		return err
	}
	c.clientService.started = true
	go http.Serve(listener, nil)
	return nil
}

func (c *Client) Stop() {
	if c.clientService.started {
		c.clientService.started = false
	}
}

func (c *Client) Subscribe(topic string, handler *EventHandler, serverAddress string, serverPath string) {
	client, err := rpc.DialHTTPPath("tcp", serverAddress, serverPath)
	if err != nil {
		fmt.Errorf("Call: %v", err)
	}
	args := &SubscribeArg{
		clientAddress: c.address,
		clientPath:    c.path,
		serviceMethod: "ClientService.PushEvent",
		topic:         topic,
	}
	reply := new(bool)
	err = client.Call("ServerService.Register", args, &reply)
	if err != nil {
		fmt.Errorf("Call: %v", err)
	}
	if *reply {
		c.bus.Subscribe(topic, *handler)
	}
}

type ClientService struct {
	client  *Client
	started bool
}

func (s *ClientService) PushEvent(arg *ClientArg, reply *bool) error {
	s.client.bus.Publish(arg.topic, arg.Args)
	*reply = true
	return nil
}
