// @Author：sunhaolong.xumu@bytedance.com
// @Date：2022/8/19

package EventBuz

import "testing"

var netTopic string
var netBuz Bus

func init() {
	topic = "test-topic"
	buz = NewEventBuz()
}

func Test_Client(t *testing.T) {
	client := NewClient(":7777", "/_client", NewEventBuz())
	client.Start()
	println(client.clientService.started)
	client.Stop()
	println(client.clientService.started)
}

func Test_Server(t *testing.T) {
	client := NewClient(":7777", "/_client", NewEventBuz())
	server := NewServer(":8888", "/_server", NewEventBuz())
	client.Start()
	server.Start()

}
