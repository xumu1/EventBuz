// @Author：sunhaolong.xumu@bytedance.com
// @Date：2022/8/18

package EventBuz

import (
	"context"
	"fmt"
	"testing"
)

func Test_base_function_1(t *testing.T) {
	ctx := context.Background()
	topic := "test-topic"
	buz := NewEventBuz(ctx)
	_ = buz.Subscribe(ctx, topic, EventHandlerFuc(func(ctx context.Context, formData string) error {
		fmt.Println(formData)
		return nil
	}))
	_ = buz.Publish(ctx, topic, map[string]interface{}{
		"a": "aa",
		"b": "bb",
	})
	_ = buz.UnSubscribe(ctx, topic, EventHandlerFuc(func(ctx context.Context, formData string) error {
		fmt.Println(formData)
		return nil
	}))
	_ = buz.Publish(ctx, topic, map[string]interface{}{
		"c": "cc",
		"d": "dd",
	})
}
func Test_base_function_2(t *testing.T) {
	ctx := context.Background()
	topic := "test-topic"
	buz := NewEventBuz(ctx)
	fun1 := EventHandlerFuc(func(ctx context.Context, formData string) error {
		fmt.Println(formData)
		return nil
	})
	_ = buz.Subscribe(ctx, topic, fun1)
	_ = buz.Publish(ctx, topic, map[string]interface{}{
		"a": "aa",
		"b": "bb",
	})
	_ = buz.UnSubscribe(ctx, topic, fun1)
	_ = buz.Publish(ctx, topic, map[string]interface{}{
		"c": "cc",
		"d": "dd",
	})
}
