package main

import (
	"encoding/json"
	"fmt"

	"github.com/hqu.edu.cn/springRain/service"
)

func mStateTest(serviceSetup service.ServiceSetup) {

	//=====设置测试数据=================================

	state1 := service.MState{
		MachineID: "201",

		CollectDate: "2020-1-1 10:30",
		UserID:      "U101",
		OwnInc:      "AMD INC.",

		Location:      "(24.6097159200,118.0910292300)",
		EngCode:       "None",
		OilAmountLeft: "50",
		TirePressure:  "[243,244,246,250]",
		CoolantTemp:   "90",
		OtherInfo:     "{'RoutineMaintenanceCheck':'True','WorkingType':'Normal'}",
	}

	state2 := service.MState{
		MachineID: "202",

		CollectDate: "2020-3-1 10:30",
		UserID:      "U102",
		OwnInc:      "NVIDA Industry INC.",

		Location:      "(25.6097159200,118.0910292300)",
		EngCode:       "None",
		OilAmountLeft: "500",
		TirePressure:  "[243,244,246,250]",
		CoolantTemp:   "96",
		OtherInfo:     "{'SafeCheck':'True','WorkingType':'Normal','Boomable':'True'}",
	}

	//==========save测试================================

	msg, err := serviceSetup.SaveMState(state1)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}

	msg, err = serviceSetup.SaveMState(state2)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}

	//===========根据使用者编号与租赁公司名称查询状态信息=========================
	result, err := serviceSetup.FindMStateByUserIDAndOwnInc("U101", "AMD INC.")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var mState service.MState
		json.Unmarshal(result, &mState)
		fmt.Println("根据使用者编号与租赁公司名称查询信息成功：")
		fmt.Println(mState)
	}

	//==========根据机器编号查询租赁信息(溯源)=================================
	result, err = serviceSetup.FindMStateByMachineID("202")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var mState service.MState
		json.Unmarshal(result, &mState)
		fmt.Println("根据机器编号查询信息成功：")
		fmt.Println(mState)
	}

	// =========修改状态=================================================
	statec := service.MState{
		MachineID: "202",

		CollectDate: "2020-3-1 10:30",
		UserID:      "U102",
		OwnInc:      "NVIDA Industry INC.",

		Location:      "(25.6097159200,118.0910292300)",
		EngCode:       "404",
		OilAmountLeft: "400",
		TirePressure:  "[243,244,246,250]",
		CoolantTemp:   "96000",
		OtherInfo:     "{'SafeCheck':'False','WorkingType':'WillBoom','Boomable':'True'}",
	}
	msg, err = serviceSetup.ModifyMState(statec)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息操作成功, 交易编号为: " + msg)
	}

	result, err = serviceSetup.FindMStateByMachineID("202")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var mState service.MState
		json.Unmarshal(result, &mState)
		fmt.Println("根据信息编号查询信息成功：")
		fmt.Println(mState)
	}
}
