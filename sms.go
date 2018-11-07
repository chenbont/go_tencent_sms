package tencent

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"models"
	"net/http"
	"strconv"
	"time"
)

type sms struct {
}

var smsAppid, smsAppKey, smsSign string

func NewSms(appid, appkey, sign string) *sms {
	var Sms sms
	smsAppid = appid
	smsAppKey = appkey
	smsSign = sign
	return &Sms
}

func (*sms) Send(tpl_id int, params []string, mobile string) error {
	type reqStruct struct {
		Ext    string   `json:"ext"`
		Extend string   `json:"extend"`
		Params []string `json:"params"`
		Sig    string   `json:"sig"`
		Sign   string   `json:"sign"`
		Tel    struct {
			Mobile     string `json:"mobile"`
			Nationcode string `json:"nationcode"`
		} `json:"tel"`
		Time  int64 `json:"time"`
		TplID int   `json:"tpl_id"`
	}
	var m reqStruct
	now := time.Now().Unix()
	random := RandomStr(12)
	m.Ext = ""
	m.Extend = ""
	m.Params = params
	m.Sig = smsCalcSign("appkey=" + smsAppKey + "&random=" + random + "&time=" + strconv.FormatInt(now, 10) + "&mobile=" + mobile)
	m.Sign = smsSign
	m.Tel.Mobile = mobile
	m.Tel.Nationcode = "86"
	m.Time = now
	m.TplID = tpl_id

	jsonStu, err := json.Marshal(m)
	if err != nil {
		return errors.New("生成json字符串错误")
	}
	var jsonStr = []byte(jsonStu)

	req, err := http.NewRequest("POST", "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid="+smsAppid+"&random="+random, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("New Http Request发生错误，原因:", err)
		return errors.New("Http Request发生错误")

	}
	req.Header.Set("Accept", "application/json")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	client := http.Client{}
	resp, _err := client.Do(req)
	if _err != nil {
		fmt.Println("短信发送失败, 原因:", _err)
		return errors.New("短信发送失败")
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	models.PrintLog(string(respBytes))
	if err != nil {
		fmt.Println("返回body错误", err)
		return errors.New("返回body错误")
	}
	type respStruct struct {
		Result int    `json:"result"`
		Errmsg string `json:"errmsg"`
		Ext    string `json:"ext"`
		Fee    int    `json:"fee"`
		Sid    string `json:"sid"`
	}
	var rr respStruct

	xml.Unmarshal(respBytes, &rr)
	models.PrintLog(rr)
	//处理return code.
	if rr.Result != 0 {
		fmt.Println("短信发送失败，原因:", rr.Errmsg, " str_req-->", rr.Ext)
		return errors.New("短信发送失败，原因:" + rr.Errmsg)
	} else {
		return nil
	}
}

func smsCalcSign(m string) string {
	models.PrintLog(m)
	h := sha256.New()
	h.Write([]byte(m))

	return fmt.Sprintf("%x", h.Sum(nil))
}

//RandomStr 随机生成字符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
