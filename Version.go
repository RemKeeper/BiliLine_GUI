package main

import (
	"BiliLine_Windows/Global"
	"encoding/json"
	"io"
	"net/http"
)

const updateUrl = "https://lineupversion.rem.asia/"

const (
	NowVersion      = "1.3.0"
	NowVersionCount = 36
)

func CheckVersion() (GlobalType.VersionSct, bool) {
	get, err := http.Get(updateUrl)
	if err != nil {
		return GlobalType.VersionSct{}, false
	}
	all, err := io.ReadAll(get.Body)
	if err != nil {
		return GlobalType.VersionSct{}, false
	}
	var VersionCache GlobalType.VersionSct
	err = json.Unmarshal(all, &VersionCache)
	if err != nil {
		return GlobalType.VersionSct{}, false
	}
	if VersionCache.VersionCount > NowVersionCount {
		return VersionCache, true
	} else {
		return VersionCache, false
	}
}
