package models

import (
	"encoding/json"
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"time"
)

type SCNDuration time.Duration

func (t SCNDuration) MarshalToDB(v SCNDuration) (int64, error) {
	return v.Duration().Milliseconds(), nil
}

func (t SCNDuration) UnmarshalToModel(v int64) (SCNDuration, error) {
	return SCNDuration(timeext.FromMilliseconds(v)), nil
}

func (t SCNDuration) Duration() time.Duration {
	return time.Duration(t)
}

func (t *SCNDuration) UnmarshalJSON(data []byte) error {
	flt := float64(0)
	if err := json.Unmarshal(data, &flt); err != nil {
		return err
	}
	d0 := timeext.FromSeconds(flt)
	*t = SCNDuration(d0)
	return nil
}

func (t SCNDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Duration().Seconds())
}
