package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hqu.edu.cn/springRain/service"
)

// 显示更新信息页面
func (app *Application) UpdateMStateShow(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "updateMState.html", data)
}

// 更新信息
func (app *Application) UpdateMState(w http.ResponseWriter, r *http.Request) {

	mState := service.MState{
		MachineID:     r.FormValue("machineID"),
		CollectDate:   r.FormValue("collectDate"),
		UserID:        r.FormValue("userID"),
		OwnInc:        r.FormValue("ownInc"),
		Location:      r.FormValue("location"),
		EngCode:       r.FormValue("engCode"),
		OilAmountLeft: r.FormValue("oilAmountLeft"),
		TirePressure:  r.FormValue("tirePressure"),
		CoolantTemp:   r.FormValue("coolantTemp"),
		OtherInfo:     r.FormValue("otherInfo"),
	}

	msg, err := app.Setup.SaveMState(mState)
	if err != nil {
		msg, err = app.Setup.ModifyMState(mState)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		fmt.Println("装填信息更新成功, 交易编号为: " + msg)
	}

	r.Form.Set("userID", mState.UserID)
	r.Form.Set("ownInc", mState.OwnInc)
	app.FindMStateByUserIDAndOwnInc(w, r)
}

func (app *Application) QueryMStateByUserIDAndOwnIncShow(w http.ResponseWriter, r *http.Request) {
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
	ShowView(w, r, "queryMStateByUserIDAndOwnInc.html", data)
}

// 根据使用者ID和租赁公司查询信息
func (app *Application) FindMStateByUserIDAndOwnInc(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userID")
	ownInc := r.FormValue("ownInc")
	result, err := app.Setup.FindMStateByUserIDAndOwnInc(userID, ownInc)
	var mState = service.MState{}
	json.Unmarshal(result, &mState)

	fmt.Println("根据使用者ID和租赁公司查询信息成功：")
	fmt.Println(mState)

	data := &struct {
		MState      service.MState
		CurrentUser User
		Msg         string
		Flag        bool
		History     bool
	}{
		MState:      mState,
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		History:     false,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "queryMStateResult.html", data)
}

func (app *Application) QueryMStateByMachineID(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "queryMStateByMachineID.html", data)
}

// 根据机器ID查询状态信息
func (app *Application) FindMStateByMachineID(w http.ResponseWriter, r *http.Request) {
	machineID := r.FormValue("machineID")
	result, err := app.Setup.FindMStateByMachineID(machineID)
	var mState = service.MState{}
	json.Unmarshal(result, &mState)

	data := &struct {
		MState      service.MState
		CurrentUser User
		Msg         string
		Flag        bool
		History     bool
	}{
		MState:      mState,
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		History:     true,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "queryMStateResult.html", data)
}
