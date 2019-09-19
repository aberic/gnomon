package gnomon

import (
	"testing"
	"time"
)

func TestTimeCommon_String2Timestamp(t *testing.T) {
	i64, err := Time().String2Timestamp("2019/09/17 10:16:56", "2006/01/02 15:04:05", time.Local)
	t.Log("i64", i64, err)
	i64, err = Time().String2Timestamp("2019/09/17 10:16:56", "2006/01/02 15:04:05", time.UTC)
	t.Log("i64", i64, err)
}

func TestTimeCommon_Timestamp2String(t *testing.T) {
	t.Log(Time().Timestamp2String(1568686616, 28889, "2006/01/02 15:04:05"))
	t.Log(Time().Timestamp2String(1568686626, 98882, "2006/01/02 15:04:05"))
}
