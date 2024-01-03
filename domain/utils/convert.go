package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	constants "github.com/pharmaniaga/auth-user/domain/constants/general"
)

func GetInt(x string) int {
	i, err := strconv.Atoi(x)
	if err != nil {
		fmt.Println("utils -> GentInt : error, ", err)
		fmt.Println("Can't convert into Integer")
		fmt.Println("Please re-check .env, you recently input")
		fmt.Println(x)
	}

	return i
}

func GetBool(x string) bool {
	i, err := strconv.ParseBool(x)
	if err != nil {
		fmt.Println("utils -> GetBool : error, ", err)
		fmt.Println("Can't convert into Boolean")
		fmt.Println("Please re-check .env, you recently input")
		fmt.Println(x)
	}

	return i
}

func GetFloat(x string) float32 {
	i, err := strconv.ParseFloat(x, 32)
	if err != nil {
		fmt.Println("utils -> GetFloat : error, ", err)
		fmt.Println("Can't convert into float32")
		fmt.Println("Please re-check .env, you recently input")
		fmt.Println(x)
	}

	return float32(i)
}

func ToFormatTime(datetime string) (string, error) {
	t, err := time.Parse(constants.DBTimeLayout, datetime)
	if err != nil {
		return datetime, err
	}

	tString := t.Format(constants.ResponseTimeLayout)

	return tString, nil
}

func GetTimeString() string {
	t := time.Now().UTC()
	tString := t.Format(constants.DBTimeLayout)

	return tString
}

func StrToInt(data string) (int, error) {
	return strconv.Atoi(data)
}

func StrToBool(data string) (bool, error) {
	return strconv.ParseBool(data)
}

func StrToInt64(data string) (int64, error) {
	return strconv.ParseInt(data, 10, 64)
}

func StrToFloat64(data string) (float64, error) {
	return strconv.ParseFloat(data, 64)
}

func Int64sJoin(data []int64) string {
	s, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return strings.Trim(string(s), "[]")
}

func GetDataFromKey(data, key string) (string, error) {
	if data == "" || key == "" {
		return "", errors.New("data/key cannot be empty")
	}

	result, err := GetDecrypt([]byte(key), data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result), nil
}

func GetKeyData(data, key string) (string, error) {
	if data == "" || key == "" {
		return "", errors.New("data/key cannot be empty")
	}

	result, err := GetEncrypt([]byte(key), fmt.Sprintf("%v", data))
	if err != nil {
		return "", err
	}

	return result, nil
}

func ConvertMonthtoRoman(month int) string {
	switch month {
	case constants.NumJan:
		return constants.RomanJan
	case constants.NumFeb:
		return constants.RomanFeb
	case constants.NumMar:
		return constants.RomanMar
	case constants.NumApr:
		return constants.RomanApr
	case constants.NumMay:
		return constants.RomanMay
	case constants.NumJune:
		return constants.RomanJune
	case constants.NumJuly:
		return constants.RomanJuly
	case constants.NumAug:
		return constants.RomanAug
	case constants.NumSep:
		return constants.RomanSep
	case constants.NumOct:
		return constants.RomanOct
	case constants.NumNov:
		return constants.RomanNov
	case constants.NumDec:
		return constants.RomanDec
	}

	return ""
}

func ConvertMonthtoString(month int) string {
	switch month {
	case constants.NumJan:
		return "January"
	case constants.NumFeb:
		return "February"
	case constants.NumMar:
		return "March"
	case constants.NumApr:
		return "April"
	case constants.NumMay:
		return "May"
	case constants.NumJune:
		return "June"
	case constants.NumJuly:
		return "July"
	case constants.NumAug:
		return "August"
	case constants.NumSep:
		return "September"
	case constants.NumOct:
		return "October"
	case constants.NumNov:
		return "November"
	case constants.NumDec:
		return "December"
	}

	return ""
}

func StructToString(data interface{}) string {
	result, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return string(result)
}

func FloatToRupiah(price float64) string {
	p := message.NewPrinter(language.Indonesian)
	moneyString := p.Sprintf("Rp %.2f", price)

	return moneyString
}

func ConvertIDs(ids string) ([]int64, error) {
	var result []int64
	idString := strings.Split(ids, ",")

	for _, val := range idString {
		id, err := StrToInt64(val)
		if err != nil {
			return result, err
		}

		result = append(result, id)
	}

	return result, nil
}

func ArrInt64Join(ids []int64, separator string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(ids), " ", separator, -1), "[]")
}

func StrToArrInt64(data, separator string) ([]int64, error) {
	splitData := strings.Split(data, separator)

	var result []int64
	for _, dt := range splitData {
		convData, err := StrToInt64(dt)
		if err != nil {
			return result, err
		}

		result = append(result, convData)
	}

	return result, nil
}

func StrToArrMapInt64(data, separator string) (map[int64]int64, error) {
	splitData := strings.Split(data, separator)

	result := make(map[int64]int64)
	for _, dt := range splitData {
		convData, err := StrToInt64(dt)
		if err != nil {
			return result, err
		}

		result[convData] = convData
	}

	return result, nil
}

func StrToArrMapString(data, separator string) (map[string]string, error) {
	splitData := strings.Split(data, separator)

	result := make(map[string]string)
	for _, dt := range splitData {
		result[dt] = dt
	}

	return result, nil
}

func FormatPhoneNumber(phone string) string {
	// Remove any leading spaces or plus sign
	phone = strings.TrimLeft(phone, " +")

	// Check if the phone number starts with "62"
	if strings.HasPrefix(phone, "62") {
		return phone
	}

	// Check if the phone number starts with "+62"
	if strings.HasPrefix(phone, "+62") {
		// Remove the "+" character and return the number
		return strings.TrimPrefix(phone, "+")
	}

	if strings.Contains(phone, "0") {
		// Replace the "0" with "62" as the prefix and return the number
		return fmt.Sprintf("62%s", strings.Replace(phone, "0", "", 1))
	}

	// If the phone number doesn't start with "62" or "+62", add "62" as the prefix
	return fmt.Sprintf("62%s", phone)
}
