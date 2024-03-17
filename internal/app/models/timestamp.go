package m

import (
	"strings"
	"time"
)

type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
	var b strings.Builder

	b.WriteString("\"")
	b.WriteString(time.Time(t).Format(time.RFC3339))
	b.WriteString("\"")

	return []byte(b.String()), nil
}

type TimeFields struct {
	CreatedAt Timestamp `json:"created_at,omitempty"`
	UpdatedAt Timestamp `json:"updated_at,omitempty"`
}
