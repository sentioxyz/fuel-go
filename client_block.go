package fuel

import (
	"context"
	"fmt"
	"github.com/sentioxyz/fuel-go/util/query"
)

type GetBlockOption struct {
	WithTransactions  bool
	OnlyTransactionID bool
}

func (c *Client) GetBlock(ctx context.Context, param QueryBlockParams, opt GetBlockOption) (*Block, error) {
	var ignoreChecker query.IgnoreChecker
	if opt.WithTransactions {
		if opt.OnlyTransactionID {
			ignoreChecker = query.IgnoreOtherField(Transaction{}, "Id")
		} else {
			// Otherwise it will create circular dependencies
			ignoreChecker = query.IgnoreObjects(SuccessStatus{}, FailureStatus{})
		}
	} else {
		ignoreChecker = query.IgnoreObjects(Transaction{})
	}
	// Contract.Bytecode is not needed
	ignoreChecker = query.MergeIgnores(ignoreChecker, query.IgnoreField(Contract{}, "Bytecode"))
	q := fmt.Sprintf("{ block(%s) { %s } }",
		query.Simple.GenParam(param),
		query.Simple.GenObjectQuery(Block{}, ignoreChecker),
	)
	//fmt.Printf("!!! query: %s\n", q)
	type resultType struct {
		Block *Block `json:"block"`
	}
	result, err := ExecuteQuery[resultType](ctx, c, q)
	if err != nil {
		return nil, err
	}
	return result.Block, nil
}
