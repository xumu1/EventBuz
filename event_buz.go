// @Author：sunhaolong.xumu@bytedance.com
// @Date：2022/8/18

package EventBuz

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"sync"
)

// BusSubscriber 总线订阅者，订阅总线上某个topic的事件。
// 包括的功能有：1. 订阅  2. 取消订阅
type BusSubscriber interface {
	Subscribe(ctx context.Context, topic string, handler EventHandler) error
	UnSubscribe(ctx context.Context, topic string, handler EventHandler) error
}

// BusPublisher 事件发布者，向总线的某个topic发布事件。
type BusPublisher interface {
	Publish(ctx context.Context, topic string, params map[string]interface{}) error
}

// BusController 总线控制台，控制总线的行为。
type BusController interface {
}

// Bus 总线本体接口，包括事件总线的各种功能，都通过bus来完成。
type Bus interface {
	BusSubscriber
	BusPublisher
	BusController
}

// EventBuz 总线的实现类
type EventBuz struct {
	handlers map[string][]*EventHandler
	lock     sync.Mutex
}

// EventHandler 事件的句柄，触发事件
type EventHandler interface {
	Handle(ctx context.Context, formData string) error
}

type EventHandlerFuc func(ctx context.Context, formData string) error

func (f EventHandlerFuc) Handle(ctx context.Context, formData string) error {
	return f(ctx, formData)
}

func NewEventBuz(ctx context.Context) Bus {
	return &EventBuz{
		handlers: make(map[string][]*EventHandler),
		lock:     sync.Mutex{},
	}
}

func (e *EventBuz) Subscribe(ctx context.Context, topic string, handler EventHandler) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.handlers[topic] = append(e.handlers[topic], &handler)
	return nil
}

func (e *EventBuz) UnSubscribe(ctx context.Context, topic string, handler EventHandler) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if _, ok := e.handlers[topic]; ok && len(e.handlers[topic]) > 0 {
		err := e.removeHandler(ctx, topic, handler)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *EventBuz) Publish(ctx context.Context, topic string, params map[string]interface{}) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	param, err := json.Marshal(params)
	if err != nil {
		return err
	}
	if _, ok := e.handlers[topic]; !ok {
		return errors.New("topic not found")
	}
	handlers := e.handlers[topic]
	for _, item := range handlers {
		err = (*item).Handle(ctx, string(param))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *EventBuz) removeHandler(ctx context.Context, topic string, handler EventHandler) error {
	handlers := e.handlers[topic]
	l := len(handlers)
	idx := -1
	for i := range handlers {
		v1 := reflect.ValueOf(*handlers[i])
		v2 := reflect.ValueOf(handler)
		if v1.Type() == v2.Type() && v1.Pointer() == v2.Pointer() {
			idx = i
			break
		}
	}
	if idx == -1 {
		return errors.New("handler not found")
	}
	copy(e.handlers[topic][idx:], e.handlers[topic][idx+1:])
	e.handlers[topic][l-1] = nil
	e.handlers[topic] = e.handlers[topic][:l-1]
	return nil
}
