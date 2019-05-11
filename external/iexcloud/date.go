package iexcloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Date models a report date
// copied from https://github.com/goinvest/iexcloud
type Date time.Time

// UnmarshalJSON implements the Unmarshaler interface for Date.
// copied from https://github.com/goinvest/iexcloud, and adjusted
func (d *Date) UnmarshalJSON(data []byte) error {
	if len(data) <= 0 || string(data) == "null" || string(data) == "{}" {
		return nil
	}

	var aux string

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return errors.Wrap(err, "error unmarshaling date to string")
	}

	t, err := time.Parse("2006-01-02", aux)
	if err != nil {
		return errors.Wrap(err, "error converting string to date")
	}

	*d = Date(t)
	return nil
}

// MarshalJSON implements the Marshaler interface for Date.
// copied from https://github.com/goinvest/iexcloud
func (d *Date) MarshalJSON() ([]byte, error) {
	t := time.Time(*d)
	return json.Marshal(t.Format("2006-01-02"))
}

// EpochTime refers to unix timestamps used for some fields in the API
// copied from https://github.com/goinvest/iexcloud
type EpochTime time.Time

// MarshalJSON implements the Marshaler interface for EpochTime.
// copied from https://github.com/goinvest/iexcloud
func (e EpochTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint(time.Time(e).Unix())), nil
}

// UnmarshalJSON implements the Unmarshaler interface for EpochTime.
// copied from https://github.com/goinvest/iexcloud, and adjusted
func (e *EpochTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) <= 0 || string(data) == "null" || string(data) == "{}" {
		return nil
	}

	ts, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	// Per docs: If the value is -1, IEX has not quoted the symbol in the trading day.
	if ts == -1 {
		return
	}

	*e = EpochTime(time.Unix(int64(ts)/1000, 0))
	return
}

// String implements the Stringer interface for EpochTime.
// copied from https://github.com/goinvest/iexcloud
func (e EpochTime) String() string { return time.Time(e).String() }
