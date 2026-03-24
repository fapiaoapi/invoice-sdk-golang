package invoice

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

// BlueTicket 数电蓝票开具接口
func (c *Client) BlueTicket(params map[string]string, items []InvoiceItem) (*InvoiceResponse, error) {
	// 构建请求参数
	requestFields := buildMultipartFields(params, items)

	resp, err := c.doRequestWithFields("POST", "/v5/enterprise/blueTicket", requestFields)
	if err != nil {
		return nil, err
	}

	var invoiceResp InvoiceResponse
	invoiceResp.Code = resp.Code
	invoiceResp.Msg = resp.Msg
	invoiceResp.Total = resp.Total

	if resp.Data != nil && len(resp.Data) > 0 && string(resp.Data) != "null" {
		if err := json.Unmarshal(resp.Data, &invoiceResp); err != nil {
			return nil, fmt.Errorf("解析数据失败: %v", err)
		}
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
	requestFields, err := buildRedInfoFields(params)
	if err != nil {
		return nil, err
	}
	return c.doRequestWithFields("POST", "/v5/enterprise/hzxxbsq", requestFields)
}

// RedTicket 数电票开负数发票
func (c *Client) RedTicket(params map[string]string, items []InvoiceItem) (*Response, error) {
	// 构建请求参数
	requestFields := buildMultipartFields(params, items)
	return c.doRequestWithFields("POST", "/v5/enterprise/hzfpkj", requestFields)
}

// SyncRedInfo 同步红字信息表
func (c *Client) SyncRedInfo(params map[string]string) (*Response, error) {
	return c.doRequest("POST", "/v5/enterprise/hzxxbtb", params)
}

func buildMultipartFields(params map[string]string, items []InvoiceItem) []formField {
	fields := make([]formField, 0, len(params)+len(items)*8)
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		value := params[k]
		if value != "" {
			fields = append(fields, formField{Key: k, Value: value})
		}
	}
	for i, item := range items {
		prefix := fmt.Sprintf("fyxm[%d]", i)
		fields = append(fields, buildInvoiceItemFields(prefix, item)...)
	}
	return fields
}

func buildRedInfoFields(params map[string]string) ([]formField, error) {
	fields := make([]formField, 0, len(params)+8)
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if k != "fyxm" {
			value := params[k]
			if value == "" {
				continue
			}
			fields = append(fields, formField{Key: k, Value: value})
			continue
		}
		raw := params[k]
		if raw == "" {
			continue
		}
		var items []map[string]string
		if err := json.Unmarshal([]byte(raw), &items); err != nil {
			var itemsAny []map[string]interface{}
			if errAny := json.Unmarshal([]byte(raw), &itemsAny); errAny != nil {
				if raw != "" {
					fields = append(fields, formField{Key: k, Value: raw})
				}
				continue
			}
			items = make([]map[string]string, 0, len(itemsAny))
			for _, item := range itemsAny {
				converted := make(map[string]string, len(item))
				for key, val := range item {
					converted[key] = fmt.Sprintf("%v", val)
				}
				items = append(items, converted)
			}
		}
		for i, item := range items {
			itemKeys := make([]string, 0, len(item))
			for itemKey := range item {
				itemKeys = append(itemKeys, itemKey)
			}
			sort.Strings(itemKeys)
			for _, itemKey := range itemKeys {
				value := item[itemKey]
				if value == "" {
					continue
				}
				fields = append(fields, formField{
					Key:   fmt.Sprintf("fyxm[%d][%s]", i, itemKey),
					Value: value,
				})
			}
		}
	}
	return fields, nil
}

func buildInvoiceItemFields(prefix string, item InvoiceItem) []formField {
	value := reflect.ValueOf(item)
	itemType := reflect.TypeOf(item)
	fields := make([]formField, 0, itemType.NumField())
	for i := 0; i < itemType.NumField(); i++ {
		field := itemType.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		fieldValue := value.Field(i)
		var strValue string
		switch fieldValue.Kind() {
		case reflect.String:
			strValue = fieldValue.String()
		default:
			strValue = fmt.Sprintf("%v", fieldValue.Interface())
		}
		if strValue == "" {
			continue
		}
		fields = append(fields, formField{Key: prefix + "[" + tag + "]", Value: strValue})
	}
	return fields
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

	return c.doRequest("POST", "/v5/enterprise/changeUser", params)
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

	return c.doRequest("POST", "/v5/enterprise/creditLine", params)
}

func (c *Client) httpPost(url string, params map[string]string) (*Response, error) {
	return c.doRequest("POST", url, params)
}

func (c *Client) httpGet(url string, params map[string]string) (*Response, error) {
	return c.doRequest("GET", url, params)
}
