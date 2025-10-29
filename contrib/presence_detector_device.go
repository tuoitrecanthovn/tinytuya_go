package contrib

import (
	"encoding/json"

	"tinytuya_go/core"
)

const (
	PRESENCE_DPS_KEY                  = "dps"
	PRESENCE_DPS_PRESENCE_KEY          = "1"
	PRESENCE_DPS_SENSITIVITY_KEY       = "2"
	PRESENCE_DPS_NEAR_DETECTION_KEY    = "3"
	PRESENCE_DPS_FAR_DETECTION_KEY     = "4"
	PRESENCE_DPS_AUTO_DETECT_RESULT_KEY = "6"
	PRESENCE_DPS_TARGET_DISTANCE_KEY   = "9"
	PRESENCE_DPS_DETECTION_DELAY_KEY   = "101"
	PRESENCE_DPS_FADING_TIME_KEY       = "102"
	PRESENCE_DPS_LIGHT_SENSE_KEY       = "104"
)

// PresenceDetectorDevice represents a Tuya-based Presence Detector.
type PresenceDetectorDevice struct {
	*core.Device
}

// StatusJSON returns a JSON string of the device status with human-readable labels.
func (d *PresenceDetectorDevice) StatusJSON() (string, error) {
	status, err := d.Status()
	if err != nil {
		return "", err
	}
	dps, ok := status[PRESENCE_DPS_KEY].(map[string]interface{})
	if !ok {
		return "", nil
	}

	jsonBytes, err := json.Marshal(map[string]interface{}{
		"Presence":        dps[PRESENCE_DPS_PRESENCE_KEY],
		"Sensitivity":     dps[PRESENCE_DPS_SENSITIVITY_KEY],
		"Near detection":  dps[PRESENCE_DPS_NEAR_DETECTION_KEY],
		"Far detection":   dps[PRESENCE_DPS_FAR_DETECTION_KEY],
		"Checking result": dps[PRESENCE_DPS_AUTO_DETECT_RESULT_KEY],
		"Target distance": dps[PRESENCE_DPS_TARGET_DISTANCE_KEY],
		"Detection delay": dps[PRESENCE_DPS_DETECTION_DELAY_KEY],
		"Fading time":     dps[PRESENCE_DPS_FADING_TIME_KEY],
		"Light sense":     dps[PRESENCE_DPS_LIGHT_SENSE_KEY],
	})
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// GetPresenceState returns the presence state of the Presence Detector.
func (d *PresenceDetectorDevice) GetPresenceState() (string, error) {
	status, err := d.Status()
	if err != nil {
		return "", err
	}
	state, _ := status[PRESENCE_DPS_KEY].(map[string]interface{})[PRESENCE_DPS_PRESENCE_KEY].(string)
	return state, nil
}
