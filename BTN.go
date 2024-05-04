package main

import (
	"encoding/json"
	"github.com/tidwall/jsonc"
)

type BTN_Ability struct {
	Interval     uint32 `json:"interval"`
	Endpoint     string `json:"endpoint"`
	RandomDelay  uint32 `json:"random_initial_delay"`
	Version      string `json:"version"`
}
type BTN_ConfigStruct struct {
	MinMainVersion uint32                 `json:"min_protocol_version"`
	MaxMainVersion uint32                 `json:"max_protocol_version"`
	Ability        map[string]BTN_Ability `json:"ability"`
}

var btnProtocol = "BTN-Protocol/0.0.0-dev"
var btnUserAgent = programUserAgent + " " + btnProtocol
var btnHeader = map[string]string {
	"User-Agent": btnUserAgent,
}

var btn_lastGetConfig int64 = 0
var btn_configureInterval = 60

var btnConfig *BTN_ConfigStruct

func BTN_GetConfig() bool {
	if config.BTNConfigureURL == "" || (btn_lastGetConfig + int64(btn_configureInterval)) > currentTimestamp  {
		return true
	}

	Log("Debug-BTN_GetConfig", "In progress..", false)

	btn_lastGetConfig = currentTimestamp

	_, _, btnConfigContent := Fetch(config.BTNConfigureURL, false, false, &btnHeader)
	if btnConfigContent == nil {
		Log("BTN_GetConfig", GetLangText("Error-FetchResponse"), true)
		return false
	}

	// Max 8MB.
	if len(btnConfigContent) > 8388608 {
		Log("BTN_GetConfig", GetLangText("Error-LargeFile"), true)
		return false
	}

	if err := json.Unmarshal(jsonc.ToJSON(btnConfigContent), btnConfig); err != nil {
		Log("SyncWithServer", GetLangText("Error-ParseConfig"), true, err.Error())
		return false
	}

	return true
}
