# 发票SDK invoice-sdk-golang

[![Go Reference](https://pkg.go.dev/badge/github.com/fapiaoapi/invoice-sdk-golang.svg)](https://pkg.go.dev/github.com/fapiaoapi/invoice-sdk-golang)

电子发票/数电发票/全电发票 Golang SDK - 支持发票开具、红冲、查询等功能
基础

[中文文档](https://fa-piao.com/doc.html "文档")

* 获取授权
* 登录数电发票平台
* 获取人脸二维码
* 获取人脸二维码认证状态
* 获取认证状态

发票开具

* 数电蓝票开具接口
* 获取销项数电版式文件

发票红冲

* 申请红字前查蓝票信息接口
* 申请红字信息表
* 开负数发票

## 安装

```bash
go get github.com/fapiaoapi/invoice-sdk-golang
```

### 开票
 ```bash
// 配置信息
	appKey := "your_app_key"
	appSecret := "you_app_secret"

	nsrsbh := "91500112MADFAQ9xxx" // 统一社会信用代码
	// title := "重庆悦江河科技有限公司"         // 名称（营业执照）
	username := "19122840406" // 手机号码（电子税务局）
	password := ""            // 个人用户密码（电子税务局）
	// sf := "01"                     // 身份（电子税务局）
	fphm := "24502000000045823936"
	kprq := ""
	token := ""

	// 创建客户端
	client := invoice.NewClient(appKey, appSecret)

	// 获取授权
	if token != "" {
		client.SetToken(token)
	} else {
		authResult, err := client.GetAuthorization(nsrsbh)
		if err != nil {
			fmt.Printf("授权失败: %v\n", err)
			return
		}
		token = authResult.Token
		client.SetToken(token)
		fmt.Printf("授权成功，Token: %s\n", token)
	}

	// 查询认证状态
	loginResult, err := client.QueryFaceAuthState(nsrsbh, map[string]string{
		"username": username,
	})
	if err != nil {
		fmt.Printf("查询认证状态失败: %v\n", err)
		return
	}

	switch loginResult.Code {
	case 200:
		fmt.Println("可以开发票了")

		// 税额计算
		amount := 200.0
		taxRate := 0.01
		isIncludeTax := true // 是否含税
		se := invoice.CalculateTax(amount, taxRate, isIncludeTax)

		fmt.Printf("价税合计：%.2f\n", amount)
		fmt.Printf("税率：%.2f\n", taxRate)
		seFloat, _ := strconv.ParseFloat(se, 64)
		fmt.Printf("合计金额：%.2f\n", amount-seFloat)
		if isIncludeTax {
			fmt.Printf("含税 合计税额：%s\n", se)
		} else {
			fmt.Printf("不含税 合计税额：%s\n", se)
		}

		// // 授信额度查询
		// creditLimitResponse, err := client.QueryCreditLimit(nsrsbh, map[string]string{
		// 	"username": username,
		// })
		// if err != nil {
		// 	fmt.Printf("授信额度查询失败: %v\n", err)
		// } else if creditLimitResponse.Code == 200 {
		// 	fmt.Printf("授信额度查询结果: %s\n", string(creditLimitResponse.Data))
		// }

		// 开具蓝票
		invoiceParams := map[string]string{
			"fpqqlsh":  appKey + strconv.FormatInt(time.Now().Unix(), 10),
			"fplxdm":   "82",
			"kplx":     "0",
			"xhdwsbh":  nsrsbh,
			"xhdwmc":   "重庆悦江河科技有限公司",
			"xhdwdzdh": "重庆市渝北区仙桃街道汇业街1号17-2 19122840xxxx",
			"xhdwyhzh": "中国工商银行 310008670920023xxxx",
			"ghdwmc":   "个人",
			"zsfs":     "0",
			"hjje":     "9.9",
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
		} else {
			fmt.Printf("%d 开票结果: %s\n", loginResult.Code, loginResult.Msg)
			fphm = invoiceResponse.Fphm
			kprq = invoiceResponse.Kprq
			fmt.Printf("发票号码: %s\n", fphm)
			fmt.Printf("开票日期: %s\n", kprq)
		}

		// 下载发票
		pdfResponse, err := client.GetVersionFile(nsrsbh, fphm, "4", map[string]string{
			"kprq":     kprq,
			"username": username,
		})
		if err != nil {
			fmt.Printf("下载发票失败: %v\n", err)
		} else if pdfResponse.Code == 200 {
			fmt.Printf("下载发票结果: %s\n", string(pdfResponse.Data))
		}

	case 420:
		fmt.Println("登录(短信认证)")

		// 1. 发短信验证码
		loginResponse, err := client.LoginDppt(nsrsbh, username, password, "")
		if err != nil {
			fmt.Printf("发送短信验证码失败: %v\n", err)
			return
		}
		if loginResponse.Code == 200 {
			fmt.Println(loginResponse.Msg)
			fmt.Printf("请%s接收验证码\n", username)
			time.Sleep(60 * time.Second) // 等待60秒
		}

		// 2. 输入验证码
		fmt.Println("请输入验证码")
		var smsCode string
		fmt.Scanln(&smsCode) // 获取用户输入的验证码

		loginResponse2, err := client.LoginDppt(nsrsbh, username, password, smsCode)
		if err != nil {
			fmt.Printf("验证码登录失败: %v\n", err)
			return
		}
		if loginResponse2.Code == 200 {
			fmt.Println(string(loginResponse2.Data))
			fmt.Println("验证成功")
		}

	case 430:
		fmt.Println("人脸认证")

		// 1. 获取人脸二维码
		qrCodeResponse, err := client.GetFaceImg(nsrsbh, map[string]string{
			"username": username,
			"type":     "1",
		})
		if err != nil {
			fmt.Printf("获取人脸二维码失败: %v\n", err)
			return
		}
		fmt.Printf("二维码: %+v\n", qrCodeResponse)

		switch qrCodeResponse.Ewmly {
		case "swj":
			fmt.Println("请使用税务局app扫码")
		case "grsds":
			fmt.Println("个人所得税app扫码")
		}

		// 2. 认证完成后获取人脸二维码认证状态
		rzid := qrCodeResponse.Rzid
		faceStatusResponse, err := client.GetFaceState(nsrsbh, rzid, map[string]string{
			"username": username,
			"type":     "1",
		})
		if err != nil {
			fmt.Printf("获取人脸认证状态失败: %v\n", err)
			return
		}
		fmt.Printf("code: %d\n", loginResult.Code)
		fmt.Printf("data: %+v\n", faceStatusResponse)

		if faceStatusResponse != nil {
			switch faceStatusResponse.Slzt {
			case "1":
				fmt.Println("未认证")
			case "2":
				fmt.Println("成功")
			case "3":
				fmt.Println("二维码过期-->重新获取人脸二维码")
			}
		}

	case 401:
		fmt.Printf("%d 授权失败: %s\n", loginResult.Code, loginResult.Msg)

	default:
		fmt.Printf("%d %s\n", loginResult.Code, loginResult.Msg)
	}
```

### 发票红冲

```bash
// 配置信息
	appKey := "your_app_key"
	appSecret := "your_app_secret"
	nsrsbh := "915101820724315989" // 纳税人识别号
	username := "19122840xxx"      // 手机号码（电子税务局）
	fphm := "25502000000038381718"
	kprq := "2025-04-13 13:35:27"
	token := ""

	// 创建客户端
	client := invoice.NewClient(appKey, appSecret)

	if token != "" {
		client.SetToken(token)
	} else {
		// 获取授权
		authResponse, err := client.GetAuthorization(nsrsbh)
		if err != nil {
			fmt.Printf("授权失败: %v\n", err)
			return
		}
		fmt.Printf("授权成功，Token: %s\n", authResponse.Token)
	}

	// 1. 数电申请红字前查蓝票信息接口
	sqyy := "2"
	queryInvoiceParams := map[string]string{
		"nsrsbh":   nsrsbh,
		"fphm":     fphm,
		"sqyy":     sqyy,
		"username": username,
		"ghdwsbh":  nsrsbh,
	}
	queryInvoiceResponse, err := client.QueryBlueInvoice(queryInvoiceParams)
	if err != nil {
		fmt.Printf("查询发票信息失败: %v\n", err)
		return
	}

	if queryInvoiceResponse.Code == 200 {
		fmt.Println("1 可以申请红字")
		time.Sleep(2 * time.Second)

		// 2. 申请红字信息表
		applyRedParams := map[string]string{
			"xhdwsbh":  nsrsbh,
			"yfphm":    fphm,
			"username": username,
			"sqyy":     "2",
			"chyydm":   "01",
		}
		applyRedResponse, err := client.ApplyRedInfo(applyRedParams)
		if err != nil {
			fmt.Printf("申请红字信息表失败: %v\n", err)
			return
		}

		if applyRedResponse.Code == 200 {
			fmt.Println("2 申请红字信息表")
			time.Sleep(2 * time.Second)

			// 从响应中提取红字信息表编号
			var redInfoData struct {
				Xxbbh string `json:"xxbbh"`
			}
			if err := invoice.ParseResponseData(applyRedResponse, &redInfoData); err != nil {
				fmt.Printf("解析红字信息表编号失败: %v\n", err)
				return
			}

			// 3. 开具红字发票
			redInvoiceParams := map[string]string{
				"fpqqlsh":  "red" + fphm,
				"username": username,
				"xhdwsbh":  nsrsbh,
				"tzdbh":    redInfoData.Xxbbh,
				"yfphm":    fphm,
			}
			redInvoiceResponse, err := client.RedTicket(redInvoiceParams, nil)
			if err != nil {
				fmt.Printf("数电票负数开具失败: %v\n", err)
				return
			}

			if redInvoiceResponse.Code == 200 {
				fmt.Println("3 负数开具成功")
			} else {
				fmt.Printf("%d 数电票负数开具失败: %s\n", redInvoiceResponse.Code, redInvoiceResponse.Msg)
				fmt.Printf("错误详情: %s\n", string(redInvoiceResponse.Data))
			}
		} else {
			fmt.Printf("%d 申请红字信息表失败: %s\n", applyRedResponse.Code, applyRedResponse.Msg)
			fmt.Printf("错误详情: %s\n", string(applyRedResponse.Data))
		}
	} else {
		fmt.Printf("%d 查询发票信息失败: %s\n", queryInvoiceResponse.Code, queryInvoiceResponse.Msg)
		fmt.Printf("错误详情: %s\n", string(queryInvoiceResponse.Data))
	}
```
