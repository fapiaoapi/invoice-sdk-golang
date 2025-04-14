package invoice

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

// 生成随机字符串
func generateRandomString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	rand.Read(b) // 安全随机数生成
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

// 计算HMAC-SHA256签名
func calculateSignature(method, path, randomString, timestamp, appKey, appSecret string) string {
	signContent := fmt.Sprintf(
		"Method=%s&Path=%s&RandomString=%s&TimeStamp=%s&AppKey=%s",
		method, path, randomString, timestamp, appKey,
	)

	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write([]byte(signContent))
	return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}

// 创建multipart/form-data请求体
func createRequestBody(params map[string]string) (io.Reader, string) {
	body := &strings.Builder{}
	writer := multipart.NewWriter(body)

	for key, value := range params {
		_ = writer.WriteField(key, value)
	}
	writer.Close()

	return strings.NewReader(body.String()), writer.FormDataContentType()
}

// 处理HTTP响应
func handleResponse(resp *http.Response) (*Response, error) {
	if resp == nil {
		return nil, fmt.Errorf("空响应")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}

// FormatAmount 格式化金额（保留2位小数）
func FormatAmount(amount interface{}) string {
	var floatAmount float64
	
	switch v := amount.(type) {
	case float64:
		floatAmount = v
	case float32:
		floatAmount = float64(v)
	case int:
		floatAmount = float64(v)
	case int64:
		floatAmount = float64(v)
	case string:
		floatAmount, _ = strconv.ParseFloat(v, 64)
	default:
		return "0.00"
	}
	
	return strconv.FormatFloat(floatAmount, 'f', 2, 64)
}

// CalculateTax 计算税额
func CalculateTax(amount, taxRate interface{}, isIncludeTax bool) string {
	amountFloat := parseFloat(amount)
	taxRateFloat := parseFloat(taxRate)
	
	var tax float64
	if isIncludeTax {
		// 含税计算：税额 = 金额 / (1 + 税率) * 税率
		tax = amountFloat / (1 + taxRateFloat) * taxRateFloat
	} else {
		// 不含税计算：税额 = 金额 * 税率
		tax = amountFloat * taxRateFloat
	}
	
	return FormatAmount(tax)
}

// CalculateAmountWithoutTax 计算不含税金额
func CalculateAmountWithoutTax(amount, taxRate interface{}) string {
	amountFloat := parseFloat(amount)
	taxRateFloat := parseFloat(taxRate)
	
	// 不含税金额 = 含税金额 / (1 + 税率)
	amountWithoutTax := amountFloat / (1 + taxRateFloat)
	
	return FormatAmount(amountWithoutTax)
}

// CalculateAmountWithTax 计算含税金额
func CalculateAmountWithTax(amount, taxRate interface{}) string {
	amountFloat := parseFloat(amount)
	taxRateFloat := parseFloat(taxRate)
	
	// 含税金额 = 不含税金额 * (1 + 税率)
	amountWithTax := amountFloat * (1 + taxRateFloat)
	
	return FormatAmount(amountWithTax)
}

// parseFloat 将各种类型转换为float64
func parseFloat(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		result, _ := strconv.ParseFloat(v, 64)
		return result
	default:
		return 0
	}
}
