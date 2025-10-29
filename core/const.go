package core

// Globals Network Settings
const (
	MAXCOUNT      = 15
	SCANTIME      = 18
	UDPPORT       = 6666
	UDPPORTS      = 6667
	UDPPORTAPP    = 7000
	TCPPORT       = 6668
	TIMEOUT       = 3.0
	TCPTIMEOUT    = 0.4
	DEFAULT_NETWORK = "192.168.0.0/24"
)

// Configuration Files
const (
	CONFIGFILE   = "tinytuya.json"
	DEVICEFILE   = "devices.json"
	RAWFILE      = "tuya-raw.json"
	SNAPSHOTFILE = "snapshot.json"
)

var DEVICEFILE_SAVE_VALUES = []string{"category", "product_name", "product_id", "biz_type", "model", "sub", "icon", "version", "last_ip", "uuid", "node_id", "sn", "mapping"}
