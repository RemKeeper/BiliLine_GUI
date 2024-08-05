package BiliUtils

import (
	"fmt"
	gjson "github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type BiliCookieConfig struct {
	//AccessKey    string `json:"accessKey"`
	Cookie       string `json:"cookie"`
	RefreshToken string `json:"refresh_token"`
}

const UserAgent = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36`

// GetLoginKeyAndLoginUrl 获取二维码内容和密钥
func GetLoginKeyAndLoginUrl() (QrKey, QrUrl string) {
	GetQrUrl := "https://passport.bilibili.com/x/passport-login/web/qrcode/generate"
	client := http.Client{}
	req, _ := http.NewRequest("GET", GetQrUrl, nil)
	req.Header.Set("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	data := gjson.ParseBytes(body)
	loginKey := data.Get("data.qrcode_key").String()
	loginUrl := data.Get("data.url").String()
	return loginKey, loginUrl
}

// GetQRCodeState 获取二维码状态并验证登录状态
func GetQRCodeState(loginKey string) (bool, error) {
	apiUrl := "https://passport.bilibili.com/x/passport-login/web/qrcode/poll"
	client := http.Client{}
	req, _ := http.NewRequest("GET", apiUrl+fmt.Sprintf("?qrcode_key=%s", loginKey), nil)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	data := gjson.ParseBytes(body)
	switch data.Get("data.code").Int() {
	case 0:
		cookieUrl := data.Get("data.url").String()
		parsedUrl, err := url.Parse(cookieUrl)
		if err != nil {
			return false, err
		}
		cookieContentList := strings.Split(parsedUrl.RawQuery, "&")
		cookieContent := ""
		for _, cookie := range cookieContentList[:len(cookieContentList)-1] {
			cookieContent = cookieContent + cookie + ";"
		}
		cookieContent = strings.TrimSuffix(cookieContent, ";")
		configInfo := BiliCookieConfig{
			Cookie:       cookieContent,
			RefreshToken: data.Get("data.refresh_token").String(),
		}

	case 86038:
		fmt.Println("二维码已失效，正在重新生成")
		return false, fmt.Errorf("二维码失效")
	case 86090:
		fmt.Println("已扫码，请确认")
	case 86101:
	default:
		return false, fmt.Errorf("未知code: %d", data.Get("data.code").Int())
	}
	return true, nil
}
