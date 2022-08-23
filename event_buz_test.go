// @Author：sunhaolong.xumu@bytedance.com
// @Date：2022/8/18

package EventBuz

import (
	"testing"
)

var topic string
var buz Bus

func init() {
	topic = "test-topic"
	buz = NewEventBuz()
}

func Test_AsyncHandler(t *testing.T) {
	//myHandler := &EventHandlerImpl{
	//	once:  false,
	//	async: false,
	//	eventHandlerFuc: func(formData string) error {
	//		fmt.Println(formData)
	//		time.Sleep(3 * time.Second)
	//		return nil
	//	},
	//}
	//buz.Subscribe(topic, myHandler)
	//buz.Publish(topic, map[string]interface{}{
	//	"a": "aa",
	//})
	//fmt.Printf("now time is: %v\n", time.Now())
	//buz.Publish(topic, map[string]interface{}{
	//	"b": "bb",
	//})
	//fmt.Printf("now time is: %v\n", time.Now())
	//buz.WaitAsync()
	//fmt.Printf("now time is: %v\n", time.Now())
	//fmt.Println("it is fine.")

}

func Test_base_function_1(t *testing.T) {
	//ctx := context.Background()
	//topic := "test-topic"
	//buz := NewEventBuz(ctx)
	//_ = buz.Subscribe( topic, EventHandlerFuc(func( formData string) error {
	//	fmt.Println(formData)
	//	return nil
	//}))
	//_ = buz.Publish( topic, map[string]interface{}{
	//	"a": "aa",
	//	"b": "bb",
	//})
	//_ = buz.UnSubscribe( topic, EventHandlerFuc(func( formData string) error {
	//	fmt.Println(formData)
	//	return nil
	//}))
	//_ = buz.Publish( topic, map[string]interface{}{
	//	"c": "cc",
	//	"d": "dd",
	//})
}
func Test_base_function_2(t *testing.T) {
	//handler := &EventHandlerImpl{
	//	eventHandlerFuc: func(formData string) error {
	//		fmt.Println(formData)
	//		return nil
	//	},
	//}
	//buz.Subscribe(topic, handler)
	//buz.Publish(topic, map[string]interface{}{
	//	"a": "aa",
	//})
	//buz.UnSubscribe(topic, handler)
	//buz.Publish(topic, map[string]interface{}{
	//	"b": "bb",
	//})
	//onceHandler := &EventHandlerImpl{
	//	once: true,
	//	eventHandlerFuc: func(formData string) error {
	//		fmt.Println(formData)
	//		return nil
	//	},
	//}
	//buz.Subscribe(topic, onceHandler)
	//buz.Publish(topic, map[string]interface{}{
	//	"c": "cc",
	//})
	//buz.Publish(topic, map[string]interface{}{
	//	"d": "dd",
	//})
}
