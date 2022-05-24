package pgxmap

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

// validates driver.Valuer interface implementation.
var _ driver.Valuer = ConfigMap{}

// ConfigMap represents a map of key-value pairs that's used as input/output.
type ConfigMap map[string]interface{}

// Value implements the driver.Valuer interface.
func (m ConfigMap) Value() (driver.Value, error) {
	items := make([]Item, 0, len(m))
	for k, v := range m {
		var val interface{}
		transformed := false

		switch v.(type) {
		case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:
			val = fmt.Sprintf("%d", v)
			transformed = true
		case float64, float32:
			val = fmt.Sprintf("%f", v)
			transformed = true
		default:
			val = v
		}

		items = append(items, Item{
			Key:         k,
			Value:       val,
			Type:        reflect.TypeOf(v).Name(),
			Transformed: transformed,
		})
	}

	return json.Marshal(items)
}

// Scan implements the sql.Scanner interface.
func (m *ConfigMap) Scan(value interface{}) error {
	source, ok := value.([]byte)
	if !ok {
		return errors.New("invalid type assertion: .([]byte)")
	}

	var list []Item
	err := json.Unmarshal(source, &list)
	if err != nil {
		return errors.Wrap(err, "error unmarshalling json")
	}

	*m = make(map[string]interface{}, len(list))
	for _, item := range list {
		if !item.Transformed {
			(*m)[item.Key] = item.Value
			continue
		}

		switch item.Type {
		case "float64", "float32":
			val, err := strconv.ParseFloat(item.Value.(string), 64)
			if err != nil {
				return errors.Wrap(err, "error parsing float64")
			}
			if item.Type == "float32" {
				(*m)[item.Key] = float32(val)
				break
			}
			(*m)[item.Key] = val

		case "int64", "int32", "int16", "int8", "int", "uint64", "uint32", "uint16", "uint8", "uint":
			i, err := strconv.ParseInt(item.Value.(string), 10, 64)
			if err != nil {
				return errors.Wrap(err, "error parsing int64")
			}

			switch item.Type {
			case reflect.Int64.String():
				(*m)[item.Key] = i
			case reflect.Int32.String():
				(*m)[item.Key] = int32(i)
			case reflect.Int16.String():
				(*m)[item.Key] = int16(i)
			case reflect.Int8.String():
				(*m)[item.Key] = int8(i)
			case reflect.Int.String():
				(*m)[item.Key] = int(i)
			case reflect.Uint.String():
				(*m)[item.Key] = uint(i)
			case reflect.Uint64.String():
				(*m)[item.Key] = uint64(i)
			case reflect.Uint32.String():
				(*m)[item.Key] = uint32(i)
			case reflect.Uint16.String():
				(*m)[item.Key] = uint16(i)
			case reflect.Uint8.String():
				(*m)[item.Key] = uint8(i)
			case reflect.Uint.String():
				(*m)[item.Key] = uint(i)
			default:
				(*m)[item.Key] = item.Value
			}
		default:
			(*m)[item.Key] = item.Value
		}
	}

	return nil
}
