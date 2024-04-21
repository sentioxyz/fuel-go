package fuel

import (
	"context"
	"fmt"
	"github.com/sentioxyz/fuel-go/util/query"
)

type GetTransactionOption struct {
	WithReceipts bool
	WithStatus   bool
}

func (c *Client) GetTransaction(
	ctx context.Context,
	param QueryTransactionParams,
	opt GetTransactionOption,
) (*Transaction, error) {
	ignoreCheckers := []query.IgnoreChecker{
		query.IgnoreField(Block{}, "Transactions"),
		query.IgnoreField(Contract{}, "Bytecode"),
		query.IgnoreField(SuccessStatus{}, "Receipts"),
		query.IgnoreField(FailureStatus{}, "Receipts"),
	}
	if !opt.WithReceipts {
		ignoreCheckers = append(ignoreCheckers, query.IgnoreField(Transaction{}, "Receipts"))
	}
	if !opt.WithStatus {
		ignoreCheckers = append(ignoreCheckers, query.IgnoreField(Transaction{}, "Status"))
	}
	q := fmt.Sprintf("{ transaction(%s) { %s } }",
		query.Simple.GenParam(param),
		query.Simple.GenObjectQuery(Transaction{}, query.MergeIgnores(ignoreCheckers...)),
	)
	//fmt.Printf("!!! query: %s\n", q)
	type resultType struct {
		Transaction *Transaction `json:"transaction"`
	}
	result, err := ExecuteQuery[resultType](ctx, c, q)
	if err != nil {
		return nil, err
	}
	return result.Transaction, nil
}
