package convert

import (
	"errors"
	"strings"

	"github.com/alecthomas/units"
)

const (
	Bytes     = iota + 1
	Kilobytes = Bytes * 1024
	Megabytes = Kilobytes * 1024
	Gigabytes = Megabytes * 1024
	Terabytes = Gigabytes * 1024
	Petabytes = Terabytes * 1024
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

func UnitConvert(s string, toUnit string) float64 {
	v, _ := units.ParseBase2Bytes(s)
	toUnit = strings.ToLower(toUnit)
	switch toUnit {
	case "byte", "b":
		return float64(v) / Bytes
	case "kilobytes", "kb", "kib":
		return float64(v) / Kilobytes
	case "megabytes", "mb", "mib":
		return float64(v) / Megabytes
	case "gigabytes", "gb", "gib":
		return float64(v) / Gigabytes
	case "terabytes", "tb", "tib":
		return float64(v) / Terabytes
	case "petabytes", "pb", "pib":
		return float64(v) / Petabytes
	}
	return -1
}
