package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hqu.edu.cn/springRain/service"
	_ "github.com/mattn/go-sqlite3"
)

func mbufferTest(serviceSetup service.ServiceSetup) {
	fmt.Println("====================农机终端农机状态监测缓冲上链测试===================")

	// 打开数据库
	db, err := sql.Open("sqlite3", "database/m.db")
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//测试数据

	state0 := service.MState{
		MachineID: "203",

		CollectDate: "2020-3-18 10:55",
		UserID:      "U103",
		OwnInc:      "CMCC",

		Location:      "(24.6097159200,118.0910292300)",
		EngCode:       "None",
		OilAmountLeft: "59",
		TirePressure:  "[243,244,246,250]",
		CoolantTemp:   "92",
		OtherInfo:     "{'RoutineMaintenanceCheck':'True','WorkingType':'Normal'}第一个缓冲",
	}
	state1 := service.MState{
		MachineID: "203",

		CollectDate: "2020-3-18 12:55",
		UserID:      "U103",
		OwnInc:      "CMCC",

		Location:      "(25.6097159200,118.0910292300)",
		EngCode:       "None",
		OilAmountLeft: "50",
		TirePressure:  "[243,244,246,250]",
		CoolantTemp:   "90",
		OtherInfo:     "{'RoutineMaintenanceCheck':'True','WorkingType':'Normal'}第二个缓冲",
	}
	state2 := service.MState{
		MachineID: "203",

		CollectDate: "2020-3-18 13:05",
		UserID:      "U103",
		OwnInc:      "CMCC",

		Location:      "(25.6097159200,118.0910292312)",
		EngCode:       "None",
		OilAmountLeft: "49",
		TirePressure:  "[243,244,246,250]",
		CoolantTemp:   "96",
		OtherInfo:     "{'RoutineMaintenanceCheck':'True','WorkingType':'Normal'}第三个缓冲",
	}

	// 向数据库插入数据
	err = insert(state0, db)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = insert(state1, db)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = insert(state2, db)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//向链中提交数据
	err = commitAll(serviceSetup, db)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 查询测试结果
	result, err := serviceSetup.FindMStateByMachineID("203")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var mState service.MState
		json.Unmarshal(result, &mState)
		fmt.Println("根据机器编号查询信息成功：")
		fmt.Println(mState)
	}

}

// 增加数据
func insert(mstate service.MState, db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO main.mstate('MachineID', 'CollectDate', 'UserID', 'OwnInc', 'Location', 'EngCode', 'OilAmountLeft', 'TirePressure', 'CoolantTemp', 'OtherInfo') VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	res, err := stmt.Exec(mstate.MachineID, mstate.CollectDate, mstate.UserID, mstate.OwnInc, mstate.Location, mstate.EngCode, mstate.OilAmountLeft, mstate.TirePressure, mstate.CoolantTemp, mstate.OtherInfo)
	_, err = res.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		return nil
	}
}

// 不提供更新数据的接口，因为每次记录都要上链，所以如果有更新需要则是进行增加新的状态到数据库中

// 删除数据（不对外开放，仅用于将数据成功上链后将数据库中对应行删除使用）
func delete(collectDate string, db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM main.mstate WHERE CollectDate=?")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	res, err := stmt.Exec(collectDate)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = res.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		return nil
	}
}

// 查询当前数据库中全部数据
func queryAll(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT * FROM main.mstate")
	return rows, err
}

// 将数据库中所有数据提交至链上
func commitAll(serviceSetup service.ServiceSetup, db *sql.DB) error {
	rows, err := queryAll(db)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	var states []service.MState
	for rows.Next() {
		var state service.MState
		err := rows.Scan(&state.MachineID, &state.CollectDate, &state.UserID, &state.OwnInc, &state.Location, &state.EngCode, &state.OilAmountLeft, &state.TirePressure, &state.CoolantTemp, &state.OtherInfo)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		states = append(states, state)
	}
	rows.Close() //db.Query()默认不会释放连接池，只有用rows.Close()才会释放连接池

	// 如果成功将一条记录提交到链上，则从数据库中删除该记录
	for i := 0; i < len(states); i++ {
		err = commitToChain(states[i], serviceSetup)
		if err != nil {
			fmt.Println(err.Error())
			return err
		} else {
			err = delete(states[i].CollectDate, db)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
		}
	}
	return nil

}

func commitToChain(mState service.MState, serviceSetup service.ServiceSetup) error {
	msg, err := serviceSetup.SaveMState(mState)
	if err != nil {
		msg, err = serviceSetup.ModifyMState(mState)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	} else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
		return nil
	}
	return nil
}
