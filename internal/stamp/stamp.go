package stamp

import (
	"strings"
	"time"
)

type Stamp struct {
	time.Time
}

func Parse(s string) Stamp {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return Stamp{}
	}

	return Stamp{t}
}

func (t Stamp) MarshalJSON() ([]byte, error) {
	var b strings.Builder

	b.WriteString("\"")
	b.WriteString(t.Format(time.RFC3339))
	b.WriteString("\"")

	return []byte(b.String()), nil
}
