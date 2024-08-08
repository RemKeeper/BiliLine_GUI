package main

import (
	"BiliLine_Windows/BiliUtils"
	"encoding/json"
	"fmt"
	"image/color"
	"os"
)

const (
	ConfigPath     = "./lineConfig_black.json"
	LineConfigPath = "./line_black.json"
)

func GetConfig() (rConfig RunConfig, err error) {
	file, err := os.ReadFile(ConfigPath)
	if err != nil {
		return RunConfig{}, err
	}
	var Config RunConfig
	err = json.Unmarshal(file, &Config)
	if err != nil {
		return RunConfig{}, err
	}
	return Config, err
}

func GetBiliCookie() (BiliUtils.BiliCookieConfig, error) {
	file, err := os.ReadFile(BiliUtils.CookiePath)
	if err != nil {
		return BiliUtils.BiliCookieConfig{}, err
	}
	var CookieConfig BiliUtils.BiliCookieConfig
	err = json.Unmarshal(file, &CookieConfig)
	if err != nil {
		return BiliUtils.BiliCookieConfig{}, err
	}
	return CookieConfig, err
}

func SetConfig(sConfig RunConfig) bool {
	ConfigJson, _ := json.MarshalIndent(sConfig, "", " ")
	lineupConfig := ConfigPath
	_, ReadConfigErr := os.Open(lineupConfig)
	if ReadConfigErr != nil {
		fmt.Println("配置文件不存在，尝试创建")
		_, ConfigErr := os.Create(lineupConfig)
		_, LineCreate := os.Create(LineConfigPath)
		if ConfigErr != nil || LineCreate != nil {
			fmt.Println("配置文件创建失败", ConfigErr.Error(), LineCreate.Error())
			return false
		} else {
			err := os.WriteFile(lineupConfig, ConfigJson, 0o666)
			if err != nil {
				fmt.Println("配置文件更新失败", err.Error())
				return false
			} else {
				return true
			}
		}
	} else {
		err := os.WriteFile(lineupConfig, ConfigJson, 0o666)
		if err != nil {
			fmt.Println("配置文件更新失败", err.Error())
			return false
		} else {
			return true
		}
	}
}

func SetLine(lp LineRow) {
	lineJson, _ := json.MarshalIndent(lp, "", " ")
	lineConfigFile := LineConfigPath
	WriteErr := os.WriteFile(lineConfigFile, lineJson, 0666)
	if WriteErr != nil {
		fmt.Println("队列文件更新失败")
	}
}

func GetLine() (line LineRow, err error) {
	lineConfigFile := LineConfigPath
	var LineGet LineRow
	file, err := os.ReadFile(lineConfigFile)
	if err != nil {
		return LineRow{}, err
	} else {
		_ = json.Unmarshal(file, &LineGet)
		return LineGet, nil
	}
}

func ToLineColor(c color.Color) LineColor {
	r, g, b, _ := c.RGBA()
	return LineColor{
		R: r,
		G: g,
		B: b,
	}
}
