package model

import (
	"database/sql/driver"
	"time"
)

type DateTime time.Time

func (t *DateTime) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		*t = DateTime(time.Time{})
		return
	}
	// 指定解析的格式
	if now, err := time.ParseInLocation(`"`+TimeLayout+`"`, string(data), time.Local); err == nil {
		*t = DateTime(now)
	} else {
		if now1, err1 := time.ParseInLocation(`"`+TimeLayoutShort+`"`, string(data), time.Local); err1 == nil {
			*t = DateTime(now1)
		} else {
			if now2, err2 := time.ParseInLocation(`"`+DateLayoutShort+`"`, string(data), time.Local); err2 == nil {
				*t = DateTime(now2)
			}
		}
	}
	return
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeLayout)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeLayout)
	b = append(b, '"')
	return b, nil
}

func (t DateTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(TimeLayout)), nil
}

func (t *DateTime) Scan(v interface{}) error {
	tTime, _ := time.Parse("2006-01-02 15:04:05 +0800 CST", v.(time.Time).String())
	*t = DateTime(tTime)
	return nil
}

func (t DateTime) String() string {
	return time.Time(t).Format(TimeLayout)
}

func (t DateTime) Time() time.Time {
	return time.Time(t)
}

var TimeLayout string = "2006-01-02 15:04:05"

var TimeLayoutShort string = "2006-01-02"

var DateLayoutShort string = "20060102"
