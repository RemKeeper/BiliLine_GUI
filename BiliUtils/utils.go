package BiliUtils

import (
	"encoding/json"
	"fmt"
	gjson "github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type BiliCookieConfig struct {
	//AccessKey    string `json:"accessKey"`
	Csrf         string `json:"csrf"`
	Cookie       string `json:"cookie"`
	RefreshToken string `json:"refresh_token"`
}

const (
	UserAgent  = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36`
	CookiePath = "biliCookie.json"
)

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
func GetQRCodeState(loginKey string) (IsLogin bool, LoginData gjson.Result, err error) {
	apiUrl := "https://passport.bilibili.com/x/passport-login/web/qrcode/poll"
	client := http.Client{}
	req, _ := http.NewRequest("GET", apiUrl+fmt.Sprintf("?qrcode_key=%s", loginKey), nil)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return false, gjson.Result{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	data := gjson.ParseBytes(body)
	switch data.Get("data.code").Int() {
	case 0:
		cookieUrl := data.Get("data.url").String()
		parsedUrl, err := url.Parse(cookieUrl)
		if err != nil {
			return false, gjson.Result{}, err
		}
		cookieContentList := strings.Split(parsedUrl.RawQuery, "&")
		cookieContent := ""
		for _, cookie := range cookieContentList[:len(cookieContentList)-1] {
			cookieContent = cookieContent + cookie + ";"
		}
		cookieContent = strings.TrimSuffix(cookieContent, ";")
		login, result, csrf := VerifyLogin(cookieContent)

		if login {
			configInfo := BiliCookieConfig{
				Csrf:         csrf,
				Cookie:       cookieContent,
				RefreshToken: data.Get("data.refresh_token").String(),
			}
			indent, err := json.MarshalIndent(configInfo, "", "    ")
			if err != nil {
				return false, gjson.Result{}, err
			}
			os.WriteFile(CookiePath, indent, 0644)
			return true, result, nil
		} else {
			return false, gjson.Result{}, fmt.Errorf("登录失败")
		}

	case 86038:
		fmt.Println("二维码已失效，正在重新生成")
		return false, gjson.Result{}, fmt.Errorf("二维码失效")
	case 86090:
		fmt.Println("已扫码，请确认")
	case 86101:
	default:
		return false, gjson.Result{}, fmt.Errorf("未知code: %d", data.Get("data.code").Int())
	}
	return false, gjson.Result{}, nil
}

// 验证 cookie 可用性
func VerifyLogin(cookie string) (bool, gjson.Result, string) {
	u := "https://api.bilibili.com/x/web-interface/nav"
	client := http.Client{}
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	data := gjson.ParseBytes(body)

	isLogin := data.Get("data.isLogin").Bool()
	var csrf string
	if isLogin {
		reg := regexp.MustCompile(`bili_jct=([0-9a-zA-Z]+)`)
		csrf = reg.FindStringSubmatch(cookie)[1]
	}
	return isLogin, data, csrf
}
