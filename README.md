# 电子发票/数电发票 Golang SDK | 开票、验真、红冲一站式集成

[![Go Reference](https://pkg.go.dev/badge/github.com/fapiaoapi/invoice-sdk-golang.svg)](https://pkg.go.dev/github.com/fapiaoapi/invoice-sdk-golang)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/fapiaoapi/invoice-sdk-golang/blob/master/LICENSE)

**发票 Golang SDK** 专为电子发票、数电发票（全电发票）场景设计，支持**开票、红冲、版式文件下载**等核心功能，快速对接税务平台API。

**关键词**: 电子发票SDK,数电票Golang,开票接口,发票api,发票开具,发票红冲,全电发票集成

---

## 📖 核心功能

### 基础认证
- ✅ **获取授权** - 快速接入税务平台身份认证
- ✅ **人脸二维码登录** - 支持数电发票平台扫码登录
- ✅ **认证状态查询** - 实时获取纳税人身份状态

### 发票开具
- 🟦 **数电蓝票开具** - 支持增值税普通/专用电子发票
- 📄 **版式文件下载** - 自动获取销项发票PDF/OFD/XML文件

### 发票红冲
- 🔍 **红冲前蓝票查询** - 精确检索待红冲的电子发票
- 🛑 **红字信息表申请** - 生成红冲凭证
- 🔄 **负数发票开具** - 自动化红冲流程

---

## 🚀 快速安装

## 安装

```bash
go get github.com/fapiaoapi/invoice-sdk-golang
```


[📚 查看完整中文文档](https://fa-piao.com/doc.html?source=github)

---

## 🔍 为什么选择此SDK？
- **精准覆盖中国数电发票标准** - 严格遵循国家最新接口规范
- **开箱即用** - 无需处理XML/签名等底层细节，专注业务逻辑
- **企业级验证** - 已在生产环境处理超100万张电子发票

---

## 📊 支持的开票类型
| 发票类型       | 状态   |
|----------------|--------|
| 数电发票（普通发票） | ✅ 支持 |
| 数电发票（增值税专用发票） | ✅ 支持 |
| 数电发票（铁路电子客票）  | ✅ 支持 |
| 数电发票（航空运输电子客票行程单） | ✅ 支持  |
| 数电票（二手车销售统一发票） | ✅ 支持  |
| 数电纸质发票（增值税专用发票） | ✅ 支持  |
| 数电纸质发票（普通发票） | ✅ 支持  |
| 数电纸质发票（机动车发票） | ✅ 支持  |
| 数电纸质发票（二手车发票） | ✅ 支持  |

---

## 🤝 贡献与支持
- 提交Issue: [问题反馈](https://github.com/fapiaoapi/invoice-sdk-golang/issues)
- 商务合作: yuejianghe@qq.com
### 开票
 ```bash
package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/fapiaoapi/invoice-sdk-golang"
	"github.com/redis/go-redis/v9"
	// "github.com/skip2/go-qrcode"
	// "github.com/shopspring/decimal"
)

func main() {
	// 配置信息
	appKey := ""
	appSecret := ""

	nsrsbh := "91500112MADFAQXXXX" // 统一社会信用代码
	title := "XXXX科技有限公司"          // 名称（营业执照）
	username := "1325580xxxx"      // 手机号码（电子税务局）
	// password := "1356325"          // 个人用户密码（电子税务局）
	// sf := "01"                     // 身份（电子税务局）
	fphm := "245020000000xxx"
	kprq := ""
	token := ""
	accountType := "6" //默认6 6基础 7标准

	// 创建客户端
	client := invoice.NewClient(appKey, appSecret)

	//一 获取授权

	// 从缓存redis中获取Token
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务地址（默认端口6379）
		Password: "test123456",     // Redis密码（无密码则留空）
		DB:       0,                // 使用默认数据库
	})
	key := nsrsbh + "@" + username + "@TOKEN"
	result, err := rdb.Get(context.Background(), key).Result()
	if err == nil && result != "" {
		token = result
		client.SetToken(token)
		fmt.Printf("Token from Redis")
	} else if err != nil && err != redis.Nil {
		fmt.Printf("从Redis获取Token异常: %v\n", err)
	}

	if token == "" {
		/**
		 * 获取授权Token文档
		 * @link https://fa-piao.com/doc.html#api1?source=github
		 */
		authResult, authErr := client.GetAuthorization(nsrsbh, accountType, "", "")
		// authResult, authErr := client.GetAuthorization(nsrsbh, accountType, username, password)
		if authErr != nil {
			fmt.Printf("授权失败: %v\n", authErr)
			return
		}
		token = authResult.Token
		client.SetToken(token)
		//缓存Token到redis 过期时间30天 key建议是nsrsbh+'@TOKEN'
		err = rdb.Set(context.Background(), key, token, 30*24*time.Hour).Err()
		if err != nil {
			fmt.Printf("缓存Token到Redis失败: %v\n", err)
			return
		}
		fmt.Printf("授权成功，Token: %s\n", token)
	}

	/**
	 * 前端模拟数电发票/电子发票开具 (蓝字发票)
	 * @link https://fa-piao.com/fapiao.html?source=github
	 *
	 */

	//二 开具蓝票

	/**
	 * 开具数电发票文档
	 * @link https://fa-piao.com/doc.html#api6?source=github
	 *
	 * 开票参数说明demo
	 * @link https://github.com/fapiaoapi/invoice-sdk-golang/blob/master/examples/tax_example.go
	 */
	invoiceParams := map[string]string{
		"username": username,
		"fpqqlsh":  appKey + strconv.FormatInt(time.Now().Unix(), 10),
		"fplxdm":   "82",
		"kplx":     "0",
		"xhdwsbh":  nsrsbh,
		"xhdwmc":   title,
		"xhdwdzdh": "重庆市两江区xxxx 1912284xxxx",
		"xhdwyhzh": "中国工商银行 31000867092xxxx",
		"ghdwmc":   "个人",
		"zsfs":     "0",
		"hjje":     "9.99",
		"hjse":     "0.1",
		"jshj":     "10",
	}

	items := []invoice.InvoiceItem{
		{
			Fphxz: "0",
			Spmc:  "*信息技术服务*软件开发服务",
			Je:    "10",
			Sl:    "0.01",
			Se:    "0.1",
			Hsbz:  "1",
			Spbm:  "3040201010000000000",
		},
	}

	invoiceResponse, err := client.BlueTicket(invoiceParams, items)
	if err != nil {
		fmt.Printf("开票失败: %v\n", err)
		return
	}
	if invoiceResponse == nil {
		fmt.Println("开票失败: 响应为空")
		return
	}
	switch invoiceResponse.Code {
	case 200:
		fmt.Printf("%d 开票结果: %s\n", invoiceResponse.Code, invoiceResponse.Msg)
		fphm = invoiceResponse.Fphm
		kprq = invoiceResponse.Kprq
		fmt.Printf("发票号码: %s\n", fphm)
		fmt.Printf("开票日期: %s\n", kprq)
		/**
		* 获取销项数电版式文件 文档 PDF/OFD/XML
		* @link https://fa-piao.com/doc.html#api7?source=github
		 */
		pdfResponse, err2 := client.GetVersionFile(nsrsbh, fphm, "4", map[string]string{
			"kprq":     kprq,
			"username": username,
		})
		if err2 != nil {
			fmt.Printf("下载发票失败: %v\n", err2)
		} else if pdfResponse.Code == 200 {
			fmt.Printf("下载发票结果: %s\n", string(pdfResponse.Data))
		}
	case 420:
		fmt.Println("登录(短信认证)")
		/**
		 * 前端模拟短信认证弹窗
		 * @link https://fa-piao.com/fapiao.html?action=sms&source=github
		 */

		// // 1. 发短信验证码
		// /**
		//  * @link https://fa-piao.com/doc.html#api2?source=github
		//  */
		// loginResponse, err3 := client.LoginDppt(nsrsbh, username, password, "")
		// if err3 != nil {
		// 	fmt.Printf("发送短信验证码失败: %v\n", err3)
		// 	return
		// }
		// if loginResponse.Code == 200 {
		// 	fmt.Println(loginResponse.Msg)
		// 	fmt.Printf("请%s接收验证码\n", username)
		// 	time.Sleep(60 * time.Second) // 等待60秒
		// }
		// // 2. 输入验证码
		// /**
		//  * @link https://fa-piao.com/doc.html#api2?source=github
		//  */
		// fmt.Println("请输入验证码")
		// var smsCode string
		// fmt.Scanln(&smsCode) // 获取用户输入的验证码

		// loginResponse2, err4 := client.LoginDppt(nsrsbh, username, password, smsCode)
		// if err4 != nil {
		// 	fmt.Printf("验证码登录失败: %v\n", err4)
		// 	return
		// }
		// if loginResponse2.Code == 200 {
		// 	fmt.Println(string(loginResponse2.Data))
		// 	fmt.Println("验证成功")
		// 	//todo 重新调用相关接口
		// }

	case 430:
		fmt.Println("人脸认证")
		/**
		 * 前端模拟人脸认证弹窗
		 * @link https://fa-piao.com/fapiao.html?action=face&source=github
		 */
		// 1. 获取人脸二维码
		/**
		 * @link https://fa-piao.com/doc.html#api3?source=github
		 */
		// qrCodeResponse, err5 := client.GetFaceImg(nsrsbh, map[string]string{
		// 	"username": username,
		// 	"type":     "1",
		// })
		// if err5 != nil {
		// 	fmt.Printf("获取人脸二维码失败: %v\n", err5)
		// 	return
		// }
		// //判断 Ewmly 不为空 长度小于500 字符串转二维码图片base64
		// if qrCodeResponse.Ewm != "" && len(qrCodeResponse.Ewm) <= 500 {
		// 	//go get github.com/skip2/go-qrcode
		// 	qrBase64, err6 := StringToQRCodeBase64(qrCodeResponse.Ewm)
		// 	if err6 != nil {
		// 		fmt.Printf("生成二维码base64失败: %v\n", err6)
		// 		return
		// 	} else {
		// 		qrCodeResponse.Ewm = qrBase64
		// 		// 可选：加上 data URL 前缀，方便在 HTML 中直接使用
		// 		qrCodeResponse.Ewm = "data:image/png;base64," + qrBase64
		// 		fmt.Printf("人脸二维码base64: %s\n", qrCodeResponse.Ewm)
		// 	}

		// }

		// switch qrCodeResponse.Ewmly {
		// case "swj":
		// 	fmt.Println("请使用税务局app扫码")
		// case "grsds":
		// 	fmt.Println("个人所得税app扫码")
		// }

		// // 2. 认证完成后获取人脸二维码认证状态
		// /**
		//  * @link https://fa-piao.com/doc.html#api4?source=github
		//  */
		// rzid := qrCodeResponse.Rzid
		// faceStatusResponse, err6 := client.GetFaceState(nsrsbh, rzid, map[string]string{
		// 	"username": username,
		// 	"type":     "1",
		// })
		// if err6 != nil {
		// 	fmt.Printf("获取人脸认证状态失败: %v\n", err6)
		// 	return
		// }
		// fmt.Printf("data: %+v\n", faceStatusResponse)

		// if faceStatusResponse != nil {
		// 	switch faceStatusResponse.Slzt {
		// 	case "1":
		// 		fmt.Println("未认证")
		// 	case "2":
		// 		fmt.Println("成功")
		// 		//todo 重新调用相关接口
		// 	case "3":
		// 		fmt.Println("二维码过期-->重新获取人脸二维码")
		// 	}
		// }
	case 401:
		fmt.Printf("%d 授权失败: %s\n", invoiceResponse.Code, invoiceResponse.Msg)
		// 重新授权获取token
	default:
		fmt.Printf("%d %s\n", invoiceResponse.Code, invoiceResponse.Msg)
		fmt.Printf("失败错误: %v\n", err)
	}

}

// func StringToQRCodeBase64(text string) (string, error) {
// 	// 生成二维码（PNG 格式写入内存 buffer）
// 	var buf bytes.Buffer
// 	qr, err := qrcode.New(text, qrcode.Medium)
// 	if err != nil {
// 		return "", err
// 	}
// 	if err := qr.Write(256, &buf); err != nil {
// 		return "", err
// 	}

// 	// 将字节转为 Base64
// 	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

// 	return encoded, nil
// }


```
[发票税额计算demo](examples/tax_example.go "发票税额计算") |
[发票红冲demo](examples/red_invoice_example.go "发票红冲")
[部分发票红冲demo](examples/part_red_invoice_example.go "部分发票红冲")
