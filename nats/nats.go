package nats

import (
	"errors"
	"sync"

	"go.uber.org/zap"
)

var (
	conn                 *nats.Conn
	once                 sync.Once
	ErrConnectionUnInit  = errors.New("连接未初始化")
	ErrInvalidConnection = errors.New("无可用Nats连接")
)

var Client = globalConnection

func globalConnection() *nats.Conn {
	return conn
}

func Init() (err error) {
	once.Do(func() {
		err = newConnect(C.Host, C.Username,
			C.Password, C.ClientName)
		if err != nil {
			zap.L().Error("new nats connect err", zap.Error(err))
		}
	})

	return
}

// 初始化nats conn
//  host server地址，以逗号分隔
//  username 用户名
//  password 密码
//  clientNamePrefix 客户端名称前缀
func newConnect(host, username, password, clientNamePrefix string) (err error) {
	conn, err = nats.Connect(host, nats.UserInfo(username, password), nats.Name(clientNamePrefix),
		nats.DisconnectErrHandler(OnDisconnection), nats.ReconnectHandler(OnReconnection))
	if err != nil {
		zap.L().Error("connect to nats-server err", zap.Error(err))
		return err
	}

	return
}

// 订阅topic
//  topic 主题
//  callback 订阅回调
func Subscribe(topic string, callback nats.MsgHandler) error {
	if conn == nil {
		zap.L().Error("Nats must be initialized")
		return ErrConnectionUnInit
	}
	if conn.IsClosed() || conn.IsDraining() {
		zap.L().Error("no valid Nats connection")
		return ErrInvalidConnection
	}
	if _, err := conn.Subscribe(topic, callback); err != nil {
		zap.L().Error("subscribe Nats topic err", zap.Error(err))
		return err
	}

	return nil
}

// 发布消息
//  topic 主题
//  data 消息内容
func Publish(topic string, data []byte) error {
	if conn == nil {
		zap.L().Error("Nats must be initialized")
		return ErrConnectionUnInit
	}
	if conn.IsClosed() || conn.IsDraining() {
		zap.L().Error("no valid Nats connection")
		return ErrInvalidConnection
	}

	if err := conn.Publish(topic, data); err != nil {
		zap.L().Error("publish Nats topic err", zap.Error(err))
		return err
	}

	return nil
}

// 客户端重新连接
func OnReconnection(conn *nats.Conn) {
	zap.L().Info("reconnected to nats-server")
}

// 客户端断开连接
func OnDisconnection(conn *nats.Conn, err error) {
	zap.L().Error("disconnect to nats-server", zap.Error(err))
}
