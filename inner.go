package main

import "fmt"

const (
	configFile  = "config.yaml"
	initialized = false
	SimpleCC    = "springCC"
	channelName = "kevinkongyixueyuan"
	ordererName = "orderer.kevin.kongyixueyuan.com"
)

func init() {
	fmt.Println("正在启动...")
}
