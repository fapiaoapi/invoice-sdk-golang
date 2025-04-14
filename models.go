package invoice

import (
	"encoding/json"
	"fmt"
)

// Response 通用响应结构
type Response struct {
	Code  int             `json:"code"`
	Msg   string          `json:"msg"`
	Data  json.RawMessage `json:"data"`
	Total int             `json:"total"`
}

// 判断响应是否成功
func (r *Response) IsSuccess() bool {
	return r.Code == 200
}

// 获取错误信息
func (r *Response) Error() string {
	if r.IsSuccess() {
		return ""
	}
	return fmt.Sprintf("错误码: %d, 错误信息: %s", r.Code, r.Msg)
}

// AuthResponse 授权响应
type AuthResponse struct {
	Token string `json:"token"`
}

// InvoiceResponse 发票开具响应
type InvoiceResponse struct {
	Fphm           string `json:"Fphm"`           // 发票号码
	Kprq           string `json:"Kprq"`           // 开票日期
	Gmfyx          string `json:"Gmfyx"`          // 购买方邮箱
	GmfSsjswjgdm   string `json:"GmfSsjswjgdm"`   // 购买方税局机关代码
	Ewm            string `json:"ewm"`            // 发票打印的二维码
	Zzfpdm         string `json:"zzfpdm"`         // 纸质发票代码
	Zzfphm         string `json:"zzfphm"`         // 纸质发票号码
}

// FaceQRCodeResponse 人脸二维码响应
type FaceQRCodeResponse struct {
	Rzid   string `json:"rzid"`   // 认证id
	Nsrsbh string `json:"nsrsbh"` // 纳税人识别号
	Ewm    string `json:"ewm"`    // 二维码
	Slzt   string `json:"slzt"`   // 受理状态
	Ewmly  string `json:"ewmly"`  // 二维码来源
}

// FaceStateResponse 人脸认证状态响应
type FaceStateResponse struct {
	Rzid   string `json:"rzid"`   // 认证id
	Nsrsbh string `json:"nsrsbh"` // 纳税人识别号
	Ewm    string `json:"ewm"`    // 二维码
	Slzt   string `json:"slzt"`   // 受理状态
}

// InvoiceItem 发票明细项
type InvoiceItem struct {
	Fphxz   string  `json:"fphxz"`   // 发票行性质
	Spmc    string  `json:"spmc"`    // 商品名称
	Ggxh    string  `json:"ggxh"`    // 规格型号
	Dw      string  `json:"dw"`      // 单位
	Spsl    string  `json:"spsl"`    // 商品数量
	Dj      string  `json:"dj"`      // 单价
	Je      string  `json:"je"`      // 金额
	Sl      string  `json:"sl"`      // 税率
	Se      string  `json:"se"`      // 税额
	Hsbz    string  `json:"hsbz"`    // 含税标志
	Spbm    string  `json:"spbm"`    // 商品编码
	Yhzcbs  string  `json:"yhzcbs"`  // 优惠赠策标识
	Lslbs   string  `json:"lslbs"`   // 零税率标识
	Zzstsgl string  `json:"zzstsgl"` // 增值税特殊管理
	MtzlDm  string  `json:"mtzlDm"`  // 煤炭种类代码
}

// VersionFileResponse 版式文件响应
type VersionFileResponse struct {
	OfdUrl string `json:"ofdUrl"` // OFD文件URL
	PdfUrl string `json:"pdfUrl"` // PDF文件URL
	XmlUrl string `json:"xmlUrl"` // XML文件URL
}