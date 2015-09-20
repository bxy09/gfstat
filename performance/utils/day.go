package utils

import "time"

//DayInBJ 时间对应的北京时间日期序号
func DayInBJ(t time.Time) int64 {
	return (t.Unix() + 8*3600) / (24 * 3600)
}

//StartOfBJTime 获取当天北京时间零点
func StartOfBJTime(t time.Time) time.Time {
	return time.Unix(DayInBJ(t)*24*3600-8*3600, 0)
}

//StartOfBJTimeI 从北京时间日期序号获得时间零点
func StartOfBJTimeI(day int64) time.Time {
	return time.Unix(day*24*3600-8*3600, 0)
}
