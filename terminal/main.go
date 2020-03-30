package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/hqu.edu.cn/springRain/service"
	_ "github.com/mattn/go-sqlite3"
)

const url string = "http://127.0.0.1:9000/bufferUpdate"

func main() {
	// 打开数据库
	db, err := sql.Open("sqlite3", "../database/m.db")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	maxTimes := 10
	count := 0

	for {
		count++
		randNum := rand.Intn(2) + 1
		insert(db)
		if randNum%2 == 0 {
			// 网络畅通
			err := commitAll(db)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// 网络环境差
			fmt.Println("无法连接至服务器，本次提交数据暂存缓存中")
			time.Sleep(time.Duration(2) * time.Second)
		}
		if count >= maxTimes {
			break
		}
	}

	defer db.Close()

}

//向缓存中添加数据
func insert(db *sql.DB) error {
	collectDate := time.Now().Format("2006-01-02 15:04:05")

	location := [2]float64{24.6097159200, 118.0910292300}
	loStr := "("
	for i := 0; i < 2; i++ {
		location[i] = location[i] + float64(rand.Intn(100)/100)
		loStr = loStr + strconv.FormatFloat(location[i], 'E', -1, 64) + ","
	}
	loStr = loStr + ")"

	oil := rand.Intn(20) + 50
	oilStr := strconv.Itoa(oil)

	mstate := service.MState{
		MachineID: "205",

		CollectDate: collectDate,
		UserID:      "U109",
		OwnInc:      "HQU.CST.LAB",

		Location:      loStr,
		EngCode:       "None",
		OilAmountLeft: oilStr,
		TirePressure:  "[243,244,246,250]",
		CoolantTemp:   "92",
		OtherInfo:     "本条信息经缓存后发出",
	}

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
		fmt.Println("向缓存中写入数据：")
		fmt.Println(mstate)
		return nil
	}
}

// 将缓存中所有数据提交至链上
func commitAll(db *sql.DB) error {
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
		err, flag := commitToChain(states[i])
		if err != nil {
			fmt.Println(err.Error())
			return err
		} else {
			if flag {
				err = delete(states[i].CollectDate, db)
			}
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
		}
	}
	return nil
}

//提交到链上的具体过程
func commitToChain(mState service.MState) (error, bool) {

	jsonStr, err := json.Marshal(mState)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("网络请求失败")
		return err, false
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	reciv := struct {
		Flag bool
		Msg  string
	}{
		Flag: false,
		Msg:  "",
	}
	json.Unmarshal(body, &reciv)

	if reciv.Flag {
		fmt.Println("信息发布成功, 编号为: " + reciv.Msg)
	}

	return err, reciv.Flag
}

// 查询当前数据库中全部数据
func queryAll(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT * FROM main.mstate")
	return rows, err
}

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
