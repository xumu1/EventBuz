// @Author：sunhaolong.xumu@bytedance.com
// @Date：2022/8/18

package EventBuz

import (
	"context"
	"fmt"
	"testing"
)

func Test_base_function(t *testing.T) {
	ctx := context.Background()
	topic := "test-topic"
	buz := NewEventBuz(ctx)
	handler := &MyHandler{}
	_ = buz.Subscribe(ctx, topic, handler)
	_ = buz.Publish(ctx, topic, map[string]interface{}{
		"a": "aa",
		"b": "bb",
	})
	_ = buz.UnSubscribe(ctx, topic, handler)
	_ = buz.Publish(ctx, topic, map[string]interface{}{
		"a": "aa",
		"b": "bb",
	})
}

type MyHandler struct {
	EventHandler
}

func (h MyHandler) Handle(ctx context.Context, formData string) error {
	fmt.Println(formData)
	return nil
}
