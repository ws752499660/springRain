package service

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type RentInfo struct {
	ObjectType string `json:"docType"`

	//租借人信息
	Name       string `json:"Name"`       // 姓名
	Gender     string `json:"Gender"`     // 性别
	Nation     string `json:"Nation"`     // 民族
	EntityID   string `json:"EntityID"`   // 身份证号
	Place      string `json:"Place"`      // 籍贯
	CreditType string `json:"CreditType"` // 信用等级

	//租借信息
	StartDate   string `json:"StartDate"`   // 租借起始日期
	EndDate     string `json:"EndDate"`     // 约定租借终止日期
	ModelName   string `json:"ModelName"`   // 型号名称
	MachineType string `json:"MachineType"` // 设备类型
	IncName     string `json:"IncName"`     // 租赁公司名称
	Price       string `json:"Price"`       // 租赁价格

	//补贴信息
	AllowanceMode   string `json:"AllowanceMode"`   // 补贴类型
	AllowanceAmount string `json:"AllowanceAmount"` // 补贴金额
	AllowanceCheck  string `json:"AllowanceCheck"`  // 是否发放补贴
	AllowanceNo     string `json:"AllowanceNo"`     // 发放账号

	InfoNo string `json:"InfoNo"` // 租赁信息编号

	Historys []RentInfoHistoryItem // 当前RentInfo的历史记录
}

type RentInfoHistoryItem struct {
	TxId     string
	RentInfo RentInfo
}

type MState struct {
	ObjectType string `json:"docType"`

	MachineID string `json:"MachineID"` //机器ID（key）

	CollectDate string `json:"CollectDate"` //信息采集时间
	UserID      string `json:"UserID"`      //当前使用者ID
	OwnInc      string `json:"OwnInc"`      //机器所有的租赁公司

	Location      string `json:"Location"`      //所在地理位置（经纬度） e.g(24.6097159200,118.0910292300)
	EngCode       string `json:"EngCode"`       //发动机故障码
	OilAmountLeft string `json:"OilAmountLeft"` //剩余油量
	TirePressure  string `json:"TirePressure"`  //胎压
	CoolantTemp   string `json:"CoolantTemp"`   //冷却液温度
	OtherInfo     string `json:"OtherInfo"`     //其他信息

	Historys []MStateHistoryItem // 当前MState的历史记录
}

type MStateHistoryItem struct {
	TxId   string
	MState MState
}

type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Printf("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}
