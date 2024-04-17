package fuel

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_executeQuery(t *testing.T) {
	cli := NewClient("https://beta-5.fuel.network/graphql")

	type result struct {
		Block Block `json:"block"`
	}
	{
		query := `
{
  block(height: "9758550") {
    id
    header {
      height
      time
    }
  }
}`
		r, err := ExecuteQuery[result](context.Background(), cli, query)

		ti, _ := time.Parse(time.RFC3339, "2024-04-15T02:44:02Z")
		assert.NoError(t, err)
		assert.Equal(t, Block{
			Id: BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
			Header: Header{
				Height: 9758550,
				Time:   Tai64Timestamp{Time: ti},
			},
		}, r.Block)
	}

	{
		query := `
{
  block(height: "9758550") {
    id
    header {
      height
      tim
    }
  }
}`
		_, err := ExecuteQuery[result](context.Background(), cli, query)

		assert.EqualError(t, err, "execute query failed: (line:7,column:7): Unknown field \"tim\" on type \"Header\". Did you mean \"time\"?")
	}
}
