package main

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

	InfoNo string `json:"InfoNo"` // 租赁信息编号（key）

	Historys []RentInfoHistoryItem // 当前RentInfo的历史记录
}

type RentInfoHistoryItem struct {
	TxId     string   //该历史变更的交易ID
	RentInfo RentInfo //该历史变更的RentInfo
}
