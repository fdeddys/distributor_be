package util

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func GetCurrDate() time.Time {

	return time.Now()
}

func GetCurrFormatDate() time.Time {

	trxDate := time.Now().Format("2006-01-02 00:00:00")
	date, err := time.Parse("2006-01-02 00:00:00", trxDate)
	fmt.Println("date ==>", date)
	if err != nil {
		fmt.Println("error ", err)
	}
	return date
}

func GetCurrFormatDateTime() string {

	return time.Now().Format("2006-01-02 15:04:05")
}

func ToSliceData(input interface{}) [][]string {
	var records [][]string
	var header []string
	object := reflect.ValueOf(input)

	if object.Len() > 0 {
		first := object.Index(0)
		typ := first.Type()

		for i := 0; i < first.NumField(); i++ {
			header = append(header, typ.Field(i).Name)
		}
		records = append(records, header)
	}

	var items []interface{}
	for i := 0; i < object.Len(); i++ {
		items = append(items, object.Index(i).Interface())
	}

	for _, v := range items {
		item := reflect.ValueOf(v)
		var record []string
		for i := 0; i < item.NumField(); i++ {
			itm := item.Field(i).Interface()
			record = append(record, fmt.Sprintf("%v", itm))
		}
		records = append(records, record)
	}
	return records
}

func Atoi64(s string) int64 {
	i64, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return 0
	}
	return i64
}

func AtoFloat64(s string) float64 {
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f64
}

func AtoFloat32(s string) float32 {
	f32, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0
	}
	return float32(f32)
}
