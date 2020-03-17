package service

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

//================================租赁信息部分===============================
func (t *ServiceSetup) SaveRentInfo(rentInfo RentInfo) (string, error) {

	eventID := "eventAddRentInfo"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将rentInfo对象序列化成为字节数组
	b, err := json.Marshal(rentInfo)
	if err != nil {
		return "", fmt.Errorf("指定的rentInfo对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "addRentInfo", Args: [][]byte{b, []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

// 用身份证和姓名查找租赁信息
func (t *ServiceSetup) FindRentInfoByEntityIDAndName(entityID, name string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryRentInfoByEntityIDAndName", Args: [][]byte{[]byte(entityID), []byte(name)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// 用信息编号查找租赁信息
func (t *ServiceSetup) FindRentInfoByInfoNo(infoNo string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryRentInfoByInfoNo", Args: [][]byte{[]byte(infoNo)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// 在区块链中，现在新增数据与原有的数据都存在
// 但在state DB只保存最新的数据，也就是最新的替代数据。
func (t *ServiceSetup) ModifyRentInfo(rentInfo RentInfo) (string, error) {

	eventID := "eventModifyRentInfo"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将rentInfo对象序列化成为字节数组
	b, err := json.Marshal(rentInfo)
	if err != nil {
		return "", fmt.Errorf("指定的rentInfo对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "updateRentInfo", Args: [][]byte{b, []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}
