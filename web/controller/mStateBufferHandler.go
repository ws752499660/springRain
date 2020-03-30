package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hqu.edu.cn/springRain/service"
)

func (app *Application) BufferUpdate(w http.ResponseWriter, r *http.Request) {
	con, _ := ioutil.ReadAll(r.Body) //获取post的数据
	mState := service.MState{}
	json.Unmarshal(con, &mState)
	resp := struct {
		Flag bool
		Msg  string
	}{
		Flag: false,
		Msg:  "",
	}

	msg, err := app.Setup.SaveMState(mState)
	if err != nil {
		msg, err = app.Setup.ModifyMState(mState)
		if err != nil {
			resp.Flag = false
			fmt.Println(err.Error())
		} else {
			resp.Flag = true
			fmt.Println("状态信息更新成功, 交易编号为: " + msg)
		}
	} else {
		resp.Flag = true
		fmt.Println("状态信息更新成功, 交易编号为: " + msg)
	}
	resp.Msg = msg

	jsonStr, err := json.Marshal(resp)

	w.Write(jsonStr)

}
