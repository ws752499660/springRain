package main

import (
	"fmt"
	"os"

	"github.com/hqu.edu.cn/springRain/sdkInit"
	"github.com/hqu.edu.cn/springRain/service"
)

func main() {

	//======创建SDK，加入通道============
	initInfo := &sdkInit.InitInfo{

		ChannelID:     channelName,
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/hqu.edu.cn/springRain/fixtures/artifacts/channel.tx",

		OrgAdmin:       "Admin",
		OrgName:        "Org1",
		OrdererOrgName: ordererName,

		ChaincodeID:     SimpleCC,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/hqu.edu.cn/springRain/chaincode/",
		UserName:        "User1",
	}

	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer sdk.Close()

	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//======实例化链码并创建通道客户端====================
	channelClient, err := sdkInit.InstallAndInstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(channelClient)

	//======启动service===============================
	serviceSetup := service.ServiceSetup{
		ChaincodeID: SimpleCC,
		Client:      channelClient,
	}

	//租赁信息部分
	rentInfoTest(serviceSetup)
	//机器状态部分
	mStateTest(serviceSetup)
	//缓冲上链部分
	mbufferTest(serviceSetup)
}
