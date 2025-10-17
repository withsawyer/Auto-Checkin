package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
)

func JsonToUrlValues(jsonStr string) (url.Values, error) {
	var data map[string]any
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	values, err := Map2UrlValues(data)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func Map2UrlValues(input map[string]any) (url.Values, error) {
	values := url.Values{}
	for key, value := range input {
		switch v := value.(type) {
		case string:
			values.Add(key, v)
		case int, int8, int16, int32, int64:
			values.Add(key, fmt.Sprintf("%d", v))
		case uint, uint8, uint16, uint32, uint64:
			values.Add(key, fmt.Sprintf("%d", v))
		case float32, float64:
			values.Add(key, fmt.Sprintf("%f", v))
		case bool:
			values.Add(key, fmt.Sprintf("%t", v))
		case []any:
			for _, item := range v {
				strItem, err := Any2String(item)
				if err != nil {
					return nil, err
				}
				values.Add(key, strItem)
			}
		case map[string]any:
			// 如果嵌套的是另一个map，可以递归处理
			nestedValues, err := Map2UrlValues(v)
			if err != nil {
				return nil, err
			}
			for nestedKey, nestedValuesList := range nestedValues {
				for _, nestedValue := range nestedValuesList {
					values.Add(fmt.Sprintf("%s[%s]", key, nestedKey), nestedValue)
				}
			}
		default:
			return nil, fmt.Errorf("unsupported type for key %s: %T", key, value)
		}
	}

	return values, nil
}

func Any2String(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v), nil
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32, float64:
		return fmt.Sprintf("%f", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	default:
		// 尝试将其他类型的值序列化为JSON字符串
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("unsupported type: %T", value)
		}
		return string(jsonBytes), nil
	}
}

// StringToInt 将字符串转换为 int 类型，并返回错误信息
func StringToInt(s string) (int, error) {
	if s == "" {
		return 0, errors.New("输入为空")
	}

	// 使用 ParseFloat 解析字符串为 float64
	fnum, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("无法将字符串转换为数字: %w", err)
	}
	// 检查浮点数是否为有效整数
	if fnum != math.Trunc(fnum) {
		return 0, fmt.Errorf("字符串表示的数字不是有效的整数: %s", s)
	}
	// 转换为 int 类型
	inum := int(fnum)
	return inum, nil
}

// StringToFloat 将字符串转换为 float64 类型，并返回错误信息
func StringToFloat(s string) (float64, error) {
	if s == "" {
		return 0.0, errors.New("输入为空")
	}
	fnum, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, fmt.Errorf("无法将字符串转换为浮点数: %w", err)
	}
	return fnum, nil
}
