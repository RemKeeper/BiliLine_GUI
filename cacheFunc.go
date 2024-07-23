package main

import (
	"encoding/json"
	"fmt"
	"github.com/vtb-link/bianka/proto"
	"io"
	"net/http"
)

var DiscountGiftData GiftDataList

func GetRoomGiftData(RoomId int) {
	DataResp, err := http.Get(fmt.Sprintf(GetRoomDataUrl, RoomId))
	if err != nil {
		return
	}
	GiftBody, err := io.ReadAll(DataResp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(GiftBody, &DiscountGiftData)
	if err != nil {
		return
	}
}

func FindAndModifyDiscountGift(LiveGiftData *proto.CmdSendGiftData) {
	for _, disCountGift := range DiscountGiftData.Data.DiscountGiftList {
		if disCountGift.GiftId == LiveGiftData.GiftID {
			LiveGiftData.Price = disCountGift.DiscountPrice
		}
	}
}
