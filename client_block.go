package fuel

import (
	"context"
	"fmt"
	"github.com/sentioxyz/fuel-go/query"
	"github.com/sentioxyz/fuel-go/types"
	"strings"
)

type GetBlockOption struct {
	WithHeader              bool
	WithConsensus           bool
	WithTransactions        bool
	WithTransactionDetail   bool
	WithTransactionReceipts bool
	WithContractBytecode    bool
	WithContractSalt        bool
}

func (o GetBlockOption) BuildIgnoreChecker() query.IgnoreChecker {
	var checkers []query.IgnoreChecker
	if !o.WithHeader {
		checkers = append(checkers, query.IgnoreField(types.Block{}, "Header"))
	}
	if !o.WithConsensus {
		checkers = append(checkers, query.IgnoreField(types.Block{}, "Consensus"))
	}
	if !o.WithTransactions {
		checkers = append(checkers, query.IgnoreField(types.Block{}, "Transactions"))
		checkers = append(checkers, query.IgnoreField(types.Block{}, "TransactionIds"))
	} else if !o.WithTransactionDetail {
		checkers = append(checkers, query.IgnoreOtherFields(types.Transaction{}, "Id"))
	} else {
		// Otherwise it will create circular dependencies
		checkers = append(checkers, query.IgnoreField(types.SuccessStatus{}, "Block"))
		checkers = append(checkers, query.IgnoreField(types.FailureStatus{}, "Block"))
		checkers = append(checkers, query.IgnoreField(types.SuccessStatus{}, "Transaction"))
		checkers = append(checkers, query.IgnoreField(types.FailureStatus{}, "Transaction"))
		if !o.WithTransactionReceipts {
			checkers = append(checkers, query.IgnoreField(types.SuccessStatus{}, "Receipts"))
			checkers = append(checkers, query.IgnoreField(types.FailureStatus{}, "Receipts"))
		}
	}
	if !o.WithContractBytecode {
		checkers = append(checkers, query.IgnoreField(types.Contract{}, "Bytecode"))
	}
	if !o.WithContractSalt {
		checkers = append(checkers, query.IgnoreField(types.Contract{}, "Salt"))
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

func (c *Client) GetBlocks(
	ctx context.Context,
	params []types.QueryBlockParams,
	opt GetBlockOption,
) ([]*types.Block, error) {
	bqs := make([]string, len(params))
	for i, param := range params {
		bqs[i] = fmt.Sprintf("b%d:block(%s) { %s}",
			i,
			query.Simple.GenParam(param),
			query.Simple.GenObjectQuery(types.Block{}, opt.BuildIgnoreChecker()),
		)
	}
	q := "{" + strings.Join(bqs, " ") + " }"
	type resultType map[string]*types.Block
	result, err := ExecuteQuery[resultType](ctx, c, q)
	if err != nil {
		return nil, err
	}
	blocks := make([]*types.Block, len(params))
	for i := range params {
		blocks[i] = result[fmt.Sprintf("b%d", i)]
	}
	return blocks, nil
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
