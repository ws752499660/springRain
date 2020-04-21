package main

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
	TxId   string //该历史变更的交易ID
	MState MState //该历史变更时的MState
}
