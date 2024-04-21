package fuel

import (
	"context"
	"fmt"
	"github.com/sentioxyz/fuel-go/types"
	"github.com/sentioxyz/fuel-go/util/query"
)

type GetTransactionOption struct {
	WithReceipts bool
	WithStatus   bool
}

func (c *Client) GetTransaction(
	ctx context.Context,
	param types.QueryTransactionParams,
	opt GetTransactionOption,
) (*types.Transaction, error) {
	ignoreCheckers := []query.IgnoreChecker{
		query.IgnoreField(types.Block{}, "Transactions"),
		query.IgnoreField(types.Contract{}, "Bytecode"),
		query.IgnoreField(types.SuccessStatus{}, "Receipts"),
		query.IgnoreField(types.FailureStatus{}, "Receipts"),
	}
	if !opt.WithReceipts {
		ignoreCheckers = append(ignoreCheckers, query.IgnoreField(types.Transaction{}, "Receipts"))
	}
	if !opt.WithStatus {
		ignoreCheckers = append(ignoreCheckers, query.IgnoreField(types.Transaction{}, "Status"))
	}
	q := fmt.Sprintf("{ transaction(%s) { %s } }",
		query.Simple.GenParam(param),
		query.Simple.GenObjectQuery(types.Transaction{}, query.MergeIgnores(ignoreCheckers...)),
	)
	type resultType struct {
		Transaction *types.Transaction `json:"transaction"`
	}
	result, err := ExecuteQuery[resultType](ctx, c, q)
	if err != nil {
		return nil, err
	}
	return result.Transaction, nil
}
