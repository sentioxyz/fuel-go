package types

import (
	"encoding/json"
	"fmt"
	"github.com/cactus/tai64"
	"math/big"
	"strings"
)

func (n *Tai64Timestamp) UnmarshalJSON(raw []byte) error {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return err
	}
	var b big.Int
	_, ok := b.SetString(s, 10)
	if !ok {
		return fmt.Errorf("invalid number %q", s)
	}
	t, err := tai64.Parse(b.Text(16))
	if err != nil {
		return err
	}
	*n = Tai64Timestamp{
		Time: t,
	}
	return nil
}

func (n Tai64Timestamp) MarshalJSON() ([]byte, error) {
	var b big.Int
	b.SetString(strings.TrimLeft(tai64.Format(n.Time), "@"), 16)
	return json.Marshal(b.Text(10))
}
