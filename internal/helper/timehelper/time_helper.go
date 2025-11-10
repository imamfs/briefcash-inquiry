package timehelper

import "time"

const ISOLayoutWithMillisAndTimezone = "2006-01-02T15:04:05.000-07.00"

func FormatTimeToISO7(t time.Time) string {
	location := time.FixedZone("WIB", 7*60*60)
	return t.In(location).Format(ISOLayoutWithMillisAndTimezone)
}

func FormatISO7ToTime(value string) (time.Time, error) {
	return time.Parse(ISOLayoutWithMillisAndTimezone, value)
}
