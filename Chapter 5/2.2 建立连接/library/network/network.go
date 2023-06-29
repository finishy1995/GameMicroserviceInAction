package network

import (
	"ProjectX/library/network/core"
	"ProjectX/library/network/src/tcpgnet"
	"ProjectX/library/network/src/tcpnet"
	"ProjectX/library/routine"
	"sync"
)

// TODO: 为 client 和 server 添加 ID，可以实现关闭某个特定的服务器/客户端

// ChatGPT 建议：
//  1. 对于 Listen 和 Connect 函数，在使用互斥锁后应该在函数的开头记录一个日志，用于追踪函数的调用情况。
//  2. 对于 DestroyAll 函数，在调用之前应该预先关闭所有客户端连接，以避免出现错误。
//  3. 对于 GetInfo 函数，可以考虑提供一个参数用于过滤信息，例如按照网络类型筛选。
//  4. 对于 GetConnNum 函数，可以考虑提供一个参数用于过滤信息，例如按照客户端类型筛选。

var (
	mutex      sync.Mutex
	infoList   map[NetType]*info
	serverList = make([]core.Server, 0, 0)
	clientList = make([]core.Client, 0, 0)
)

func init() {
	infoList = map[NetType]*info{
		TcpNet: {
			client:       func() core.Client { return new(tcpnet.Client) },
			server:       func() core.Server { return new(tcpnet.Server) },
			codecSupport: true,
		},
		TcpGNet: {
			client:       nil,
			server:       func() core.Server { return new(tcpgnet.Server) },
			codecSupport: false,
		},
	}
}

// Listen 使用指定网络类型监听指定端口
func Listen(netType NetType, address string, newAgent core.GetAgent, opts ...core.ServerOption) (core.Server, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if i, ok := infoList[netType]; ok && i.server != nil {
		s := i.server()
		err := s.Start(address, newAgent, opts...)
		if err != nil {
			return nil, err
		}
		serverList = append(serverList, s)
		err = routine.Run(false, s.Run)
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	return nil, ErrUnsupportedNetType
}

// Connect 使用指定网络类型连接指定端口
func Connect(netType NetType, address string, newAgent core.GetAgent, opts ...core.ClientOption) (core.Client, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if i, ok := infoList[netType]; ok && i.client != nil {
		c := i.client()
		err := c.Start(address, newAgent, opts...)
		if err != nil {
			return nil, err
		}
		clientList = append(clientList, c)
		err = routine.Run(false, func() {
			c.Run()
		})
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	return nil, ErrUnsupportedNetType
}

// DestroyAll 摧毁所有服务器客户端
func DestroyAll() {
	mutex.Lock()
	defer mutex.Unlock()
	for _, server := range serverList {
		server.Close()
	}
	for _, client := range clientList {
		client.Close()
	}
	serverList = make([]core.Server, 0, 0)
	clientList = make([]core.Client, 0, 0)
}

// GetConnNum 获取所有连接数
func GetConnNum() (num int) {
	num = 0
	for _, server := range serverList {
		num += server.GetConnNum()
	}
	num += len(clientList)
	//for _, client := range clientList {
	//	if client.IsConnected() {
	//		num++
	//	}
	//}
	return
}

// GetInfo 获取所有允许的网络类型信息
func GetInfo() map[NetType]*info {
	return infoList
}
