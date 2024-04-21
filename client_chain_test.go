package fuel

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getLastestBlockHeight(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	h, err := cli.GetLatestBlockHeight(context.Background())
	assert.NoError(t, err)
	fmt.Printf("%d\n", h)
}

func Test_getLastestBlockHeader(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	h, err := cli.GetLatestBlockHeader(context.Background())
	assert.NoError(t, err)
	fmt.Printf("%#v\n", h)
}
