package main

import (
	"image/color"
)

const (
	// GuardLineType 舰长队列标识码
	GuardLineType = 0
	// GiftLineType 礼物队列标识码
	GiftLineType = 1
	// CommonLineType 普通队列标识码
	CommonLineType = 2

	// OpDelete 删除操作标识码
	OpDelete = 0
	// OpAdd 添加操作标识码
	OpAdd = 1
)

// RoomInfo 直播间信息
type RoomInfo struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		Uid              int      `json:"uid"`
		RoomId           int      `json:"room_id"`
		ShortId          int      `json:"short_id"`
		Attention        int      `json:"attention"`
		Online           int      `json:"online"`
		IsPortrait       bool     `json:"is_portrait"`
		Description      string   `json:"description"`
		LiveStatus       int      `json:"live_status"`
		AreaId           int      `json:"area_id"`
		ParentAreaId     int      `json:"parent_area_id"`
		ParentAreaName   string   `json:"parent_area_name"`
		OldAreaId        int      `json:"old_area_id"`
		Background       string   `json:"background"`
		Title            string   `json:"title"`
		UserCover        string   `json:"user_cover"`
		Keyframe         string   `json:"keyframe"`
		IsStrictRoom     bool     `json:"is_strict_room"`
		LiveTime         string   `json:"live_time"`
		Tags             string   `json:"tags"`
		IsAnchor         int      `json:"is_anchor"`
		RoomSilentType   string   `json:"room_silent_type"`
		RoomSilentLevel  int      `json:"room_silent_level"`
		RoomSilentSecond int      `json:"room_silent_second"`
		AreaName         string   `json:"area_name"`
		Pendants         string   `json:"pendants"`
		AreaPendants     string   `json:"area_pendants"`
		HotWords         []string `json:"hot_words"`
		HotWordsStatus   int      `json:"hot_words_status"`
		Verify           string   `json:"verify"`
		NewPendants      struct {
			Frame struct {
				Name       string `json:"name"`
				Value      string `json:"value"`
				Position   int    `json:"position"`
				Desc       string `json:"desc"`
				Area       int    `json:"area"`
				AreaOld    int    `json:"area_old"`
				BgColor    string `json:"bg_color"`
				BgPic      string `json:"bg_pic"`
				UseOldArea bool   `json:"use_old_area"`
			} `json:"frame"`
			Badge       interface{} `json:"badge"`
			MobileFrame struct {
				Name       string `json:"name"`
				Value      string `json:"value"`
				Position   int    `json:"position"`
				Desc       string `json:"desc"`
				Area       int    `json:"area"`
				AreaOld    int    `json:"area_old"`
				BgColor    string `json:"bg_color"`
				BgPic      string `json:"bg_pic"`
				UseOldArea bool   `json:"use_old_area"`
			} `json:"mobile_frame"`
			MobileBadge interface{} `json:"mobile_badge"`
		} `json:"new_pendants"`
		UpSession            string `json:"up_session"`
		PkStatus             int    `json:"pk_status"`
		PkId                 int    `json:"pk_id"`
		BattleId             int    `json:"battle_id"`
		AllowChangeAreaTime  int    `json:"allow_change_area_time"`
		AllowUploadCoverTime int    `json:"allow_upload_cover_time"`
		StudioInfo           struct {
			Status     int           `json:"status"`
			MasterList []interface{} `json:"master_list"`
		} `json:"studio_info"`
	} `json:"data"`
}

// LineRow 队列信息
type LineRow struct {
	GuardLine   []Line
	GiftLine    []GiftLine
	CommonLine  []Line
	GuardIndex  map[int]int
	GiftIndex   map[int]int
	CommonIndex map[int]int
}

// UpdateIndex 更新队列索引，传入UpdateType 0为GuardLine 1为GiftLine 2为CommonLine
func (r LineRow) UpdateIndex(UpdateType int) {
	switch UpdateType {
	case 0:
		for i, l := range r.GuardLine {
			r.GuardIndex[l.Uid] = i + 1
		}
	case 1:
		for i, l := range r.GiftLine {
			r.GiftIndex[l.Uid] = i + 1
		}
	case 2:
		for i, l := range r.CommonLine {
			r.CommonIndex[l.Uid] = i + 1
		}
	}
}

// Line 单一队列基础信息
type Line struct {
	Uid        int       `json:"Uid"`
	UserName   string    `json:"UserName"`
	Avatar     string    `json:"Avatar"`
	PrintColor LineColor `json:"PrintColor"`
}

// GiftLine 礼物用户队列信息
type GiftLine struct {
	Uid        int       `json:"Uid"`
	UserName   string    `json:"UserName"`
	Avatar     string    `json:"Avatar"`
	PrintColor LineColor `json:"PrintColor"`
	GiftPrice  float64   `json:"GiftPrice"`
}

// WsPack 前端通讯Websocket包结构
type WsPack struct {
	OpMessage int
	Index     int
	LineType  int
	Line      Line
	GiftLine  GiftLine
}

// RunConfig 配置格式
type RunConfig struct {
	IdCode                  string
	GuardPrintColor         LineColor
	GiftPrintColor          LineColor
	GiftLinePrice           float64
	CommonPrintColor        LineColor
	LineKey                 string
	GiftPriceDisplay        bool
	IsOnlyGift              bool
	AutoJoinGiftLine        bool
	TransparentBackground   bool
	CurrentQueueSizeDisplay bool
	MaxLineCount            int
}

type LineColor struct {
	R uint32
	G uint32
	B uint32
}

func (lc LineColor) ToRGBA() color.RGBA {
	return color.RGBA{
		R: uint8(lc.R),
		G: uint8(lc.G),
		B: uint8(lc.B),
		A: 255, // 设置透明度为不透明
	}
}

func (lc LineColor) IsEmpty() bool {
	return lc.R == 0 && lc.G == 0 && lc.B == 0
}

func (r LineRow) IsEmpty() bool {
	return len(r.GuardLine) == 0 &&
		len(r.GiftLine) == 0 &&
		len(r.CommonLine) == 0
}

// VersionSct 版本检查结构
type VersionSct struct {
	Version      string   `json:"version"`
	VersionCount int      `json:"versionCount"`
	UpdateDate   string   `json:"update_date"`
	Changelog    []string `json:"changelog"`
	UpdateUrl    string   `json:"update_url"`
}
