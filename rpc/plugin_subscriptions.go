package rpc

import (
	"context"
	"reflect"
	"github.com/ethereum/go-ethereum/log"
)


func isChanType(t reflect.Type) bool {
	// Pointers to channels are weird, but whatever
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// Make sure we have a channel
	if t.Kind() != reflect.Chan {
		return false
	}
	// Make sure it is a receivable channel
	return (t.ChanDir() & reflect.RecvDir) == reflect.RecvDir
}

func isChanPubsub(methodType reflect.Type) bool {
	if methodType.NumIn() < 2 || methodType.NumOut() != 2 {
		return false
	}
	return isContextType(methodType.In(1)) &&
		isChanType(methodType.Out(0)) &&
		isErrorType(methodType.Out(1))
}

func callbackifyChanPubSub(receiver, fn reflect.Value) *callback {
	c := &callback{rcvr: receiver, errPos: 1, isSubscribe: true}
	fntype := fn.Type()
	// Skip receiver and context.Context parameter (if present).
	firstArg := 0
	if c.rcvr.IsValid() {
		firstArg++
	}
	if fntype.NumIn() > firstArg && fntype.In(firstArg) == contextType {
		c.hasCtx = true
		firstArg++
	}
	// Add all remaining parameters.
	c.argTypes = make([]reflect.Type, fntype.NumIn()-firstArg)
	for i := firstArg; i < fntype.NumIn(); i++ {
		c.argTypes[i-firstArg] = fntype.In(i)
	}

	retFnType := reflect.FuncOf(append([]reflect.Type{receiver.Type(), contextType}, c.argTypes...), []reflect.Type{subscriptionType, errorType}, false)

// // What follows uses reflection to construct a dynamically typed function equivalent to:
// func(receiver <T>, cctx context.Context, args ...<T>) (rpc.Subscription, error) {
// 	notifier, supported := NotifierFromContext(cctx)
// 	if !supported { return Subscription{}, ErrNotificationsUnsupported}
// 	ctx, cancel := context.WithCancel(context.Background())
// 	ch, err := fn()
// 	if err != nil { return Subscription{}, err }
// 	rpcSub := notifier.CreateSubscription()
// 	go func() {
// 		select {
// 		case v, ok := <- ch:
// 			if !ok { return }
// 			notifier.Notify(rpcSub.ID, v)
// 		case <-rpcSub.Err():
// 			cancel()
// 			return
// 		case <-notifier.Closed():
// 			cancel()
// 			return
// 		}
// 	}()
// 	return rpcSub, nil
// }
//

	c.fn = reflect.MakeFunc(retFnType, func(args []reflect.Value) ([]reflect.Value) {
		notifier, supported := NotifierFromContext(args[1].Interface().(context.Context))
		if !supported {
			return []reflect.Value{reflect.Zero(subscriptionType), reflect.ValueOf(ErrNotificationsUnsupported)}
		}
		ctx, cancel := context.WithCancel(context.Background())
		args[1] = reflect.ValueOf(ctx)
		out := fn.Call(args)
		if !out[1].IsNil() {
			// This amounts to: if err != nil { return nil, err }
			return []reflect.Value{reflect.Zero(subscriptionType), out[1]}
		}
		// Geth's provided context is done once we've returned the subscription id.
		// This new context will cancel when the notifier closes.

		rpcSub := notifier.CreateSubscription()
		go func() {
			defer log.Info("Plugin subscription goroutine closed")
			selectCases := []reflect.SelectCase{
				{Dir: reflect.SelectRecv, Chan: out[0]},
				{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(rpcSub.Err())},
				{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(notifier.Closed())},
			}
			for {
				chosen, val, recvOK := reflect.Select(selectCases)
				switch chosen {
				case 0: // val, ok := <-ch
					if !recvOK {
						return
					}
					notifier.Notify(rpcSub.ID, val.Interface())
				case 1:
					cancel()
					return
				case 2:
					cancel()
					return
				}
			}
		}()
		return []reflect.Value{reflect.ValueOf(*rpcSub), reflect.Zero(errorType)}
	})
	return c
}

func pluginExtendedCallbacks(callbacks map[string]*callback, receiver reflect.Value) {
	typ := receiver.Type()
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		if method.PkgPath != "" {
			continue // method not exported
		}
		if method.Name == "Timer" {
			methodType := method.Func.Type()
			log.Info("Timer method", "in", methodType.NumIn(), "out", methodType.NumOut(), "contextType", isContextType(methodType.In(1)), "chanType", isChanType(methodType.Out(0)), "chandir", methodType.Out(0).ChanDir() & reflect.RecvDir == reflect.RecvDir, "errorType", isErrorType(methodType.Out(1)))
		}
		if isChanPubsub(method.Type) {
			cb := callbackifyChanPubSub(receiver, method.Func)
			name := formatName(method.Name)
			callbacks[name] = cb
			log.Info("Added chanPubsub", "name", name, "args", cb.argTypes)
		}
	}
}
