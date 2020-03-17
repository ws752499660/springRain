package service

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// 保存状态信息
func (t *ServiceSetup) SaveMState(mState MState) (string, error) {

	eventID := "eventAddMState"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将mState对象序列化成为字节数组
	b, err := json.Marshal(mState)
	if err != nil {
		return "", fmt.Errorf("指定的mState对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "addMState", Args: [][]byte{b, []byte(eventID)}}
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

// 用使用者ID和租赁公司名称查找租赁信息
func (t *ServiceSetup) FindMStateByUserIDAndOwnInc(userID, ownInc string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryMStateByUserIDAndOwnInc", Args: [][]byte{[]byte(userID), []byte(ownInc)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// 用机器号码编号查找状态信息
func (t *ServiceSetup) FindMStateByMachineID(machineID string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryMStateByMachineID", Args: [][]byte{[]byte(machineID)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// 在区块链中，现在新增数据与原有的数据都存在
// 但在state DB只保存最新的数据，也就是最新的替代数据。
func (t *ServiceSetup) ModifyMState(mState MState) (string, error) {

	eventID := "eventModifyMState"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将rmState对象序列化成为字节数组
	b, err := json.Marshal(mState)
	if err != nil {
		return "", fmt.Errorf("指定的mState对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "updateMState", Args: [][]byte{b, []byte(eventID)}}
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