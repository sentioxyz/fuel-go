package fuel

import (
	"context"
	"fmt"
	"github.com/sentioxyz/fuel-go/query"
	"github.com/sentioxyz/fuel-go/types"
)

type GetTransactionOption struct {
	WithReceipts         bool
	WithStatus           bool
	WithContractBytecode bool
	WithContractSalt     bool
}

func (o GetTransactionOption) BuildIgnoreChecker() query.IgnoreChecker {
	ignoreCheckers := []query.IgnoreChecker{
		query.IgnoreField(types.Block{}, "Transactions"),
	}
	if !o.WithContractBytecode {
		ignoreCheckers = append(ignoreCheckers, query.IgnoreField(types.Contract{}, "Bytecode"))
	}
	if !o.WithContractSalt {
		ignoreCheckers = append(ignoreCheckers, query.IgnoreField(types.Contract{}, "Salt"))
	}
	if !o.WithReceipts {
		ignoreCheckers = append(ignoreCheckers, query.IgnoreObjects(types.Receipt{}))
	}
	if !o.WithStatus {
		ignoreCheckers = append(ignoreCheckers, query.IgnoreField(types.Transaction{}, "Status"))
	}
	return query.MergeIgnores(ignoreCheckers...)
}

func (c *Client) GetTransaction(
	ctx context.Context,
	param types.QueryTransactionParams,
	opt GetTransactionOption,
) (*types.Transaction, error) {
	q := fmt.Sprintf("{ transaction(%s) { %s} }",
		query.Simple.GenParam(param),
		query.Simple.GenObjectQuery(types.Transaction{}, opt.BuildIgnoreChecker()),
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
