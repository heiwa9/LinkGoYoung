package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"os"
)

type userInfo struct {
	UserHard     string `json:"user_hard"`
	UserAccount  string `json:"user_account"`
	PassWord     string `json:"pass_word"`
	LastLoginURL string `json:"last_login_url"`
}

func (info *userInfo) SaveUserInfoJson(label *widget.Label) {
	data, err := json.Marshal(info)
	if err != nil {
		fmt.Println("JSON 序列化失败" + err.Error())
		label.SetText("JSON 序列化失败")
		return
	}
	err = ioutil.WriteFile("data.json", data, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
		label.SetText("JSON 写入失败")
		return
	}
}

func (info *userInfo) ReadUserInfoJson(label *widget.Label) {
	// 读取JSON文件内容 返回字节切片
	bytes, err := ioutil.ReadFile("./data.json")
	if err != nil {
		fmt.Println("data.json 读取失败" + err.Error())
		label.SetText("data.json 读取失败")
		return
	}
	//fmt.Println("data.json content:")
	//// 打印时需要转为字符串
	//fmt.Println(string(bytes))
	//// 将字节切片映射到指定结构上
	err = json.Unmarshal(bytes, &info)
	if err != nil {
		fmt.Println("JSON 反序列化失败" + err.Error())
		label.SetText("JSON 反序列化失败")
		return
	}
}
