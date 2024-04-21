package fuel

import "context"

func (c *Client) GetLatestBlockHeight(ctx context.Context) (U32, error) {
	query := `
{
  chain {
    latestBlock {
      header {
		height
      }
    }
  }
}
`
	type chainInfoResult struct {
		Chain ChainInfo `json:"chain"`
	}
	chainInfo, err := ExecuteQuery[chainInfoResult](ctx, c, query)
	if err != nil {
		return 0, err
	}
	return chainInfo.Chain.LatestBlock.Header.Height, nil
}

func (c *Client) GetLatestBlockHeader(ctx context.Context) (Header, error) {
	query := `
{
  chain {
    latestBlock {
      header {
		id
		daHeight
		transactionsCount
		messageReceiptCount
		transactionsRoot
		messageReceiptRoot
		height
		prevRoot
		time
		applicationHash
      }
    }
  }
}
`
	type resultType struct {
		Chain ChainInfo `json:"chain"`
	}
	result, err := ExecuteQuery[resultType](ctx, c, query)
	if err != nil {
		return Header{}, err
	}
	return result.Chain.LatestBlock.Header, nil
}
