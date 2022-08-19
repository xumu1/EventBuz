// @Author：sunhaolong.xumu@bytedance.com
// @Date：2022/8/18

package EventBuz

import (
	"context"
	"fmt"
	"testing"
)

func Test_base_function_1(t *testing.T) {
	//ctx := context.Background()
	//topic := "test-topic"
	//buz := NewEventBuz(ctx)
	//_ = buz.Subscribe(ctx, topic, EventHandlerFuc(func(ctx context.Context, formData string) error {
	//	fmt.Println(formData)
	//	return nil
	//}))
	//_ = buz.Publish(ctx, topic, map[string]interface{}{
	//	"a": "aa",
	//	"b": "bb",
	//})
	//_ = buz.UnSubscribe(ctx, topic, EventHandlerFuc(func(ctx context.Context, formData string) error {
	//	fmt.Println(formData)
	//	return nil
	//}))
	//_ = buz.Publish(ctx, topic, map[string]interface{}{
	//	"c": "cc",
	//	"d": "dd",
	//})
}
func Test_base_function_2(t *testing.T) {
	ctx := context.Background()
	topic := "test-topic"
	buz := NewEventBuz(ctx)
	handler := &EventHandlerImpl{
		eventHandlerFuc: func(ctx context.Context, formData string) error {
			fmt.Println(formData)
			return nil
		},
	}
	buz.Subscribe(ctx, topic, handler)
	buz.Publish(ctx, topic, map[string]interface{}{
		"a": "aa",
	})
	buz.UnSubscribe(ctx, topic, handler)
	buz.Publish(ctx, topic, map[string]interface{}{
		"b": "bb",
	})
	onceHandler := &EventHandlerImpl{
		once: true,
		eventHandlerFuc: func(ctx context.Context, formData string) error {
			fmt.Println(formData)
			return nil
		},
	}
	buz.Subscribe(ctx, topic, onceHandler)
	buz.Publish(ctx, topic, map[string]interface{}{
		"c": "cc",
	})
	buz.Publish(ctx, topic, map[string]interface{}{
		"d": "dd",
	})
}
