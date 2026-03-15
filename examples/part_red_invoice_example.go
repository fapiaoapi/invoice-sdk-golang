package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fapiaoapi/invoice-sdk-golang"
)

func main() {
	// 配置信息
	appKey := ""
	appSecret := ""
	nsrsbh := ""   // 纳税人识别号
	username := "" // 手机号码（电子税务局）
	fphm := ""
	// kprq := "2025-04-13 13:35:27"
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

	/**
	 * 1. 数电申请红字前查蓝票信息接口
	 * @link https://fa-piao.com/doc.html#api8?source=github
	 */
	sqyy := "2" //2 销方申请
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

		/**
		 * 2. 申请红字信息表
		 * @link https://fa-piao.com/doc.html#api9?source=github
		 */
		fyxmJSON, err := json.Marshal([]map[string]string{
			{
				"hsbz": "0",
				"je":   "-70.80",
				"se":   "-9.20",
				"spsl": "-1.00",
				"xh":   "1",
			},
		})
		if err != nil {
			fmt.Printf("序列化 fyxm 失败: %v\n", err)
			return
		}
		applyRedParams := map[string]string{
			"bfch":     "1",
			"chyydm":   "02",
			"fyxm":     string(fyxmJSON),
			"hjje":     "-70.80",
			"hjse":     "-9.20",
			"sqyy":     "2",
			"username": username,
			"xhdwsbh":  nsrsbh,
			"yfphm":    fphm,
			"jyrzzt":   "1",
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

			/**
			 * 3. 开具红字发票
			 * @link https://fa-piao.com/doc.html#api10?source=github
			 */
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
}
