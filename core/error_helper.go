package core

import (
	"encoding/json"
	"fmt"
)

// TinyTuya Error Response Codes
const (
	ERR_JSON       = 900
	ERR_CONNECT    = 901
	ERR_TIMEOUT    = 902
	ERR_RANGE      = 903
	ERR_PAYLOAD    = 904
	ERR_OFFLINE    = 905
	ERR_STATE      = 906
	ERR_FUNCTION   = 907
	ERR_DEVTYPE    = 908
	ERR_CLOUDKEY   = 909
	ERR_CLOUDRESP  = 910
	ERR_CLOUDTOKEN = 911
	ERR_PARAMS     = 912
	ERR_CLOUD      = 913
	ERR_KEY_OR_VER = 914
)

var errorCodes = map[int]string{
	ERR_JSON:       "Invalid JSON Response from Device",
	ERR_CONNECT:    "Network Error: Unable to Connect",
	ERR_TIMEOUT:    "Timeout Waiting for Device",
	ERR_RANGE:      "Specified Value Out of Range",
	ERR_PAYLOAD:    "Unexpected Payload from Device",
	ERR_OFFLINE:    "Network Error: Device Unreachable",
	ERR_STATE:      "Device in Unknown State",
	ERR_FUNCTION:   "Function Not Supported by Device",
	ERR_DEVTYPE:    "Device22 Detected: Retry Command",
	ERR_CLOUDKEY:   "Missing Tuya Cloud Key and Secret",
	ERR_CLOUDRESP:  "Invalid JSON Response from Cloud",
	ERR_CLOUDTOKEN: "Unable to Get Cloud Token",
	ERR_PARAMS:     "Missing Function Parameters",
	ERR_CLOUD:      "Error Response from Tuya Cloud",
	ERR_KEY_OR_VER: "Check device key or version",
}

func ErrorJSON(number int, payload interface{}) map[string]interface{} {
	var spayload string
	b, err := json.Marshal(payload)
	if err != nil {
		spayload = "\"\""
	} else {
		spayload = string(b)
	}

	msg, ok := errorCodes[number]
	if !ok {
		msg = "Unknown Error"
	}

	jsonStr := fmt.Sprintf(`{ "Error":"%s", "Err":"%d", "Payload":%s }`, msg, number, spayload)

	var result map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &result)
	return result
}
