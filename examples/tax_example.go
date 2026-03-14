package main

import (
	"encoding/json"
	"fmt"

	"github.com/fapiaoapi/invoice-sdk-golang"
	"github.com/shopspring/decimal"
)

/**
 * 含税金额计算示例
 *
 *   不含税单价 = 含税单价/(1+ 税率)  noTaxDj = dj / (1 + sl)
 *   不含税金额 = 不含税单价*数量  noTaxJe = noTaxDj * spsl
 *   含税金额 = 含税单价*数量  je = dj * spsl
 *   税额 = 税额 = 1 / (1 + 税率) * 税率 * 含税金额  se = 1  / (1 + sl) * sl * je
 *    hjse= se1 + se2 + ... + seN
 *    jshj= je1 + je2 + ... + jeN
 *   价税合计 =合计金额+合计税额 jshj = hjje + hjse
 *
 */
func main44() {
	/**
	 * 含税计算示例1  无价格  无数量
	 * @link `https://fa-piao.com/fapiao.html?action=data1&source=github`
	 *
	 */

	hsbz := 1
	amount := 200.0
	sl := 0.01
	se := invoice.CalculateTax(amount, sl, hsbz == 1, 2)
	data := map[string]interface{}{
		"hjje": 0.0,
		"hjse": 0.0,
		"jshj": 0.0,
		"fyxm": []map[string]interface{}{
			{
				"fphxz": 0,
				"hsbz":  hsbz,
				"spmc":  "*软件维护服务*接口服务费",
				"spbm":  "3040201030000000000",
				"je":    amount,
				"sl":    sl,
				"se":    se,
			},
		},
	}

	for _, item := range data["fyxm"].([]map[string]interface{}) {
		data["jshj"] = add(data["jshj"].(float64), item["je"].(float64), 2)
		data["hjse"] = add(data["hjse"].(float64), item["se"].(float64), 2)
	}
	data["hjje"] = sub(data["jshj"].(float64), data["hjse"].(float64), 2)

	fmt.Println("含税计算示例1  无价格  无数量: ")
	printJSON(data)
	fmt.Println("---------------------------------------------")

	/**
	 * 含税计算示例2  有价格 有数量
	 * @link `https://fa-piao.com/fapiao.html?action=data3&source=github`
	 *
	 */

	hsbz = 1
	spsl := 1.0
	dj := 2.0
	sl = 0.03

	spsl2 := 1.0
	dj2 := 3.0
	sl2 := 0.01

	je := mul(dj, spsl, 2)
	se = invoice.CalculateTax(je, sl, hsbz == 1, 2)

	je2 := mul(dj2, spsl2, 2)
	se2 := invoice.CalculateTax(je2, sl2, hsbz == 1, 2)
	data = map[string]interface{}{
		"hjje": 0.0,
		"hjse": 0.0,
		"jshj": 0.0,
		"fyxm": []map[string]interface{}{
			{
				"fphxz": 0,
				"hsbz":  hsbz,
				"spmc":  "*水冰雪*一阶水费",
				"spbm":  "1100301030000000000",
				"ggxh":  "",
				"dw":    "吨",
				"dj":    dj,
				"spsl":  spsl,
				"je":    je,
				"sl":    sl,
				"se":    se,
			},
			{
				"fphxz": 0,
				"hsbz":  hsbz,
				"spmc":  "*水冰雪*二阶水费",
				"spbm":  "1100301030000000000",
				"ggxh":  "",
				"dw":    "吨",
				"dj":    dj2,
				"spsl":  spsl2,
				"je":    je2,
				"sl":    sl2,
				"se":    se2,
			},
		},
	}

	for _, item := range data["fyxm"].([]map[string]interface{}) {
		data["jshj"] = add(data["jshj"].(float64), item["je"].(float64), 2)
		data["hjse"] = add(data["hjse"].(float64), item["se"].(float64), 2)
	}
	data["hjje"] = sub(data["jshj"].(float64), data["hjse"].(float64), 2)

	fmt.Println("含税计算示例2  有价格 有数量: ")
	printJSON(data)
	fmt.Println("---------------------------------------------")

	/**
	 * 含税计算示例3  有价格自动算数量  购买猪肉1000元,16.8元/斤
	 * @link `https://fa-piao.com/fapiao.html?action=data5&source=github`
	 *
	 */
	hsbz = 1
	amount = 1000.0
	dj = 16.8
	sl = 0.01
	se = invoice.CalculateTax(amount, sl, hsbz == 1, 2)

	data = map[string]interface{}{
		"hjje": 0.0,
		"hjse": 0.0,
		"jshj": 0.0,
		"fyxm": []map[string]interface{}{
			{
				"fphxz": 0,
				"hsbz":  hsbz,
				"spmc":  "*肉类*猪肉",
				"spbm":  "1030107010100000000",
				"ggxh":  "",
				"dw":    "斤",
				"dj":    dj,
				"spsl":  div(amount, dj, 13),
				"je":    amount,
				"sl":    sl,
				"se":    se,
			},
		},
	}
	for _, item := range data["fyxm"].([]map[string]interface{}) {
		data["jshj"] = add(data["jshj"].(float64), item["je"].(float64), 2)
		data["hjse"] = add(data["hjse"].(float64), item["se"].(float64), 2)
	}
	data["hjje"] = sub(data["jshj"].(float64), data["hjse"].(float64), 2)

	fmt.Println("含税计算示例3  有价格自动算数量 购买猪肉1000元,16.8元/斤: ")
	printJSON(data)
	fmt.Println("---------------------------------------------")

	/**
	 * 含税计算示例4  有数量自动算价格  购买接口服务1000元7次
	 * @link `https://fa-piao.com/fapiao.html?action=data7&source=github`
	 *
	 */

	hsbz = 1
	amount = 1000.0
	spsl = 7.0
	sl = 0.01
	se = invoice.CalculateTax(amount, sl, hsbz == 1, 2)

	data = map[string]interface{}{
		"hjje": 0.0,
		"hjse": 0.0,
		"jshj": 0.0,
		"fyxm": []map[string]interface{}{
			{
				"fphxz": 0,
				"hsbz":  hsbz,
				"spmc":  "*软件维护服务*接口服务费",
				"spbm":  "3040201030000000000",
				"ggxh":  "",
				"dw":    "次",
				"dj":    div(amount, spsl, 13),
				"spsl":  spsl,
				"je":    amount,
				"sl":    sl,
				"se":    se,
			},
		},
	}
	for _, item := range data["fyxm"].([]map[string]interface{}) {
		data["jshj"] = add(data["jshj"].(float64), item["je"].(float64), 2)
		data["hjse"] = add(data["hjse"].(float64), item["se"].(float64), 2)
	}
	data["hjje"] = sub(data["jshj"].(float64), data["hjse"].(float64), 2)

	fmt.Println("含税计算示例4  有数量自动算价格 购买接口服务1000元7次: ")
	printJSON(data)
	fmt.Println("---------------------------------------------")

	/**
	 * 不含税计算示例
	 *  金额 = 单价 * 数量  je = dj * spsl
	 *  税额 = 金额 * 税率  se = je * sl
	 *   hjse= se1 + se2 + ... + seN
	 *   hjje= je1 + je2 + ... + jeN
	 *  价税合计 =合计金额+合计税额 jshj = hjje + hjse
	 *
	 */

	/**
	 *
	 * 不含税计算示例1 无价格 无数量
	 * @link `https://fa-piao.com/fapiao.html?action=data2&source=github`
	 */

	hsbz = 0
	amount = 200.0
	sl = 0.01
	se = invoice.CalculateTax(amount, sl, hsbz == 1, 2)
	data = map[string]interface{}{
		"hjje": 0.0,
		"hjse": 0.0,
		"jshj": 0.0,
		"fyxm": []map[string]interface{}{
			{
				"fphxz": 0,
				"hsbz":  hsbz,
				"spmc":  "*软件维护服务*接口服务费",
				"spbm":  "3040201030000000000",
				"je":    amount,
				"sl":    sl,
				"se":    se,
			},
		},
	}

	for _, item := range data["fyxm"].([]map[string]interface{}) {
		data["hjje"] = add(data["hjje"].(float64), item["je"].(float64), 2)
		data["hjse"] = add(data["hjse"].(float64), item["se"].(float64), 2)
	}
	data["jshj"] = add(data["hjje"].(float64), data["hjse"].(float64), 2)

	fmt.Println("不含税计算示例1 无价格 无数量: ")
	printJSON(data)
	fmt.Println("---------------------------------------------")

	/**
	 *
	 * 不含税计算示例2  有价格 有数量
	 * @link `https://fa-piao.com/fapiao.html?action=data4&source=github`
	 */
	// 一阶水费 1吨，单价2元/吨，税率0.03
	// 二阶水费 1吨，单价3元/吨，税率0.01
	hsbz = 0
	spsl = 1.0
	dj = 2.0
	sl = 0.03

	spsl2 = 1.0
	dj2 = 3.0
	sl2 = 0.01

	je = mul(dj, spsl, 2)
	se = invoice.CalculateTax(je, sl, hsbz == 1, 2)

	je2 = mul(dj2, spsl2, 2)
	se2 = invoice.CalculateTax(je2, sl2, hsbz == 1, 2)
	data = map[string]interface{}{
		"hjje": 0.0,
		"hjse": 0.0,
		"jshj": 0.0,
		"fyxm": []map[string]interface{}{
			{
				"fphxz": 0,
				"hsbz":  hsbz,
				"spmc":  "*水冰雪*一阶水费",
				"spbm":  "1100301030000000000",
				"ggxh":  "",
				"dw":    "吨",
				"dj":    dj,
				"spsl":  spsl,
				"je":    je,
				"sl":    sl,
				"se":    se,
			},
			{
				"fphxz": 0,
				"hsbz":  hsbz,
				"spmc":  "*水冰雪*而阶水费",
				"spbm":  "1100301030000000000",
				"ggxh":  "",
				"dw":    "吨",
				"dj":    dj2,
				"spsl":  spsl2,
				"je":    je2,
				"sl":    sl2,
				"se":    se2,
			},
		},
	}

	for _, item := range data["fyxm"].([]map[string]interface{}) {
		data["hjje"] = add(data["hjje"].(float64), item["je"].(float64), 2)
		data["hjse"] = add(data["hjse"].(float64), item["se"].(float64), 2)
	}
	data["jshj"] = add(data["hjje"].(float64), data["hjse"].(float64), 2)

	fmt.Println("不含税计算示例2  有价格 有数量: ")
	printJSON(data)
	fmt.Println("---------------------------------------------")

	/**
	 * 不含税计算示例3  有价格自动算数量  购买猪肉1000元,16.8元/斤
	 * @link `https://fa-piao.com/fapiao.html?action=data6&source=github`
	 *
	 */
	hsbz = 0
	amount = 1000.0
	dj = 16.8
	sl = 0.01
	se = invoice.CalculateTax(amount, sl, hsbz == 1, 2)

	data = map[string]interface{}{
		"hjje": 0.0,
		"hjse": 0.0,
		"jshj": 0.0,
		"fyxm": []map[string]interface{}{
			{
				"fphxz": 0,
				"hsbz":  hsbz,
				"spmc":  "*肉类*猪肉",
				"spbm":  "1030107010100000000",
				"ggxh":  "",
				"dw":    "斤",
				"dj":    dj,
				"spsl":  div(amount, dj, 13),
				"je":    amount,
				"sl":    sl,
				"se":    se,
			},
		},
	}
	for _, item := range data["fyxm"].([]map[string]interface{}) {
		data["hjje"] = add(data["hjje"].(float64), item["je"].(float64), 2)
		data["hjse"] = add(data["hjse"].(float64), item["se"].(float64), 2)
	}
	data["jshj"] = add(data["hjje"].(float64), data["hjse"].(float64), 2)
	fmt.Println("不含税计算示例3  有价格自动算数量 购买猪肉1000元,16.8元/斤: ")
	printJSON(data)
	fmt.Println("---------------------------------------------")

	/**
	 * 不含税计算示例4  有数量自动算价格  购买接口服务1000元7次
	 *
	 * @link `https://fa-piao.com/fapiao.html?action=data8&source=github`
	 *
	 */

	hsbz = 0
	amount = 1000.0
	spsl = 7.0
	sl = 0.01
	se = invoice.CalculateTax(amount, sl, hsbz == 1, 2)

	data = map[string]interface{}{
		"hjje": 0.0,
		"hjse": 0.0,
		"jshj": 0.0,
		"fyxm": []map[string]interface{}{
			{
				"fphxz": 0,
				"hsbz":  hsbz,
				"spmc":  "*软件维护服务*接口服务费",
				"spbm":  "1030107010100000000",
				"ggxh":  "",
				"dw":    "次",
				"dj":    div(amount, spsl, 13),
				"spsl":  spsl,
				"je":    amount,
				"sl":    sl,
				"se":    se,
			},
		},
	}
	for _, item := range data["fyxm"].([]map[string]interface{}) {
		data["hjje"] = add(data["hjje"].(float64), item["je"].(float64), 2)
		data["hjse"] = add(data["hjse"].(float64), item["se"].(float64), 2)
	}
	data["jshj"] = add(data["hjje"].(float64), data["hjse"].(float64), 2)
	fmt.Println("不含税计算示例4  有数量自动算价格 购买接口服务1000元7次: ")
	printJSON(data)
	fmt.Println("---------------------------------------------")

	/**
	 * 免税计算示例
	 *  金额 = 单价 * 数量  je = dj * spsl
	 *  税额 = 0
	 *  hjse = se1 + se2 + ... + seN
	 *  jshj = je1 + je2 + ... + jeN
	 *  价税合计 =合计金额+合计税额 jshj = hjje + hjse
	 * @link `https://fa-piao.com/fapiao.html?action=data9&source=github`
	 */

	hsbz = 0
	dj = 32263.98
	sl = 0.0
	se = 0.0
	data = map[string]interface{}{
		"hjje": 0.0,
		"hjse": 0.0,
		"jshj": 0.0,
		"fyxm": []map[string]interface{}{
			{
				"fphxz":   0,
				"hsbz":    hsbz,
				"spmc":    "*经纪代理服务*国际货物运输代理服务",
				"spbm":    "3040802010200000000",
				"ggxh":    "",
				"dw":      "次",
				"spsl":    1,
				"dj":      dj,
				"je":      dj,
				"sl":      sl,
				"se":      se,
				"yhzcbs":  1,
				"lslbs":   1,
				"zzstsgl": "免税",
			},
		},
	}
	for _, item := range data["fyxm"].([]map[string]interface{}) {
		data["hjje"] = add(data["hjje"].(float64), item["je"].(float64), 2)
		data["hjse"] = add(data["hjse"].(float64), item["se"].(float64), 2)
	}

	data["jshj"] = add(data["hjje"].(float64), data["hjse"].(float64), 2)

	fmt.Print("免税计算示例: ")
	printJSON(data)
}

func add(a, b float64, scale int32) float64 {
	v, _ := decimal.NewFromFloat(a).Add(decimal.NewFromFloat(b)).Round(scale).Float64()
	return v
}

func sub(a, b float64, scale int32) float64 {
	v, _ := decimal.NewFromFloat(a).Sub(decimal.NewFromFloat(b)).Round(scale).Float64()
	return v
}

func mul(a, b float64, scale int32) float64 {
	v, _ := decimal.NewFromFloat(a).Mul(decimal.NewFromFloat(b)).Round(scale).Float64()
	return v
}

func div(a, b float64, scale int32) float64 {
	v, _ := decimal.NewFromFloat(a).Div(decimal.NewFromFloat(b)).Round(scale).Float64()
	return v
}

func printJSON(data map[string]interface{}) {
	encoder := json.NewEncoder(consoleWriter{})
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(data)
}

type consoleWriter struct{}

func (consoleWriter) Write(p []byte) (int, error) {
	return fmt.Print(string(p))
}
