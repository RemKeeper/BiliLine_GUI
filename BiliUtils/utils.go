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

type UserInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Card struct {
			Mid         string        `json:"mid"`
			Name        string        `json:"name"`
			Approve     bool          `json:"approve"`
			Sex         string        `json:"sex"`
			Rank        string        `json:"rank"`
			Face        string        `json:"face"`
			FaceNft     int           `json:"face_nft"`
			FaceNftType int           `json:"face_nft_type"`
			DisplayRank string        `json:"DisplayRank"`
			Regtime     int           `json:"regtime"`
			Spacesta    int           `json:"spacesta"`
			Birthday    string        `json:"birthday"`
			Place       string        `json:"place"`
			Description string        `json:"description"`
			Article     int           `json:"article"`
			Attentions  []interface{} `json:"attentions"`
			Fans        int           `json:"fans"`
			Friend      int           `json:"friend"`
			Attention   int           `json:"attention"`
			Sign        string        `json:"sign"`
			LevelInfo   struct {
				CurrentLevel int `json:"current_level"`
				CurrentMin   int `json:"current_min"`
				CurrentExp   int `json:"current_exp"`
				NextExp      int `json:"next_exp"`
			} `json:"level_info"`
			Pendant struct {
				Pid               int    `json:"pid"`
				Name              string `json:"name"`
				Image             string `json:"image"`
				Expire            int    `json:"expire"`
				ImageEnhance      string `json:"image_enhance"`
				ImageEnhanceFrame string `json:"image_enhance_frame"`
				NPid              int    `json:"n_pid"`
			} `json:"pendant"`
			Nameplate struct {
				Nid        int    `json:"nid"`
				Name       string `json:"name"`
				Image      string `json:"image"`
				ImageSmall string `json:"image_small"`
				Level      string `json:"level"`
				Condition  string `json:"condition"`
			} `json:"nameplate"`
			Official struct {
				Role  int    `json:"role"`
				Title string `json:"title"`
				Desc  string `json:"desc"`
				Type  int    `json:"type"`
			} `json:"Official"`
			OfficialVerify struct {
				Type int    `json:"type"`
				Desc string `json:"desc"`
			} `json:"official_verify"`
			Vip struct {
				Type       int   `json:"type"`
				Status     int   `json:"status"`
				DueDate    int64 `json:"due_date"`
				VipPayType int   `json:"vip_pay_type"`
				ThemeType  int   `json:"theme_type"`
				Label      struct {
					Path                  string `json:"path"`
					Text                  string `json:"text"`
					LabelTheme            string `json:"label_theme"`
					TextColor             string `json:"text_color"`
					BgStyle               int    `json:"bg_style"`
					BgColor               string `json:"bg_color"`
					BorderColor           string `json:"border_color"`
					UseImgLabel           bool   `json:"use_img_label"`
					ImgLabelUriHans       string `json:"img_label_uri_hans"`
					ImgLabelUriHant       string `json:"img_label_uri_hant"`
					ImgLabelUriHansStatic string `json:"img_label_uri_hans_static"`
					ImgLabelUriHantStatic string `json:"img_label_uri_hant_static"`
				} `json:"label"`
				AvatarSubscript    int    `json:"avatar_subscript"`
				NicknameColor      string `json:"nickname_color"`
				Role               int    `json:"role"`
				AvatarSubscriptUrl string `json:"avatar_subscript_url"`
				TvVipStatus        int    `json:"tv_vip_status"`
				TvVipPayType       int    `json:"tv_vip_pay_type"`
				TvDueDate          int    `json:"tv_due_date"`
				AvatarIcon         struct {
					IconType     int `json:"icon_type"`
					IconResource struct {
					} `json:"icon_resource"`
				} `json:"avatar_icon"`
				VipType   int `json:"vipType"`
				VipStatus int `json:"vipStatus"`
			} `json:"vip"`
			IsSeniorMember int         `json:"is_senior_member"`
			NameRender     interface{} `json:"name_render"`
		} `json:"card"`
		Following    bool `json:"following"`
		ArchiveCount int  `json:"archive_count"`
		ArticleCount int  `json:"article_count"`
		Follower     int  `json:"follower"`
		LikeNum      int  `json:"like_num"`
	} `json:"data"`
}

const (
	UserAgent      = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36`
	CookiePath     = "biliCookie.json"
	GetQrUrl       = "https://passport.bilibili.com/x/passport-login/web/qrcode/generate"
	GetUserInfoUrl = "https://api.bilibili.com/x/web-interface/card?mid=%d"
)

// GetLoginKeyAndLoginUrl 获取二维码内容和密钥
func GetLoginKeyAndLoginUrl() (QrKey, QrUrl string) {

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

func GetUserInfo(mid int) UserInfo {
	u := fmt.Sprintf(GetUserInfoUrl, mid)
	client := http.Client{}
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data UserInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		return UserInfo{}
	}
	return data
}
