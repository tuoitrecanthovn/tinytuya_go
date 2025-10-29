package core

// Tuya Command Types
const (
	AP_CONFIG           = 1  // FRM_TP_CFG_WF
	ACTIVE              = 2  // FRM_TP_ACTV (discard)
	SESS_KEY_NEG_START  = 3  // FRM_SECURITY_TYPE3
	SESS_KEY_NEG_RESP   = 4  // FRM_SECURITY_TYPE4
	SESS_KEY_NEG_FINISH = 5  // FRM_SECURITY_TYPE5
	UNBIND              = 6  // FRM_TP_UNBIND_DEV
	CONTROL             = 7  // FRM_TP_CMD
	STATUS              = 8  // FRM_TP_STAT_REPORT
	HEART_BEAT          = 9  // FRM_TP_HB
	DP_QUERY            = 0x0a // 10
	QUERY_WIFI          = 0x0b // 11
	TOKEN_BIND          = 0x0c // 12
	CONTROL_NEW         = 0x0d // 13
	ENABLE_WIFI         = 0x0e // 14
	WIFI_INFO           = 0x0f // 15
	DP_QUERY_NEW        = 0x10 // 16
	SCENE_EXECUTE       = 0x11 // 17
	UPDATEDPS           = 0x12 // 18
	UDP_NEW             = 0x13 // 19
	AP_CONFIG_NEW       = 0x14 // 20
	BOARDCAST_LPV34     = 0x23 // 35
	REQ_DEVINFO         = 0x25 // 37
	LAN_EXT_STREAM      = 0x40 // 64
)
