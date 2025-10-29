package contrib

import (
	"fmt"
	"strconv"
	"strings"

	"tinytuya_go/core"
)

const (
	BLANKET_DPS_BODY_LEVEL     = "14"
	BLANKET_DPS_FEET_LEVEL     = "15"
	BLANKET_DPS_BODY_TIME      = "16"
	BLANKET_DPS_FEET_TIME      = "17"
	BLANKET_DPS_BODY_COUNTDOWN = "18"
	BLANKET_DPS_FEET_COUNTDOWN = "19"
	BLANKET_LEVEL_PREFIX       = "level_"
)

// BlanketDevice represents a Tuya based Electric Blanket Device.
type BlanketDevice struct {
	*core.Device
}

func (d *BlanketDevice) numberToLevel(num int) string {
	return fmt.Sprintf("%s%d", BLANKET_LEVEL_PREFIX, num+1)
}

func (d *BlanketDevice) levelToNumber(level string) (int, error) {
	numStr := strings.TrimPrefix(level, BLANKET_LEVEL_PREFIX)
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}
	return num - 1, nil
}

// GetFeetLevel returns the feet level.
func (d *BlanketDevice) GetFeetLevel() (int, error) {
	status, err := d.Status()
	if err != nil {
		return 0, err
	}
	level, _ := status["dps"].(map[string]interface{})[BLANKET_DPS_FEET_LEVEL].(string)
	return d.levelToNumber(level)
}

// GetBodyLevel returns the body level.
func (d *BlanketDevice) GetBodyLevel() (int, error) {
	status, err := d.Status()
	if err != nil {
		return 0, err
	}
	level, _ := status["dps"].(map[string]interface{})[BLANKET_DPS_BODY_LEVEL].(string)
	return d.levelToNumber(level)
}

// SetFeetLevel sets the feet level.
func (d *BlanketDevice) SetFeetLevel(num int) (map[string]interface{}, error) {
	if num < 0 || num > 6 {
		return nil, fmt.Errorf("level needs to be between 0 and 6")
	}
	return d.SetValue(15, d.numberToLevel(num))
}

// SetBodyLevel sets the body level.
func (d *BlanketDevice) SetBodyLevel(num int) (map[string]interface{}, error) {
	if num < 0 || num > 6 {
		return nil, fmt.Errorf("level needs to be between 0 and 6")
	}
	return d.SetValue(14, d.numberToLevel(num))
}
