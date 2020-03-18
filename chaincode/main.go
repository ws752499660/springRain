package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SpringChaincode struct {
}

func (t *SpringChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

func (t *SpringChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()

	if fun == "addRentInfo" {
		return t.addRentInfo(stub, args) // 添加租赁信息
	} else if fun == "queryRentInfoByEntityIDAndName" {
		return t.queryRentInfoByEntityIDAndName(stub, args) // 根据身份证号及姓名查询信息
	} else if fun == "queryRentInfoByInfoNo" {
		return t.queryRentInfoByInfoNo(stub, args) // 根据信息编号查询详情
	} else if fun == "updateRentInfo" {
		return t.updateRentInfo(stub, args) // 根据信息编号更新信息
	} else if fun == "delRentInfo" {
		return t.delRentInfo(stub, args) // 根据信息编号删除信息
	} else if fun == "addMState" {
		return t.addMState(stub, args) // 添加农机状态
	} else if fun == "queryMStateByUserIDAndOwnInc" {
		return t.queryMStateByUserIDAndOwnInc(stub, args) // 根据用户ID和租赁公司查询信息
	} else if fun == "queryMStateByMachineID" {
		return t.queryMStateByMachineID(stub, args) // 根据机器编号查询详情
	} else if fun == "updateMState" {
		return t.updateMState(stub, args) // 根据机器编号更新信息
	} else if fun == "delMState" {
		return t.delMState(stub, args) // 根据机器编号删除信息
	}

	return shim.Error("指定的函数名称错误")

}

func main() {
	err := shim.Start(new(SpringChaincode))
	if err != nil {
		fmt.Printf("启动SpringChaincode时发生错误: %s", err)
	}
}
