package fuel

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sentioxyz/fuel-go/types"
	"github.com/sentioxyz/fuel-go/util"
	"github.com/sentioxyz/fuel-go/util/query"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const testnetEndpoint = "https://beta-5.fuel.network/graphql"

func Test_ExecuteQuery(t *testing.T) {
	cli := NewClient(testnetEndpoint)

	type result struct {
		Block types.Block `json:"block"`
	}
	{
		q := `
{
  block(height: "9758550") {
    id
    header {
      height
      time
    }
  }
}`
		r, err := ExecuteQuery[result](context.Background(), cli, q)

		ti, _ := time.Parse(time.RFC3339, "2024-04-15T02:44:02Z")
		assert.NoError(t, err)
		assert.Equal(t, types.Block{
			Id: types.BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
			Header: types.Header{
				Height: 9758550,
				Time:   types.Tai64Timestamp{Time: ti},
			},
		}, r.Block)
	}

	{
		q := `
{
  block(height: "9758550") {
    id
    header {
      height
      tim
    }
  }
}`
		_, err := ExecuteQuery[result](context.Background(), cli, q)

		assert.EqualError(t, err, "execute query failed: (line:7,column:7): Unknown field \"tim\" on type \"Header\". Did you mean \"time\"?")
	}
}

func Test_GenObjectQuery(t *testing.T) {
	assert.Equal(t,
		"id header { id daHeight transactionsCount messageReceiptCount transactionsRoot messageReceiptRoot height prevRoot time applicationHash } consensus { __typename ... on Genesis { chainConfigHash coinsRoot contractsRoot messagesRoot } ... on PoAConsensus { signature } } ",
		query.Simple.GenObjectQuery(types.Block{}, query.IgnoreObjects(types.Transaction{})),
	)
	assert.Equal(t, `id
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
consensus {
  __typename
  ... on Genesis {
    chainConfigHash
    coinsRoot
    contractsRoot
    messagesRoot
  }
  ... on PoAConsensus {
    signature
  }
}
`,
		query.Beauty.GenObjectQuery(types.Block{}, query.IgnoreObjects(types.Transaction{})),
	)

	assert.Equal(t,
		"id header { id daHeight transactionsCount messageReceiptCount transactionsRoot messageReceiptRoot height prevRoot time applicationHash } consensus { __typename ... on Genesis { chainConfigHash coinsRoot contractsRoot messagesRoot } ... on PoAConsensus { signature } } transactions { id inputAssetIds inputContracts { id bytecode salt } inputContract { utxoId balanceRoot stateRoot txPointer contract { id bytecode salt } } policies { gasPrice witnessLimit maturity maxFee } gasPrice scriptGasLimit maturity mintAmount mintAssetId txPointer isScript isCreate isMint inputs { __typename ... on InputCoin { utxoId owner amount assetId txPointer witnessIndex maturity predicateGasUsed predicate predicateData } ... on InputContract { utxoId balanceRoot stateRoot txPointer contract { id bytecode salt } } ... on InputMessage { sender recipient amount nonce witnessIndex predicateGasUsed data predicate predicateData } } outputs { __typename ... on CoinOutput { to amount assetId } ... on ContractOutput { inputIndex balanceRoot stateRoot } ... on ChangeOutput { to amount assetId } ... on VariableOutput { to amount assetId } ... on ContractCreated { contract { id bytecode salt } stateRoot } } outputContract { inputIndex balanceRoot stateRoot } witnesses receiptsRoot status { __typename ... on SubmittedStatus { time } ... on SqueezedOutStatus { reason } } receipts { contract { id bytecode salt } pc is to { id bytecode salt } toAddress amount assetId gas param1 param2 val ptr digest reason ra rb rc rd len receiptType result gasUsed data sender recipient nonce contractId subId } script scriptData bytecodeWitnessIndex bytecodeLength salt storageSlots rawPayload } ",
		query.Simple.GenObjectQuery(types.Block{}, query.IgnoreObjects(types.SuccessStatus{}, types.FailureStatus{})),
	)
}

func Test_GenParam(t *testing.T) {
	assert.Equal(t,
		`id: "0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c" height: "1234" `,
		query.Simple.GenParam(types.QueryBlockParams{
			Id:     &types.BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
			Height: util.GetPointer[types.U32](1234),
		}),
	)
	assert.Equal(t,
		"id: \"0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c\"\nheight: \"1234\"\n",
		query.Beauty.GenParam(types.QueryBlockParams{
			Id:     &types.BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
			Height: util.GetPointer[types.U32](1234),
		}),
	)
	assert.Equal(t,
		"height: \"1234\"\n",
		query.Beauty.GenParam(types.QueryBlockParams{
			Height: util.GetPointer[types.U32](1234),
		}),
	)
	assert.Equal(t,
		"id: \"0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c\"\n",
		query.Beauty.GenParam(types.QueryBlockParams{
			Id: &types.BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
		}),
	)
}
