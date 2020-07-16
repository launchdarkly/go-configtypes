package configtypes

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

// OptDuration represents an optional time.Duration parameter.
//
// When setting this value from a string representation, the following formats are allowed, where 9
// represents any number of digits: "9ms" (milliseconds), "9s" (seconds), "9m" (minutes), "9h" (hours),
// ":99" (seconds), "9:99" (minutes:seconds), "9:99:99" (hours:minutes:seconds).
//
// Converting to a string uses whichever format most compactly represents the value.
//
// When converting to or from JSON, an empty value is null, and all other values are strings.
//
// See the package documentation for the general contract for methods that have no specific documentation
// here.
type OptDuration struct {
	hasValue bool
	value    time.Duration
}

func NewOptDuration(value time.Duration) OptDuration {
	return OptDuration{hasValue: true, value: value}
}

func NewOptDurationFromString(s string) (OptDuration, error) {
	if s == "" {
		return OptDuration{}, nil
	}
	var n, hh, mm, ss int
	// note that the newlines in these format strings mean there should be no other characters after the format
	if count, err := fmt.Sscanf(s, "%dms\n", &n); err == nil && count == 1 {
		return NewOptDuration(time.Duration(n) * time.Millisecond), nil
	}
	if count, err := fmt.Sscanf(s, "%ds\n", &n); err == nil && count == 1 {
		return NewOptDuration(time.Duration(n) * time.Second), nil
	}
	if count, err := fmt.Sscanf(s, "%dm\n", &n); err == nil && count == 1 {
		return NewOptDuration(time.Duration(n) * time.Minute), nil
	}
	if count, err := fmt.Sscanf(s, "%dh\n", &n); err == nil && count == 1 {
		return NewOptDuration(time.Duration(n) * time.Hour), nil
	}
	if count, err := fmt.Sscanf(s, ":%d\n", &ss); err == nil && count == 1 {
		return NewOptDuration(time.Duration(ss) * time.Second), nil
	}
	if count, err := fmt.Sscanf(s, "%d:%d\n", &mm, &ss); err == nil && count == 2 {
		secs := mm*60 + ss
		return NewOptDuration(time.Duration(secs) * time.Second), nil
	}
	if count, err := fmt.Sscanf(s, "%d:%d:%d\n", &hh, &mm, &ss); err == nil && count == 3 {
		secs := (hh*60+mm)*60 + ss
		return NewOptDuration(time.Duration(secs) * time.Second), nil
	}
	return OptDuration{}, errDurationFormat()
}

func (o OptDuration) IsDefined() bool {
	return o.hasValue
}

func (o OptDuration) GetOrElse(orElseValue time.Duration) time.Duration {
	if !o.hasValue {
		return orElseValue
	}
	return o.value
}

func (o OptDuration) String() string {
	if !o.hasValue {
		return ""
	}
	d := o.value
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second
	d -= seconds * time.Second
	millis := d / time.Millisecond
	if hours > 0 && o.value == hours*time.Hour {
		return fmt.Sprintf("%dh", hours)
	}
	if minutes > 0 && o.value == minutes*time.Minute {
		return fmt.Sprintf("%dm", minutes)
	}
	if seconds > 0 && o.value == seconds*time.Second {
		return fmt.Sprintf("%ds", seconds)
	}
	if millis == 0 {
		if hours == 0 {
			return fmt.Sprintf("%d:%2d", minutes, seconds)
		}
		return fmt.Sprintf("%d:%2d:%2d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%dms", d/time.Millisecond)
}

func (o OptDuration) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o *OptDuration) UnmarshalText(data []byte) error {
	opt, err := NewOptDurationFromString(string(data))
	if err == nil {
		*o = opt
	}
	return err
}

func (o OptDuration) MarshalJSON() ([]byte, error) {
	if o.hasValue {
		return json.Marshal(o.String())
	}
	return json.Marshal(nil)
}

func (o *OptDuration) UnmarshalJSON(data []byte) error {
	var v ldvalue.Value
	var err error
	if err = v.UnmarshalJSON(data); err != nil {
		return err
	}
	switch {
	case v.IsNull():
		*o = OptDuration{}
		return nil
	case v.IsString():
		*o, err = NewOptDurationFromString(v.StringValue())
		return err
	default:
		return errDurationFormat()
	}
}
