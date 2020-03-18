package main

import (
	"encoding/json"
	"fmt"

	"github.com/hqu.edu.cn/springRain/service"
)

func rentInfoTest(serviceSetup service.ServiceSetup) {
	//租赁信息部分

	fmt.Println("====================租赁信息测试=======================")

	//=====设置测试数据=================================

	info1 := service.RentInfo{
		Name:       "张三",
		Gender:     "男",
		Nation:     "汉",
		EntityID:   "342222199311230045",
		Place:      "朝鲜新义州",
		CreditType: "极好",

		StartDate:   "2020-1-1",
		EndDate:     "2020-2-23",
		ModelName:   "R5-3500X",
		MachineType: "高端拖拉机",
		IncName:     "超威半导体",
		Price:       "899",

		AllowanceMode:   "百亿补贴",
		AllowanceAmount: "100",
		AllowanceCheck:  "是",
		AllowanceNo:     "6262626262626",

		InfoNo: "101",
	}

	info2 := service.RentInfo{
		Name:       "李四",
		Gender:     "男",
		Nation:     "汉",
		EntityID:   "342222199310020010",
		Place:      "朝鲜开城",
		CreditType: "极好",

		StartDate:   "2020-1-1",
		EndDate:     "2020-2-23",
		ModelName:   "GTX690",
		MachineType: "大当量核武器",
		IncName:     "英伟达军事工业",
		Price:       "670",

		AllowanceMode:   "闲鱼优惠券",
		AllowanceAmount: "20",
		AllowanceCheck:  "是",
		AllowanceNo:     "626262446262626",

		InfoNo: "102",
	}

	//==========save测试================================

	msg, err := serviceSetup.SaveRentInfo(info1)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}

	msg, err = serviceSetup.SaveRentInfo(info2)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}

	//===========根据身份证号与姓名查询租赁信息=========================
	result, err := serviceSetup.FindRentInfoByEntityIDAndName("342222199310020010", "李四")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var rentInfo service.RentInfo
		json.Unmarshal(result, &rentInfo)
		fmt.Println("根据身份证号与姓名查询信息成功：")
		fmt.Println(rentInfo)
	}

	//==========根据信息编号查询租赁信息(溯源)=================================
	result, err = serviceSetup.FindRentInfoByInfoNo("101")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var rentInfo service.RentInfo
		json.Unmarshal(result, &rentInfo)
		fmt.Println("根据信息编号查询信息成功：")
		fmt.Println(rentInfo)
	}

	// =========修改信息=================================================
	infoc := service.RentInfo{
		Name:       "张三",
		Gender:     "男",
		Nation:     "汉",
		EntityID:   "342222199311230045",
		Place:      "朝鲜新义州",
		CreditType: "较差",

		StartDate:   "2020-1-1",
		EndDate:     "2020-2-23",
		ModelName:   "R5-3500X",
		MachineType: "高端拖拉机",
		IncName:     "超威半导体",
		Price:       "899",

		AllowanceMode:   "百亿补贴",
		AllowanceAmount: "100",
		AllowanceCheck:  "否",
		AllowanceNo:     "6262626262626",

		InfoNo: "101",
	}

	msg, err = serviceSetup.ModifyRentInfo(infoc)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息操作成功, 交易编号为: " + msg)
	}

	result, err = serviceSetup.FindRentInfoByInfoNo("101")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var rentInfo service.RentInfo
		json.Unmarshal(result, &rentInfo)
		fmt.Println("根据信息编号查询信息成功：")
		fmt.Println(rentInfo)
	}
}
