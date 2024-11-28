package main

import (
	"encoding/json"
	"io"
	"net/http"
)

const updateUrl = "https://lineupversion.rem.asia/"

const (
	NowVersion      = "1.4.3"
	NowVersionCount = 43
)

func CheckVersion() (VersionSct, bool) {
	get, err := http.Get(updateUrl)
	if err != nil {
		return VersionSct{}, false
	}
	all, err := io.ReadAll(get.Body)
	if err != nil {
		return VersionSct{}, false
	}
	var VersionCache VersionSct
	err = json.Unmarshal(all, &VersionCache)
	if err != nil {
		return VersionSct{}, false
	}
	if VersionCache.VersionCount > NowVersionCount {
		return VersionCache, true
	} else {
		return VersionCache, false
	}
}
