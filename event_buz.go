// @Author：sunhaolong.xumu@bytedance.com
// @Date：2022/8/18

package EventBuz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	handlers map[string][]EventHandler
	lock     sync.Mutex
}

type handlerSetting struct {
	once bool
}

// EventHandler 事件的句柄，触发事件
type EventHandler interface {
	Handle(ctx context.Context, formData string) error
	isOnce(ctx context.Context) bool
	isAsync(ctx context.Context) bool
	isTransactional(ctx context.Context) bool
}

type EventHandlerFuc func(ctx context.Context, formData string) error

func (f EventHandlerFuc) Handle(ctx context.Context, formData string) error {
	return f(ctx, formData)
}

type EventHandlerImpl struct {
	once            bool
	async           bool
	transactional   bool
	eventHandlerFuc EventHandlerFuc
}

func (h *EventHandlerImpl) Handle(ctx context.Context, formData string) error {
	return h.eventHandlerFuc(ctx, formData)
}

func (h *EventHandlerImpl) isOnce(ctx context.Context) bool {
	return h.once
}

func (h *EventHandlerImpl) isAsync(ctx context.Context) bool {
	return h.async
}

func (h *EventHandlerImpl) isTransactional(ctx context.Context) bool {
	return h.transactional
}

func NewEventBuz(ctx context.Context) Bus {
	return &EventBuz{
		handlers: make(map[string][]EventHandler),
		lock:     sync.Mutex{},
	}
}

func (buz *EventBuz) Subscribe(ctx context.Context, topic string, handler EventHandler) error {
	buz.lock.Lock()
	defer buz.lock.Unlock()
	return buz.doSubscribe(ctx, topic, handler, false)
}

func (buz *EventBuz) SubscribeOnce(ctx context.Context, topic string, handler EventHandler) error {
	buz.lock.Lock()
	defer buz.lock.Unlock()
	return buz.doSubscribe(ctx, topic, handler, true)
}

func (buz *EventBuz) doSubscribe(ctx context.Context, topic string, handler EventHandler, once bool) error {
	buz.handlers[topic] = append(buz.handlers[topic], handler)
	return nil
}

func (buz *EventBuz) UnSubscribe(ctx context.Context, topic string, handler EventHandler) error {
	buz.lock.Lock()
	defer buz.lock.Unlock()
	if _, ok := buz.handlers[topic]; ok && len(buz.handlers[topic]) > 0 {
		err := buz.removeHandler(ctx, topic, handler)
		if err != nil {
			return err
		}
	}
	return fmt.Errorf("topic %s not found", topic)
}

func (buz *EventBuz) Publish(ctx context.Context, topic string, params map[string]interface{}) error {
	buz.lock.Lock()
	defer buz.lock.Unlock()
	param, _ := json.Marshal(params)
	if _, ok := buz.handlers[topic]; !ok {
		return fmt.Errorf("handlers in %s topic not found", topic)
	}
	handlers := buz.handlers[topic]
	for idx, item := range handlers {
		if item.isOnce(ctx) {
			if err := buz.removeHandlerByIndex(ctx, topic, idx); err != nil {
				return err
			}
		}
		if err := item.Handle(ctx, string(param)); err != nil {
			return err
		}
	}
	return nil
}

func (buz *EventBuz) removeHandler(ctx context.Context, topic string, handler EventHandler) error {
	handlers := buz.handlers[topic]
	idx := -1
	for i := range handlers {
		v1 := reflect.ValueOf(handlers[i])
		v2 := reflect.ValueOf(handler)
		if v1.Type() == v2.Type() && v1.Pointer() == v2.Pointer() {
			idx = i
			break
		}
	}
	if idx == -1 {
		return errors.New("handler not found")
	}
	return buz.removeHandlerByIndex(ctx, topic, idx)
}

func (buz *EventBuz) removeHandlerByIndex(ctx context.Context, topic string, idx int) error {
	handlers := buz.handlers[topic]
	l := len(handlers)
	if idx == -1 {
		return errors.New("handler not found")
	}
	copy(buz.handlers[topic][idx:], buz.handlers[topic][idx+1:])
	buz.handlers[topic][l-1] = nil
	buz.handlers[topic] = buz.handlers[topic][:l-1]
	return nil
}
