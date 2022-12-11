package utils

import (
	"encoding/json"
	"os"
	"zyloo.com/zinx/ziface"
)

type GlobalObj struct {
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string

	Version          string
	MaxConn          int
	MaxPackageSize   uint32
	WorkerPoolSize   uint32 // 工作池中的 worker 数量
	MaxWorkerTaskLen uint32 // 每个 worker 所对应的消息队列的长度
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.9",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          3,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	//GlobalObject.Reload()
}
