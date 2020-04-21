package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hqu.edu.cn/springRain/service"
)

// 显示添加信息页面
func (app *Application) AddRentShow(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "addRentInfo.html", data)
}

// 添加信息
func (app *Application) AddRent(w http.ResponseWriter, r *http.Request) {

	rent := service.RentInfo{
		Name:            r.FormValue("name"),
		Gender:          r.FormValue("gender"),
		Nation:          r.FormValue("nation"),
		EntityID:        r.FormValue("entityID"),
		Place:           r.FormValue("place"),
		CreditType:      r.FormValue("creditType"),
		StartDate:       r.FormValue("startDate"),
		EndDate:         r.FormValue("endDate"),
		ModelName:       r.FormValue("modelName"),
		MachineType:     r.FormValue("machineType"),
		IncName:         r.FormValue("incName"),
		Price:           r.FormValue("price"),
		AllowanceMode:   r.FormValue("allowanceMode"),
		AllowanceAmount: r.FormValue("allowanceAmount"),
		AllowanceCheck:  r.FormValue("allowanceCheck"),
		AllowanceNo:     r.FormValue("allowanceNo"),
		InfoNo:          r.FormValue("infoNo"),
	}

	app.Setup.SaveRentInfo(rent)

	r.Form.Set("entityID", rent.EntityID)
	r.Form.Set("name", rent.Name)
	app.FindRentInfoByEntityIDAndName(w, r)
}

func (app *Application) QueryRentByEntityIDAndNameShow(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
		IsQuery     bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		IsQuery:     true,
	}
	ShowView(w, r, "queryRentByEntityIDAndName.html", data)
}

// 根据身份证号与姓名查询信息
func (app *Application) FindRentInfoByEntityIDAndName(w http.ResponseWriter, r *http.Request) {
	entityID := r.FormValue("entityID")
	name := r.FormValue("name")
	result, err := app.Setup.FindRentInfoByEntityIDAndName(entityID, name)
	var rent = service.RentInfo{}
	json.Unmarshal(result, &rent)

	fmt.Println("根据身份证号与姓名查询信息成功：")
	fmt.Println(rent)

	data := &struct {
		Rent        service.RentInfo
		CurrentUser User
		Msg         string
		Flag        bool
		History     bool
	}{
		Rent:        rent,
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		History:     false,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "queryRentResult.html", data)
}

func (app *Application) QueryRentByInfoNo(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "queryRentByInfoNo.html", data)
}

// 根据身份证号码查询信息
func (app *Application) FindRentByInfoNo(w http.ResponseWriter, r *http.Request) {
	infoNo := r.FormValue("infoNo")
	result, err := app.Setup.FindRentInfoByInfoNo(infoNo)
	var rent = service.RentInfo{}
	json.Unmarshal(result, &rent)

	data := &struct {
		Rent        service.RentInfo
		CurrentUser User
		Msg         string
		Flag        bool
		History     bool
	}{
		Rent:        rent,
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		History:     true,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "queryRentResult.html", data)
}

func (app *Application) ModifyRentInputShow(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
		IsQuery     bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		IsQuery:     false,
	}
	ShowView(w, r, "queryRentByEntityIDAndName.html", data)
}

// 显示修改/添加新租赁信息
func (app *Application) ModifyRentShow(w http.ResponseWriter, r *http.Request) {
	// 根据身份证号与姓名查询信息
	entityID := r.FormValue("entityID")
	name := r.FormValue("name")
	result, err := app.Setup.FindRentInfoByEntityIDAndName(entityID, name)

	var rent = service.RentInfo{}
	json.Unmarshal(result, &rent)

	data := &struct {
		Rent        service.RentInfo
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		Rent:        rent,
		CurrentUser: cuser,
		Flag:        true,
		Msg:         "",
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "modifyRent.html", data)
}

// 修改/添加新租赁信息
func (app *Application) ModifyRent(w http.ResponseWriter, r *http.Request) {
	rent := service.RentInfo{
		Name:            r.FormValue("name"),
		Gender:          r.FormValue("gender"),
		Nation:          r.FormValue("nation"),
		EntityID:        r.FormValue("entityID"),
		Place:           r.FormValue("place"),
		CreditType:      r.FormValue("creditType"),
		StartDate:       r.FormValue("startDate"),
		EndDate:         r.FormValue("endDate"),
		ModelName:       r.FormValue("modelName"),
		MachineType:     r.FormValue("machineType"),
		IncName:         r.FormValue("incName"),
		Price:           r.FormValue("price"),
		AllowanceMode:   r.FormValue("allowanceMode"),
		AllowanceAmount: r.FormValue("allowanceAmount"),
		AllowanceCheck:  r.FormValue("allowanceCheck"),
		AllowanceNo:     r.FormValue("allowanceNo"),
		InfoNo:          r.FormValue("infoNo"),
	}

	app.Setup.ModifyRentInfo(rent)

	r.Form.Set("infoNo", rent.InfoNo)
	app.FindRentByInfoNo(w, r)
}
