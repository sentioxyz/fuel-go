package fuel

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetChain(t *testing.T) {
	cli := NewClientWithLogger(testnetEndpoint, SimpleStdoutLogger{})
	_, err := cli.GetChain(context.Background(), GetChainOption{Simple: true})
	assert.NoError(t, err)
}

func Test_GetLatestBlockHeight(t *testing.T) {
	cli := NewClientWithLogger(testnetEndpoint, SimpleStdoutLogger{})
	h, err := cli.GetLatestBlockHeight(context.Background())
	assert.NoError(t, err)
	fmt.Printf("%d\n", h)
}

func Test_GetLatestBlockHeader(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	h, err := cli.GetLatestBlockHeader(context.Background())
	assert.NoError(t, err)
	fmt.Printf("%#v\n", h)
}
