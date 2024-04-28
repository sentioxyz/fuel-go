package fuel

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sentioxyz/fuel-go/query"
	"github.com/sentioxyz/fuel-go/types"
	"github.com/sentioxyz/fuel-go/util"
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

func Test_Union_marshalJSON(t *testing.T) {
	status := types.TransactionStatus{
		TypeName_: "SuccessStatus",
		SuccessStatus: &types.SuccessStatus{
			TransactionId: types.TransactionId{
				Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000123"),
			},
			Block: types.Block{
				Id: types.BlockId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
				Header: types.Header{
					Id:                  types.BlockId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
					DaHeight:            2,
					TransactionsCount:   3,
					MessageReceiptCount: 4,
					Height:              5,
					Time:                types.Tai64Timestamp{Time: time.Date(2024, time.April, 15, 2, 44, 2, 0, time.UTC)},
				},
			},
		},
	}
	text, err := json.MarshalIndent(status, "", "  ")
	assert.NoError(t, err)
	assert.Equal(t, `{
  "__typename": "SuccessStatus",
  "transactionId": "0x0000000000000000000000000000000000000000000000000000000000000123",
  "block": {
    "id": "0x0000000000000000000000000000000000000000000000000000000000000001",
    "header": {
      "id": "0x0000000000000000000000000000000000000000000000000000000000000001",
      "daHeight": "2",
      "transactionsCount": "3",
      "messageReceiptCount": "4",
      "transactionsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "messageReceiptRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "height": "5",
      "prevRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "time": "4611686020140536983",
      "applicationHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
    },
    "consensus": {
      "__typename": ""
    },
    "transactions": null
  },
  "time": "4611685956291791114",
  "programState": null,
  "receipts": null
}`, string(text))

	con := types.Consensus{
		TypeName_: "PoAConsensus",
		PoAConsensus: &types.PoAConsensus{
			Signature: types.Signature{Bytes: common.FromHex("0x724028a0724428785d451000724828802d41148a24040000")},
		},
		Genesis: &types.Genesis{
			ChainConfigHash: types.Bytes32{Hash: common.HexToHash("0x1100000000000000000000000000000000000000000000000000000000000011")},
		},
	}
	text, err = json.Marshal(con)
	assert.NoError(t, err)
	assert.Equal(t, `{"__typename":"PoAConsensus","signature":"0x724028a0724428785d451000724828802d41148a24040000"}`, string(text))
}
