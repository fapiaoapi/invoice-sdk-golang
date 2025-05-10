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
	"strings"

	"github.com/shopspring/decimal"
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

// CalculateTax 计算税额
func CalculateTax(amount float64, taxRate float64, isIncludeTax bool, newScale int) (se float64) {
	// 转换为 decimal 类型
	amountDecimal := decimal.NewFromFloat(amount)
	taxRateDecimal := decimal.NewFromFloat(taxRate)

	var tax decimal.Decimal
	if isIncludeTax {
		// 含税计算：税额 = 金额 / (1 + 税率) * 税率
		one := decimal.NewFromInt(1)
		denominator := one.Add(taxRateDecimal)

		// 计算税额
		tax = amountDecimal.Div(denominator).Mul(taxRateDecimal)
	} else {
		// 不含税计算：税额 = 金额 * 税率
		tax = amountDecimal.Mul(taxRateDecimal)
	}

	// 设置精度并四舍五入
	tax = tax.Round(int32(newScale))
	// // 计算不含税金额
	// exje = amountDecimal.Sub(tax).Round(int32(newScale)).InexactFloat64()

	// 转换为 float64 类型
	se, _ = tax.Float64()

	return se
}
