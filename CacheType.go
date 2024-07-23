package main

const (
	GetRoomDataUrl = "https://api.live.bilibili.com/xlive/web-room/v1/giftPanel/giftData?room_id=%d&platform=pc"
)

type GiftDataList struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		RoomGiftList struct {
			GoldList []struct {
				Position int `json:"position"`
				GiftId   int `json:"gift_id"`
				Id       int `json:"id"`
				PlanId   int `json:"plan_id"`
				Special  struct {
					SpecialType int    `json:"special_type"`
					IsUse       int    `json:"is_use"`
					Tips        string `json:"tips"`
				} `json:"special"`
				UpgradeGift []struct {
					GiftId  int    `json:"gift_id"`
					Alias   string `json:"alias"`
					Desc    string `json:"desc"`
					Locked  bool   `json:"locked"`
					LockTip string `json:"lock_tip"`
				} `json:"upgrade_gift"`
				GiftTag   []int `json:"gift_tag"`
				ExtraInfo struct {
					IsFixed        int    `json:"is_fixed"`
					IconBottomTips string `json:"icon_bottom_tips"`
					Alg            string `json:"alg"`
				} `json:"extra_info"`
				GiftScene interface{} `json:"gift_scene"`
			} `json:"gold_list"`
			SilverList []struct {
				Position int `json:"position"`
				GiftId   int `json:"gift_id"`
				Id       int `json:"id"`
				PlanId   int `json:"plan_id"`
				Special  struct {
					SpecialType int    `json:"special_type"`
					IsUse       int    `json:"is_use"`
					Tips        string `json:"tips"`
				} `json:"special"`
				UpgradeGift interface{} `json:"upgrade_gift"`
				GiftTag     interface{} `json:"gift_tag"`
				ExtraInfo   struct {
					IsFixed        int    `json:"is_fixed"`
					IconBottomTips string `json:"icon_bottom_tips"`
					Alg            string `json:"alg"`
				} `json:"extra_info"`
				GiftScene interface{} `json:"gift_scene"`
			} `json:"silver_list"`
			NeedOddsOffline bool `json:"need_odds_offline"`
			AbResult        struct {
				GiftFirstScreen string `json:"gift_first_screen"`
			} `json:"ab_result"`
		} `json:"room_gift_list"`
		DiscountGiftList []struct {
			GiftId         int    `json:"gift_id"`
			Price          int    `json:"price"`
			DiscountPrice  int    `json:"discount_price"`
			CornerMark     string `json:"corner_mark"`
			CornerPosition int    `json:"corner_position"`
			CornerColor    string `json:"corner_color"`
			Id             int    `json:"id"`
		} `json:"discount_gift_list"`
		TabList []struct {
			TabId    int    `json:"tab_id"`
			TabName  string `json:"tab_name"`
			Position int    `json:"position"`
			List     []struct {
				Position int `json:"position"`
				GiftId   int `json:"gift_id"`
				Id       int `json:"id"`
				PlanId   int `json:"plan_id"`
				Special  struct {
					SpecialType int    `json:"special_type"`
					IsUse       int    `json:"is_use"`
					Tips        string `json:"tips"`
				} `json:"special"`
				UpgradeGift interface{} `json:"upgrade_gift"`
				GiftTag     []int       `json:"gift_tag"`
				ExtraInfo   interface{} `json:"extra_info"`
				GiftScene   *struct {
					Scene   string `json:"scene"`
					PayType string `json:"pay_type"`
				} `json:"gift_scene"`
			} `json:"list"`
		} `json:"tab_list"`
		MaxSendGift       int         `json:"max_send_gift"`
		ComboIntervalTime int         `json:"combo_interval_time"`
		LotteryGiftConfig interface{} `json:"lottery_gift_config"`
		Privilege         struct {
			BuyGuardBtn   string `json:"buy_guard_btn"`
			IsExpired     int    `json:"is_expired"`
			PrivilegeType int    `json:"privilege_type"`
		} `json:"privilege"`
		RedDot []struct {
			TabId     int    `json:"tab_id"`
			HasRedDot int    `json:"has_red_dot"`
			RedDotId  int    `json:"red_dot_id"`
			RoomId    int    `json:"room_id"`
			Uid       int    `json:"uid"`
			Module    string `json:"module"`
		} `json:"red_dot"`
		PayLimitIcon    string      `json:"pay_limit_icon"`
		NamingGift      interface{} `json:"naming_gift"`
		SpecialShowGift interface{} `json:"special_show_gift"`
		BagTabDisable   int         `json:"bag_tab_disable"`
		SpecialTag      []struct {
			GiftId            int    `json:"gift_id"`
			CornerMark        string `json:"corner_mark"`
			CornerBackground  string `json:"corner_background"`
			HasBling          bool   `json:"has_bling"`
			SpecialGiftBanner struct {
				AppPic     string `json:"app_pic"`
				WebPic     string `json:"web_pic"`
				JumpUrl    string `json:"jump_url"`
				WebJumpUrl string `json:"web_jump_url"`
				HasSet     bool   `json:"has_set"`
			} `json:"special_gift_banner"`
			Corner struct {
				HasSet           bool   `json:"has_set"`
				Source           string `json:"source"`
				CornerMark       string `json:"corner_mark"`
				CornerMarkColor  string `json:"corner_mark_color"`
				CornerBackground string `json:"corner_background"`
				CornerColorBg    string `json:"corner_color_bg"`
				WebLight         struct {
					CornerMark       string `json:"corner_mark"`
					CornerBackground string `json:"corner_background"`
					CornerMarkColor  string `json:"corner_mark_color"`
					CornerColorBg    string `json:"corner_color_bg"`
				} `json:"web_light"`
				WebDark struct {
					CornerMark       string `json:"corner_mark"`
					CornerBackground string `json:"corner_background"`
					CornerMarkColor  string `json:"corner_mark_color"`
					CornerColorBg    string `json:"corner_color_bg"`
				} `json:"web_dark"`
			} `json:"corner"`
		} `json:"special_tag"`
	} `json:"data"`
}
