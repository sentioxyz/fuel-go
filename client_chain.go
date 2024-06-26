package fuel

import (
	"context"
	"fmt"
	"github.com/sentioxyz/fuel-go/query"
	"github.com/sentioxyz/fuel-go/types"
)

type GetChainOption struct {
	Simple bool
	GetBlockOption
}

func (o GetChainOption) BuildIgnoreChecker() query.IgnoreChecker {
	if o.Simple {
		return query.MergeIgnores(
			query.IgnoreOtherFields(types.ChainInfo{}, "Name", "LatestBlock"),
			o.GetBlockOption.BuildIgnoreChecker(),
		)
	}
	return o.GetBlockOption.BuildIgnoreChecker()
}

func (c *Client) GetChain(ctx context.Context, opt GetChainOption) (types.ChainInfo, error) {
	q := fmt.Sprintf("{ chain { %s} }",
		query.Simple.GenObjectQuery(types.ChainInfo{}, opt.BuildIgnoreChecker()),
	)
	type resultType struct {
		Chain types.ChainInfo `json:"chain"`
	}
	result, err := ExecuteQuery[resultType](ctx, c, q)
	if err != nil {
		return types.ChainInfo{}, err
	}
	return result.Chain, nil
}

func (c *Client) GetLatestBlockHeight(ctx context.Context) (types.U32, error) {
	info, err := c.GetChain(ctx, GetChainOption{
		Simple:         true,
		GetBlockOption: GetBlockOption{},
	})
	return info.LatestBlock.Height, err
}

func (c *Client) GetLatestBlockHeader(ctx context.Context) (types.Header, error) {
	info, err := c.GetChain(ctx, GetChainOption{
		Simple:         true,
		GetBlockOption: GetBlockOption{WithHeader: true},
	})
	return info.LatestBlock.Header, err
}
