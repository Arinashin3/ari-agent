package convert

import "errors"

const (
	Bytes = iota + 1
	Kilobytes
	Megabytes
	Gigabytes
	Terabytes
	Petabytes
)

// BytesConvert
func BytesConvert(f float64, fromUnit int, toUnit int) float64 {
	if fromUnit < toUnit {
		for i := fromUnit - toUnit; i != 0; i++ {
			f /= 1024
		}
	} else if fromUnit > toUnit {
		for i := fromUnit - toUnit; i != 0; i-- {
			f *= 1024
		}
	}
	return f
}

func InterfaceToFloat64(t interface{}) (float64, error) {
	switch v := t.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case float32:
		return float64(v), nil
	default:
		return 0, errors.New("Unknown Type")

	}
}
