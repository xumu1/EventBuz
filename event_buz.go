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
	WaitAsync(ctx context.Context)
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
	wg       sync.WaitGroup
}

// EventHandler 事件的句柄，触发事件
type EventHandler interface {
	Handle(ctx context.Context, formData string) error
	isOnce() bool
	isAsync() bool
	isTransactional() bool
	transactionLock()
	transactionUnLock()
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
	sync.Mutex
}

func (h *EventHandlerImpl) transactionLock() {
	h.Lock()
}

func (h *EventHandlerImpl) transactionUnLock() {
	h.Unlock()
}

func (h *EventHandlerImpl) Handle(ctx context.Context, formData string) error {
	return h.eventHandlerFuc(ctx, formData)
}

func (h *EventHandlerImpl) isOnce() bool {
	return h.once
}

func (h *EventHandlerImpl) isAsync() bool {
	return h.async
}

func (h *EventHandlerImpl) isTransactional() bool {
	return h.transactional
}

func NewEventBuz(ctx context.Context) Bus {
	return &EventBuz{
		handlers: make(map[string][]EventHandler),
		lock:     sync.Mutex{},
		wg:       sync.WaitGroup{},
	}
}

func (buz *EventBuz) Subscribe(ctx context.Context, topic string, handler EventHandler) error {
	buz.lock.Lock()
	defer buz.lock.Unlock()
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
	for idx, item := range buz.handlers[topic] {
		if item.isOnce() {
			if err := buz.removeHandlerByIndex(ctx, topic, idx); err != nil {
				return err
			}
		}
		if !item.isAsync() {
			return buz.doPublish(ctx, item, string(param))
		}
		buz.wg.Add(1)
		item.transactionLock()
		go buz.doPublish(ctx, item, string(param))
	}
	return nil
}

func (buz *EventBuz) doPublish(ctx context.Context, handler EventHandler, params string) error {
	if handler.isAsync() {
		defer buz.wg.Done()
	}
	if handler.isTransactional() {
		handler.transactionLock()
		defer handler.transactionUnLock()
	}
	err := handler.Handle(ctx, params)
	return err
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

func (buz *EventBuz) WaitAsync(ctx context.Context) {
	buz.wg.Wait()
}
