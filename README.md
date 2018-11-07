# go_tencent_sms
腾讯云短信

##使用方法
```go
//创建对象参数分别为：appid, appKey 和要使用的签名，签名必须事先审核通过。
Sms := tencent.NewSms(smsAppid, smsAppKey, smsSign)
//发送短信时参数为：短信模板ID、模板参数数组、接收的手机号码
err := Sms.Send(222124, []string{"352146", "10"}, "手机号码")
if err != nil {
  c.JSON(http.StatusOK, gin.H{
    "code": e.ERROR,
    "msg":  err.Error(),
  })
} else {
  c.JSON(http.StatusOK, gin.H{
    "code": 0,
  })
}
```
