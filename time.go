package gnomon

import "time"

type timeCommon struct{}

// String2Timestamp 字符串转时间戳
//
// date 待转换时间字符串 如：2019/09/17 10:16:56
//
// format 时间字符串格式化类型 如：2006/01/02 15:04:05
//
// zone 时区 如：time.Local / time.UTC
func (t *timeCommon) String2Timestamp(date, format string, zone *time.Location) (int64, error) {
	var (
		theTime time.Time
		err     error
	)
	if theTime, err = time.ParseInLocation(format, date, zone); nil != err {
		return 0, err
	}
	return theTime.Unix(), nil
}

// Timestamp2String 时间戳转字符串
//
// timestampSec 时间戳秒值
//
// timestampNSec 时间戳纳秒值
//
// format 时间字符串格式化类型 如：2006/01/02 15:04:05
func (t *timeCommon) Timestamp2String(timestampSec, timestampNSec int64, format string) string {
	return time.Unix(timestampSec, timestampNSec).Format(format) //设置时间戳 使用模板格式化为日期字符串
}
