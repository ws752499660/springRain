package web

import (
	"fmt"
	"net/http"

	"github.com/hqu.edu.cn/springRain/web/controller"
)

// 启动Web服务并指定路由信息
func WebStart(app controller.Application) {

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 指定路由信息(匹配请求)
	http.HandleFunc("/", app.LoginView)
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/loginout", app.LoginOut)

	http.HandleFunc("/index", app.Index)
	http.HandleFunc("/help", app.Help)

	//=========================================

	http.HandleFunc("/addRentInfo", app.AddRentShow) // 显示添加信息页面
	http.HandleFunc("/addRent", app.AddRent)         // 提交信息请求

	http.HandleFunc("/queryRentByEntityIDAndName", app.QueryRentByEntityIDAndNameShow) // 转至根据证书编号与姓名查询信息页面
	http.HandleFunc("/findRentByEntityIDAndName", app.FindRentInfoByEntityIDAndName)   // 根据证书编号与姓名查询信息

	http.HandleFunc("/queryRentByInfoNo", app.QueryRentByInfoNo) // 转至信息编号查询租赁信息页面
	http.HandleFunc("/findRentByInfoNo", app.FindRentByInfoNo)   // 根据信息编号查询租赁信息

	http.HandleFunc("/modifyRentInputShow", app.ModifyRentInputShow) //修改页面输入身份证号与姓名页面
	http.HandleFunc("/modifyRentShow", app.ModifyRentShow)           // 修改信息页面
	http.HandleFunc("/modifyRent", app.ModifyRent)                   //  修改信息

	//============================================

	http.HandleFunc("/updateMStateShow", app.UpdateMStateShow) // 显示更新农机状态页面
	http.HandleFunc("/updateMState", app.UpdateMState)         // 提交更新请求

	http.HandleFunc("/queryMStateByUserIDAndOwnInc", app.QueryMStateByUserIDAndOwnIncShow) // 转至根据使用者ID和所在公司查询农机状态页面
	http.HandleFunc("/findMStateByUserIDAndOwnInc", app.FindMStateByUserIDAndOwnInc)       // 根据使用者ID和所在公司查询信息

	http.HandleFunc("/queryMStateByMachineID", app.QueryMStateByMachineID) // 转至根据农机的机器ID查询农机状态页面
	http.HandleFunc("/findMStateByMachineID", app.FindMStateByMachineID)   // 根据农机的机器ID查询状态信息

	//=============================================

	http.HandleFunc("/bufferUpdate", app.BufferUpdate) //处理缓存过来的农机状态更新

	fmt.Println("启动Web服务, 监听端口号为: 9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("Web服务启动失败: %v", err)
	}

}
