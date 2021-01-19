/*
package nats 基于nats消息组件，提供客户端连接、消息发布和消息订阅等接口.
for example:
	// first init connect
	nats.Init()

	// Subscribe
	nats.Subscribe("hello", func(msg *nats.Msg) {
		fmt.Printf("Received a message: %+v\n", msg)
	})

	// Publish
	nats.Publish("hello", []byte("world"))
*/
package nats
