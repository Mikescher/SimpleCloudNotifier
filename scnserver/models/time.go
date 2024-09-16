package models

import (
	"encoding/json"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/rfctime"
	"time"
)

type SCNTime time.Time

func (t SCNTime) MarshalToDB(v SCNTime) (int64, error) {
	return v.Time().UnixMilli(), nil
}

func (t SCNTime) UnmarshalToModel(v int64) (SCNTime, error) {
	return NewSCNTime(time.UnixMilli(v)), nil
}

func (t SCNTime) Time() time.Time {
	return time.Time(t)
}

func (t *SCNTime) UnmarshalJSON(data []byte) error {
	str := ""
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	t0, err := time.Parse(time.RFC3339Nano, str)
	if err != nil {
		return err
	}
	*t = SCNTime(t0)
	return nil
}

func (t SCNTime) MarshalJSON() ([]byte, error) {
	str := t.Time().Format(time.RFC3339Nano)
	return json.Marshal(str)
}

func NewSCNTime(t time.Time) SCNTime {
	return SCNTime(t)
}

func NewSCNTimePtr(t *time.Time) *SCNTime {
	if t == nil {
		return nil
	}
	return langext.Ptr(SCNTime(*t))
}

func NowSCNTime() SCNTime {
	return SCNTime(time.Now())
}

func tt(v rfctime.AnyTime) time.Time {
	if r, ok := v.(time.Time); ok {
		return r
	}
	if r, ok := v.(rfctime.RFCTime); ok {
		return r.Time()
	}
	return time.Unix(0, v.UnixNano()).In(v.Location())
}
