package znet

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Host    string
	TcpPort int
	Name    string

	Version          string
	MaxConn          int
	MaxDataLen       uint32
	WorkerPoolSize   uint32 // 工作池中的 worker 数量
	MaxWorkerTaskLen uint32 // 每个 worker 所对应的消息队列的长度
}

func (cfg *Config) Reload() {
	data, err := os.ReadFile("zynx.json")
	if err != nil {
		fmt.Println("Read File Error, Use Default Config")
		return
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}
}

var config *Config

func init() {
	config = &Config{
		Name:             "AppServer",
		Version:          "V1",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          3,
		MaxDataLen:       4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	config.Reload()
}
