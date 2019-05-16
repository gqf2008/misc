package misc

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pborman/uuid"
)

//Exist ....
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//UUID ....
func UUID() string {
	return strings.Replace(uuid.NewUUID().String(), "-", "", -1)
}

//TimeFormat ....
func TimeFormat(format string, t time.Time) string {
	//2006-01-02T15:04:05
	//yy/MM/dd HH:mm:ss
	format = strings.Replace(format, "yyyy", "2006", -1)
	format = strings.Replace(format, "yy", "06", -1)
	format = strings.Replace(format, "MM", "01", -1)
	format = strings.Replace(format, "M", "1", -1)
	format = strings.Replace(format, "dd", "02", -1)
	format = strings.Replace(format, "d", "2", -1)
	format = strings.Replace(format, "HH", "15", -1)
	format = strings.Replace(format, "mm", "04", -1)
	format = strings.Replace(format, "m", "4", -1)
	format = strings.Replace(format, "ss", "05", -1)
	format = strings.Replace(format, "s", "5", -1)
	return t.Format(format)
}

//GetString ....
func GetString(m map[string]interface{}, key string) (string, bool) {
	v, found := m[key]
	if !found {
		return "", false
	}
	if vv, ok := v.(string); ok {
		return vv, true
	}
	return "", false
}

//GetInt ....
func GetInt(m map[string]interface{}, key string) (int64, bool) {
	v, found := m[key]
	if !found {
		return 0, false
	}
	if vv, ok := v.(float64); ok {
		return int64(vv), true
	}
	if vv, ok := v.(int64); ok {
		return vv, true
	}
	return 0, false
}

//GetFloat ....
func GetFloat(m map[string]interface{}, key string) (float64, bool) {
	v, found := m[key]
	if !found {
		return 0, false
	}
	if vv, ok := v.(float64); ok {
		return vv, true
	}
	if vv, ok := v.(int64); ok {
		return float64(vv), true
	}
	return 0.0, false
}

//GetObject ....
func GetObject(m map[string]interface{}, key string, obj interface{}) bool {
	v, found := m[key]
	if !found {
		return false
	}
	b, _ := json.Marshal(v)
	err := json.Unmarshal(b, obj)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

//GetMap ....
func GetMap(m map[string]interface{}, key string) (map[string]interface{}, bool) {
	v, found := m[key]
	if !found {
		return nil, false
	}
	if vv, ok := v.(map[string]interface{}); ok {
		return vv, true
	}
	return nil, false
}

//GetArray ....
func GetArray(m map[string]interface{}, key string) ([]interface{}, bool) {
	v, found := m[key]
	if !found {
		return nil, false
	}
	if vv, ok := v.([]interface{}); ok {
		return vv, true
	}
	return nil, false
}
