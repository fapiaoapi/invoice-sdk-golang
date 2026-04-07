package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"time"

	"github.com/fapiaoapi/invoice-sdk-golang"
	"github.com/redis/go-redis/v9"
	"github.com/skip2/go-qrcode"
)

var (
	appKey    = ""
	appSecret = ""
	nsrsbh    = "" // 统一社会信用代码
	username  = "" // 手机号码（电子税务局）
	password  = "" // 个人用户密码（电子税务局）

	ewmlx = "10" // 1 电子税务局app人脸二维码登录，10 电子税务局app扫码登录2 个税人脸二维码登录，3 个税 app 扫码确认登录

	token       = ""
	accountType = "7"  //默认6 6基础 7标准
	debug       = true // true 是否开启调试模式 false 不开启
)

type loginDpptQRData struct {
	Qrcode string `json:"qrcode"`
	Ewmid  string `json:"ewmid"`
}

func main66() {

	// 从缓存redis中获取Token
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务地址（默认端口6379）
		Password: "test123456",     // Redis密码（无密码则留空）
		DB:       0,                // 使用默认数据库
	})
	// 创建客户端
	client := invoice.NewClient(appKey, appSecret, debug)
	// 获取授权
	GetToken(rdb, client, false)

	faceStateResponse, err := client.QueryFaceAuthState(nsrsbh, map[string]string{
		"username": username,
	})
	if err != nil {
		fmt.Printf("获取状态失败: %v\n", err)
		return
	}
	switch faceStateResponse.Code {
	case 200:
		fmt.Println("不需要认证")
		break
	case 420:
		fmt.Println("420 扫码认证")
		loginResponse, err3 := client.LoginDppt(nsrsbh, username, password, "", map[string]string{
			"ewmlx": ewmlx,
		})
		if err3 != nil {
			fmt.Printf("发送扫码登录失败: %v\n", err3)
			return
		}
		if loginResponse.Code == 200 {
			fmt.Println("成功做完人脸认证,请输入数字 1")
			var qrData loginDpptQRData
			if err := json.Unmarshal(loginResponse.Data, &qrData); err != nil {
				fmt.Printf("解析扫码数据失败: %v\n", err)
				return
			}
			qrcode := qrData.Qrcode
			ewmid := qrData.Ewmid
			if qrcode == "" || ewmid == "" {
				fmt.Printf("扫码数据缺失: %s\n", string(loginResponse.Data))
				return
			}
			length := len(qrcode)
			if length < 500 {
				StringToQrcode(qrcode)
			} else {
				Base64StringToQrcode(qrcode)
			}
			fmt.Printf("请在300秒内(%s前)输入: ", time.Now().Add(300*time.Second).Format("2006-01-02 15:04:05"))
			input := make(chan string, 1)
			go func() {
				sms, _ := bufio.NewReader(os.Stdin).ReadString('\n')
				input <- sms
			}()
			select {
			case sms := <-input:
				fmt.Printf("\n你输入了: %s", sms)
				loginResponse2, err4 := client.LoginDppt(nsrsbh, username, password, "", map[string]string{
					"ewmlx": ewmlx,
					"ewmid": ewmid,
				})
				if err4 != nil {
					fmt.Printf("扫码验证失败: %v\n", err4)
					return
				}
				//税务app登陆信息与当前纳税人不符，请您查看税务app信息或使用短信验证码登陆
				//在电子税务局app首页右上角，点击“身份切换”按钮，选择对应的企业，点击“切换”
				if loginResponse2.Code == 200 {
					fmt.Println(string(loginResponse2.Data))
					fmt.Println("扫码验证成功")
					faceStateResponse2, err := client.QueryFaceAuthState(nsrsbh, map[string]string{
						"username": username,
					})
					if err != nil {
						fmt.Printf("获取状态失败: %v\n", err)
						return
					}
					fmt.Printf("QueryFaceAuthState: %d %s\n", faceStateResponse2.Code, faceStateResponse2.Msg)
				} else {
					fmt.Printf("扫码验证失败: %s\n", loginResponse2.Msg)
					return
				}
			case <-time.After(300 * time.Second):
				fmt.Println("\n\n[错误] 输入超时！")
			}
		}

	case 430:
		fmt.Println("人脸认证")
		break
	case 401:
		fmt.Printf("%d 授权失败: %s\n", faceStateResponse.Code, faceStateResponse.Msg)
		// 重新授权获取token
		break
	default:
		fmt.Printf("异常 %d %s\n", faceStateResponse.Code, faceStateResponse.Msg)
		break
	}

}

// 获取token 从redis获取或重新获取
func GetToken(rdb *redis.Client, client *invoice.Client, forceUpdate bool) {
	key := nsrsbh + "@" + username + "@TOKEN"
	if forceUpdate {
		/**
		 * 获取授权Token文档
		 * @link https://fa-piao.com/doc.html#api1?source=github
		 */
		authResult, authErr := client.GetAuthorization(nsrsbh, accountType, "", "")
		// authResult, authErr := client.GetAuthorization(nsrsbh, accountType, username, password)
		if authErr != nil {
			fmt.Printf("授权失败: %v\n", authErr)
		} else {
			token = authResult.Token
			client.SetToken(token)
			//缓存Token到redis 过期时间30天 key建议是nsrsbh+'@TOKEN'
			err := rdb.Set(context.Background(), key, token, 30*24*time.Hour).Err()
			if err != nil {
				fmt.Printf("缓存Token到Redis失败: %v\n", err)
			}
			fmt.Printf("授权成功，Token: %s\n", token)
		}
	} else {
		result, err := rdb.Get(context.Background(), key).Result()
		if err == nil && result != "" {
			token = result
			client.SetToken(token)
			fmt.Printf("Token from Redis")
		} else if err != nil && err != redis.Nil {
			fmt.Printf("从Redis获取Token异常: %v\n", err)
		} else {
			authResult, authErr := client.GetAuthorization(nsrsbh, accountType, "", "")
			// authResult, authErr := client.GetAuthorization(nsrsbh, accountType, username, password)
			if authErr != nil {
				fmt.Printf("授权失败: %v\n", authErr)
			} else {
				token = authResult.Token
				client.SetToken(token)
				//缓存Token到redis 过期时间30天 key建议是nsrsbh+'@TOKEN'
				err := rdb.Set(context.Background(), key, token, 30*24*time.Hour).Err()
				if err != nil {
					fmt.Printf("缓存Token到Redis失败: %v\n", err)
				}
				fmt.Printf("授权成功，Token: %s\n", token)
			}
		}
	}
}

func StringToQrcode(text string) {
	qr, err := qrcode.New(text, qrcode.Low)
	if err != nil {
		panic(err)
	}

	matrix := qr.Bitmap()

	for y := 0; y < len(matrix); y += 2 {
		for x := 0; x < len(matrix[y]); x++ {
			top := matrix[y][x]
			bottom := false
			if y+1 < len(matrix) {
				bottom = matrix[y+1][x]
			}

			switch {
			case top && bottom:
				fmt.Print("█")
			case top && !bottom:
				fmt.Print("▀")
			case !top && bottom:
				fmt.Print("▄")
			default:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

}

func Base64StringToQrcode(base64Str string) {
	// 1. 解码 Base64 数据
	imageData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		fmt.Printf("解码 Base64 失败: %v\n", err)
		return
	}

	// 2. 从二进制数据中解码 PNG 图片
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		fmt.Printf("解码图片失败: %v\n", err)
		return
	}

	// 3. 从图片中提取二维码矩阵（黑白）
	matrix := extractQRMatrix(img)
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		fmt.Println("二维码矩阵提取失败")
		return
	}

	// 4. 用与 stringToQrcode 相同的字符策略输出，保证尺寸和清晰度一致
	for y := 0; y < len(matrix); y += 2 {
		for x := 0; x < len(matrix[y]); x++ {
			top := matrix[y][x]
			bottom := false
			if y+1 < len(matrix) {
				bottom = matrix[y+1][x]
			}

			switch {
			case top && bottom:
				fmt.Print("█")
			case top && !bottom:
				fmt.Print("▀")
			case !top && bottom:
				fmt.Print("▄")
			default:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func extractQRMatrix(img image.Image) [][]bool {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width == 0 || height == 0 {
		return nil
	}

	binary := make([][]bool, height)
	for y := 0; y < height; y++ {
		binary[y] = make([]bool, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			luma := (299*(r>>8) + 587*(g>>8) + 114*(b>>8)) / 1000
			binary[y][x] = luma < 128
		}
	}

	minX, minY := width, height
	maxX, maxY := -1, -1
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if !binary[y][x] {
				continue
			}
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}
		}
	}
	if maxX < minX || maxY < minY {
		return nil
	}

	cropW := maxX - minX + 1
	cropH := maxY - minY + 1
	moduleCount := 21
	bestScore := 1e9
	for c := 21; c <= 177; c += 4 {
		cellW := float64(cropW) / float64(c)
		cellH := float64(cropH) / float64(c)
		roundW := float64(int(cellW + 0.5))
		roundH := float64(int(cellH + 0.5))
		if roundW < 1 || roundH < 1 {
			continue
		}
		score := absFloat(cellW-roundW) + absFloat(cellH-roundH) + absFloat(cellW-cellH)
		if score < bestScore {
			bestScore = score
			moduleCount = c
		}
	}

	matrix := make([][]bool, moduleCount)
	for y := 0; y < moduleCount; y++ {
		matrix[y] = make([]bool, moduleCount)
		for x := 0; x < moduleCount; x++ {
			sampleX := minX + (x*cropW+cropW/(2*moduleCount))/moduleCount
			sampleY := minY + (y*cropH+cropH/(2*moduleCount))/moduleCount
			if sampleX > maxX {
				sampleX = maxX
			}
			if sampleY > maxY {
				sampleY = maxY
			}
			matrix[y][x] = binary[sampleY][sampleX]
		}
	}
	return matrix
}

func absFloat(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
