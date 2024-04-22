package fuel

import (
	"context"
	"fmt"
	"github.com/sentioxyz/fuel-go/query"
	"github.com/sentioxyz/fuel-go/types"
)

type GetBlockOption struct {
	OnlyIdAndHeight      bool
	WithTransactions     bool
	OnlyTransactionID    bool
	WithContractBytecode bool
}

func (o GetBlockOption) BuildIgnoreChecker() query.IgnoreChecker {
	var checkers []query.IgnoreChecker
	if o.OnlyIdAndHeight {
		checkers = []query.IgnoreChecker{
			query.IgnoreOtherFields(types.Block{}, "Id", "Header"),
			query.IgnoreOtherFields(types.Header{}, "Id", "Height"),
		}
	} else if o.WithTransactions {
		if o.OnlyTransactionID {
			checkers = []query.IgnoreChecker{query.IgnoreOtherFields(types.Transaction{}, "Id")}
		} else {
			// Otherwise it will create circular dependencies
			checkers = []query.IgnoreChecker{query.IgnoreObjects(types.SuccessStatus{}, types.FailureStatus{})}
		}
	} else {
		checkers = []query.IgnoreChecker{query.IgnoreObjects(types.Transaction{})}
	}
	if !o.WithContractBytecode {
		checkers = append(checkers, query.IgnoreField(types.Contract{}, "Bytecode"))
	}
	return query.MergeIgnores(checkers...)
}

func (c *Client) GetBlock(ctx context.Context, param types.QueryBlockParams, opt GetBlockOption) (*types.Block, error) {
	q := fmt.Sprintf("{ block(%s) { %s} }",
		query.Simple.GenParam(param),
		query.Simple.GenObjectQuery(types.Block{}, opt.BuildIgnoreChecker()),
	)
	type resultType struct {
		Block *types.Block `json:"block"`
	}
	result, err := ExecuteQuery[resultType](ctx, c, q)
	if err != nil {
		return nil, err
	}
	return result.Block, nil
}

func (c *Client) GetBlockHeader(ctx context.Context, param types.QueryBlockParams) (*types.Header, error) {
	block, err := c.GetBlock(ctx, param, GetBlockOption{})
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, err
	}
	return &block.Header, nil
}
