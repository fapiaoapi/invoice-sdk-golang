package invoice

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// BlueTicket 数电蓝票开具接口
func (c *Client) BlueTicket(params map[string]string, items []InvoiceItem) (*InvoiceResponse, error) {
	// 构建请求参数
	requestParams := make(map[string]string)
	for k, v := range params {
		requestParams[k] = v
	}

	// 添加发票明细项
	for i, item := range items {
		prefix := fmt.Sprintf("fyxm[%d]", i)
		
		// 使用反射获取结构体字段和值
		itemData, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("序列化发票项失败: %v", err)
		}
		
		var itemMap map[string]interface{}
		if err := json.Unmarshal(itemData, &itemMap); err != nil {
			return nil, fmt.Errorf("解析发票项失败: %v", err)
		}
		
		// 添加到请求参数
		for k, v := range itemMap {
			if v != nil && v != "" {
				var strValue string
				switch val := v.(type) {
				case string:
					strValue = val
				case float64:
					strValue = strconv.FormatFloat(val, 'f', -1, 64)
				default:
					strValue = fmt.Sprintf("%v", val)
				}
				requestParams[prefix+"["+k+"]"] = strValue
			}
		}
	}

	resp, err := c.doRequest("POST", "/v5/enterprise/blueTicket", requestParams)
	if err != nil {
		return nil, err
	}

	var invoiceResp InvoiceResponse
	if err := parseResponseData(resp, &invoiceResp); err != nil {
		return nil, err
	}

	return &invoiceResp, nil
}

// GetVersionFile 获取销项数电版式文件
func (c *Client) GetVersionFile(nsrsbh, fphm, downflag string, options ...map[string]string) (*Response, error) {
	params := map[string]string{
		"nsrsbh":   nsrsbh,
		"fphm":     fphm,
		"downflag": downflag,
	}

	// 添加可选参数
	if len(options) > 0 {
		for k, v := range options[0] {
			params[k] = v
		}
	}

	return c.doRequest("POST", "/v5/enterprise/pdfOfdXml", params)
}

// QueryBlueInvoice 查蓝票信息接口
func (c *Client) QueryBlueInvoice(params map[string]string) (*Response, error) {
	return c.doRequest("POST", "/v5/enterprise/retMsg", params)
}

// ApplyRedInfo 申请红字信息表
func (c *Client) ApplyRedInfo(params map[string]string) (*Response, error) {
	return c.doRequest("POST", "/v5/enterprise/applyRedInfo", params)
}

// RedTicket 数电票开负数发票
func (c *Client) RedTicket(params map[string]string, items []InvoiceItem) (*Response, error) {
	// 构建请求参数
	requestParams := make(map[string]string)
	for k, v := range params {
		requestParams[k] = v
	}

	// 添加发票明细项
	for i, item := range items {
		prefix := fmt.Sprintf("fyxm[%d]", i)
		
		// 使用反射获取结构体字段和值
		itemData, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("序列化发票项失败: %v", err)
		}
		
		var itemMap map[string]interface{}
		if err := json.Unmarshal(itemData, &itemMap); err != nil {
			return nil, fmt.Errorf("解析发票项失败: %v", err)
		}
		
		// 添加到请求参数
		for k, v := range itemMap {
			if v != nil && v != "" {
				var strValue string
				switch val := v.(type) {
				case string:
					strValue = val
				case float64:
					strValue = strconv.FormatFloat(val, 'f', -1, 64)
				default:
					strValue = fmt.Sprintf("%v", val)
				}
				requestParams[prefix+"["+k+"]"] = strValue
			}
		}
	}

	return c.doRequest("POST", "/v5/enterprise/redTicket", requestParams)
}

// SwitchAccount 切换电子税务局账号
func (c *Client) SwitchAccount(nsrsbh string, options ...map[string]string) (*Response, error) {
	params := map[string]string{
		"nsrsbh": nsrsbh,
	}

	// 添加可选参数
	if len(options) > 0 {
		for k, v := range options[0] {
			params[k] = v
		}
	}

	return c.doRequest("POST", "/v5/enterprise/switchAccount", params)
}

// QueryCreditLimit 授信额度查询
func (c *Client) QueryCreditLimit(nsrsbh string, options ...map[string]string) (*Response, error) {
	params := map[string]string{
		"nsrsbh": nsrsbh,
	}

	// 添加可选参数
	if len(options) > 0 {
		for k, v := range options[0] {
			params[k] = v
		}
	}

	return c.doRequest("POST", "/v5/enterprise/queryCreditLimit", params)
}