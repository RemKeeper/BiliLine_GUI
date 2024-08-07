package main

import (
	"encoding/json"
	"os"

	"github.com/vtb-link/bianka/proto"
)

var DiscountGiftData GiftDataList

func GetRoomGiftData(RoomId int) {
	//DataResp, err := http.Get(fmt.Sprintf(GetRoomDataUrl, RoomId))
	//if err != nil {
	//	return
	//}
	//GiftBody, err := io.ReadAll(DataResp.Body)
	//if err != nil {
	//	return
	//}
	file, err := os.ReadFile(GiftJsonPath)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &DiscountGiftData)
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
