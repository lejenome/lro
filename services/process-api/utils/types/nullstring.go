package types
import (
	"database/sql"
	"encoding/json"
)

type NullTime struct {
	sql.NullTime
}

func (t *NullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time)
}
