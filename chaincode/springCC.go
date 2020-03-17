package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const RENT_DOC_TYPE = "rentObj"
const STATE_DOC_TYPE = "stateObj"

// 保存RentInfo
// args: RentInfo
func PutRentInfo(stub shim.ChaincodeStubInterface, rentinfo RentInfo) ([]byte, bool) {

	rentinfo.ObjectType = RENT_DOC_TYPE

	b, err := json.Marshal(rentinfo)
	if err != nil {
		return nil, false
	}

	// 保存RentInfo状态,将InfoNo做为Key
	err = stub.PutState(rentinfo.InfoNo, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

// 根据租赁信息号码查询信息状态
// args: infoNo
func GetRentInfo(stub shim.ChaincodeStubInterface, InfoNo string) (RentInfo, bool) {
	var rentinfo RentInfo
	// 根据租赁信息号码查询信息状态
	b, err := stub.GetState(InfoNo)
	if err != nil {
		return rentinfo, false
	}

	if b == nil {
		return rentinfo, false
	}

	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &rentinfo)
	if err != nil {
		return rentinfo, false
	}

	// 返回结果
	return rentinfo, true
}

// 根据指定的查询字符串实现富查询
func getRentInfoByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil

}

// 添加信息
// args: rentInfoObject
// 信息编号为 key, RentInfo 为 value
func (t *SpringChaincode) addRentInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var rentinfo RentInfo
	err := json.Unmarshal([]byte(args[0]), &rentinfo)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	// 查重: 信息编号必须唯一
	_, exist := GetRentInfo(stub, rentinfo.InfoNo)
	if exist {
		return shim.Error("要添加的信息编号已存在")
	}

	_, bl := PutRentInfo(stub, rentinfo)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))
}

// 根据身份证号及姓名查询信息
// args: EntityID, name
// 利用富查询实现，需要手动拼接一个给CouchDB的查询字符串
func (t *SpringChaincode) queryRentInfoByEntityIDAndName(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	entityID := args[0]
	name := args[1]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"eduObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"EntityID\":\"%s\", \"Name\":\"%s\"}}", RENT_DOC_TYPE, entityID, name)

	// 查询数据
	result, err := getRentInfoByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据根据身份证号及姓名查询信息发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的根据身份证号及姓名没有查询到相关的信息")
	}
	return shim.Success(result)
}

// 根据信息编号查询详情（并将当前状态保存为历史状态）
// args: infoNo
func (t *SpringChaincode) queryRentInfoByInfoNo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}

	// 根据信息编号查询rentInfo状态
	b, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据信息编号查询信息失败")
	}

	if b == nil {
		return shim.Error("根据信息编号没有查询到相关的信息")
	}

	// 对查询到的状态进行反序列化
	var rentinfo RentInfo
	err = json.Unmarshal(b, &rentinfo)
	if err != nil {
		return shim.Error("反序列化rentInfo信息失败")
	}

	// 获取历史变更数据
	iterator, err := stub.GetHistoryForKey(rentinfo.InfoNo)
	/*
		历史数据查询GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error)
		对于某个我们更改和删除过数据的对象，现在要查询这个对象的更改记录，直接返回区块链上数据（状态数据库上只有最新记录）
	*/
	if err != nil {
		return shim.Error("根据指定的信息编号查询对应的历史变更数据失败")
	}
	defer iterator.Close()

	// 迭代处理
	var historys []RentInfoHistoryItem
	var hisRentInfo RentInfo
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("获取rentInfo的历史变更数据失败")
		}

		var historyItem RentInfoHistoryItem
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisRentInfo)

		if hisData.Value == nil {
			var empty RentInfo
			historyItem.RentInfo = empty
		} else {
			historyItem.RentInfo = hisRentInfo
		}

		historys = append(historys, historyItem)

	}

	historys = append(historys[:len(historys)-1])

	rentinfo.Historys = historys

	// 返回
	result, err := json.Marshal(rentinfo)
	if err != nil {
		return shim.Error("序列化rentinfo信息时发生错误")
	}
	return shim.Success(result)
}

// 根据信息编号更新信息
// args: rentInfoObject
// GetRentInfo除了查询，还将原状态放入历史中返回给调用
func (t *SpringChaincode) updateRentInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var rentinfo RentInfo
	err := json.Unmarshal([]byte(args[0]), &rentinfo)
	if err != nil {
		return shim.Error("反序列化rentInfo信息失败")
	}

	// 根据信息编号查询信息
	result, bl := GetRentInfo(stub, rentinfo.InfoNo)
	if !bl {
		return shim.Error("根据信息编号查询信息时发生错误")
	}

	result.CreditType = rentinfo.CreditType

	result.StartDate = rentinfo.StartDate
	result.EndDate = rentinfo.EndDate
	result.ModelName = rentinfo.ModelName
	result.IncName = rentinfo.IncName
	result.MachineType = rentinfo.MachineType
	result.ModelName = rentinfo.ModelName
	result.Price = rentinfo.Price

	result.AllowanceAmount = rentinfo.AllowanceAmount
	result.AllowanceMode = rentinfo.AllowanceMode
	result.AllowanceCheck = rentinfo.AllowanceCheck

	_, bl = PutRentInfo(stub, result)
	if !bl {
		return shim.Error("保存信息信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息更新成功"))
}

// 根据信息编号删除信息（暂不对外提供）
// args: infoNo
func (t *SpringChaincode) delRentInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	/*var edu Education
	result, bl := GetEduInfo(stub, info.EntityID)
	err := json.Unmarshal(result, &edu)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}*/

	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("删除信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息删除成功"))
}

//=============================================================================================================

//MState部分
// 保存MState
// args: MState
func PutMState(stub shim.ChaincodeStubInterface, mState MState) ([]byte, bool) {

	mState.ObjectType = STATE_DOC_TYPE
	fmt.Print(mState)
	b, err := json.Marshal(mState)
	if err != nil {
		return nil, false
	}

	// 保存mState状态,将MachineID做为Key
	err = stub.PutState(mState.MachineID, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

// 根据机器编号查询信息状态
// args: machineID
func GetMState(stub shim.ChaincodeStubInterface, machineID string) (MState, bool) {
	var mState MState
	// 根据机器编号查询信息状态
	b, err := stub.GetState(machineID)
	if err != nil {
		return mState, false
	}

	if b == nil {
		return mState, false
	}

	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &mState)
	if err != nil {
		return mState, false
	}

	// 返回结果
	return mState, true
}

// 根据指定的查询字符串实现富查询
func getMStateByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil

}

// 添加信息
// args: mStateObject
// 机器编号为 key, MState 为 value
func (t *SpringChaincode) addMState(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var mState MState
	err := json.Unmarshal([]byte(args[0]), &mState)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	// 查重: 信息编号必须唯一
	_, exist := GetMState(stub, mState.MachineID)
	if exist {
		return shim.Error("要添加的机器编号已存在")
	}

	_, bl := PutMState(stub, mState)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))
}

// 根据使用者ID及租赁公司名称查询信息
// args: UserID, OwnInc
// 利用富查询实现，需要手动拼接一个给CouchDB的查询字符串
func (t *SpringChaincode) queryMStateByUserIDAndOwnInc(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	userID := args[0]
	ownInc := args[1]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"eduObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"UserID\":\"%s\", \"OwnInc\":\"%s\"}}", STATE_DOC_TYPE, userID, ownInc)

	// 查询数据
	result, err := getMStateByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据使用者ID及租赁公司名称查询信息发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的使用者ID及租赁公司名称没有查询到相关的信息")
	}
	return shim.Success(result)
}

// 根据机器编号查询详情（并将当前状态保存为历史状态）
// args: machineID
func (t *SpringChaincode) queryMStateByMachineID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}

	// 根据信息编号查询MState状态
	b, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据信息编号查询信息失败")
	}

	if b == nil {
		return shim.Error("根据信息编号没有查询到相关的信息")
	}

	// 对查询到的状态进行反序列化
	var mState MState
	err = json.Unmarshal(b, &mState)
	if err != nil {
		return shim.Error("反序列化MState信息失败")
	}

	// 获取历史变更数据
	iterator, err := stub.GetHistoryForKey(mState.MachineID)
	/*
		历史数据查询GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error)
		对于某个我们更改和删除过数据的对象，现在要查询这个对象的更改记录，直接返回区块链上数据（状态数据库上只有最新记录）
	*/
	if err != nil {
		return shim.Error("根据指定的信息编号查询对应的历史变更数据失败")
	}
	defer iterator.Close()

	// 迭代处理
	var historys []MStateHistoryItem
	var hisMState MState
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("获取MState的历史变更数据失败")
		}

		var historyItem MStateHistoryItem
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisMState)

		if hisData.Value == nil {
			var empty MState
			historyItem.MState = empty
		} else {
			historyItem.MState = hisMState
		}

		historys = append(historys, historyItem)

	}

	historys = append(historys[:len(historys)-1])

	mState.Historys = historys

	// 返回
	result, err := json.Marshal(mState)
	if err != nil {
		return shim.Error("序列化mState信息时发生错误")
	}
	return shim.Success(result)
}

// 根据机器编号更新信息
// args: mStateObject
func (t *SpringChaincode) updateMState(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var mState MState
	err := json.Unmarshal([]byte(args[0]), &mState)
	if err != nil {
		return shim.Error("反序列化mState信息失败")
	}

	// 根据机器编号查询信息
	result, bl := GetMState(stub, mState.MachineID)
	if !bl {
		return shim.Error("根据机器编号查询信息时发生错误")
	}

	result.CollectDate = mState.CollectDate
	result.UserID = mState.UserID
	result.OwnInc = mState.OwnInc

	result.Location = mState.Location
	result.EngCode = mState.EngCode
	result.OilAmountLeft = mState.OilAmountLeft
	result.TirePressure = mState.TirePressure
	result.CoolantTemp = mState.CoolantTemp
	result.OtherInfo = mState.OtherInfo

	_, bl = PutMState(stub, result)
	if !bl {
		return shim.Error("保存信息信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息更新成功"))
}

// 根据机器编号删除状态（暂不对外提供）
// args: machineID
func (t *SpringChaincode) delMState(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	/*var edu Education
	result, bl := GetEduInfo(stub, info.EntityID)
	err := json.Unmarshal(result, &edu)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}*/

	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("删除信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息删除成功"))
}
