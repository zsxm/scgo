package date

import (
	"strconv"
	"strings"
	"time"
)

type Duration time.Duration

const (
	Default int64 = -62135596800

	YYMD    = "2006-01-02"
	YYMDHMS = "2006-01-02 15:04:05"
	YYMDHM  = "2006-01-02 15:04"
	GMT     = "Mon, 02 Jan 2006 15:04:05 GMT"

	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
	Day                  = 24 * Hour
)

func Now() time.Time {
	return time.Now()
}

func NowUnix() int64 {
	return time.Now().Unix()
}

func NowUnixStr() string {
	return strconv.FormatInt(NowUnix(), 10)
}

func Unix(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

func Format(layout string, t time.Time) string {
	if t.IsZero() {
		return ""
	} else {
		return t.Format(layout)
	}
}

func FormatGMT(t time.Time) string {
	return Format(GMT, t)
}

func FormatYYMD(t time.Time) string {
	return Format(YYMD, t)
}

func FormatYYMDHMS(t time.Time) string {
	return Format(YYMDHMS, t)
}

func FormatYYMDHM(t time.Time) string {
	return Format(YYMDHM, t)
}

func FormatUnix(layout string, i int64) string {
	if i == 0 || i == Default {
		return ""
	}
	t := time.Unix(i, 0)
	return Format(layout, t)
}

func FormatUnixGMT(i int64) string {
	return FormatUnix(GMT, i)
}

func FormatUnixYYMD(i int64) string {
	return FormatUnix(YYMD, i)
}

func FormatUnixYYMDHMS(i int64) string {
	return FormatUnix(YYMDHMS, i)
}

func FormatUnixYYMDHM(i int64) string {
	return FormatUnix(YYMDHM, i)
}

func FormatNumString(v string) string {
	str := strings.Replace(v, " ", "", -1)
	str = strings.Replace(str, ":", "", -1)
	str = strings.Replace(str, "-", "", -1)
	return str
}
