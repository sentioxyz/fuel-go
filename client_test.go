package fuel

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sentioxyz/fuel-go/query"
	"github.com/sentioxyz/fuel-go/types"
	"github.com/sentioxyz/fuel-go/util"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"testing"
	"time"
)

const testnetEndpoint = "https://testnet.fuel.network/v1/graphql"

func Test_ExecuteQuery(t *testing.T) {
	cli := NewClient(testnetEndpoint)

	type result struct {
		Block types.Block `json:"block"`
	}
	{
		q := `
{
  block(height: "1067005") {
    id
    header {
      height
      time
    }
  }
}`
		r, err := ExecuteQuery[result](context.Background(), cli, q)

		ti, _ := time.Parse(time.RFC3339, "2024-05-29T01:43:37Z")
		assert.NoError(t, err)
		assert.Equal(t, types.Block{
			Id: types.BlockId{Hash: common.HexToHash("0x4e02668366cbdc2ea9197fa3a84e57e723028de5fe4f574ccaa7b6b744ced495")},
			Header: types.Header{
				Height: 1067005,
				Time:   types.Tai64Timestamp{Time: ti},
			},
		}, r.Block)
	}

	{
		q := `
{
  block(height: "1067005") {
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
		"version id height header { version id daHeight consensusParametersVersion stateTransitionBytecodeVersion transactionsCount messageReceiptCount transactionsRoot messageOutboxRoot eventInboxRoot height prevRoot time applicationHash } consensus { __typename ... on Genesis { chainConfigHash coinsRoot contractsRoot messagesRoot transactionsRoot } ... on PoAConsensus { signature } } transactionIds ",
		query.Simple.GenObjectQuery(types.Block{}, query.IgnoreObjects(types.Transaction{})),
	)
	assert.Equal(t, `version
id
height
header {
  version
  id
  daHeight
  consensusParametersVersion
  stateTransitionBytecodeVersion
  transactionsCount
  messageReceiptCount
  transactionsRoot
  messageOutboxRoot
  eventInboxRoot
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
    transactionsRoot
  }
  ... on PoAConsensus {
    signature
  }
}
transactionIds
`,
		query.Beauty.GenObjectQuery(types.Block{}, query.IgnoreObjects(types.Transaction{})),
	)

	assert.Equal(t,
		"version id height header { version id daHeight consensusParametersVersion stateTransitionBytecodeVersion transactionsCount messageReceiptCount transactionsRoot messageOutboxRoot eventInboxRoot height prevRoot time applicationHash } consensus { __typename ... on Genesis { chainConfigHash coinsRoot contractsRoot messagesRoot transactionsRoot } ... on PoAConsensus { signature } } transactionIds transactions { id inputAssetIds inputContracts inputContract { utxoId balanceRoot stateRoot txPointer contractId } policies { tip witnessLimit maturity maxFee } scriptGasLimit maturity mintAmount mintAssetId mintGasPrice txPointer isScript isCreate isMint isUpgrade isUpload isBlob inputs { __typename ... on InputCoin { utxoId owner amount assetId txPointer witnessIndex predicateGasUsed predicate predicateData } ... on InputContract { utxoId balanceRoot stateRoot txPointer contractId } ... on InputMessage { sender recipient amount nonce witnessIndex predicateGasUsed data predicate predicateData } } outputs { __typename ... on CoinOutput { to amount assetId } ... on ContractOutput { inputIndex balanceRoot stateRoot } ... on ChangeOutput { to amount assetId } ... on VariableOutput { to amount assetId } ... on ContractCreated { contract stateRoot } } outputContract { inputIndex balanceRoot stateRoot } witnesses receiptsRoot status { __typename ... on SubmittedStatus { time } ... on SqueezedOutStatus { reason } } script scriptData bytecodeWitnessIndex blobId salt storageSlots bytecodeRoot subsectionIndex subsectionsNumber proofSet upgradePurpose { __typename ... on ConsensusParametersPurpose { witnessIndex checksum } ... on StateTransitionPurpose { root } } rawPayload } ",
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
			BlockHeight: 5,
			Block: types.Block{
				Version: "V1",
				Id:      types.BlockId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
				Height:  5,
				Header: types.Header{
					Version:             "V1",
					Id:                  types.BlockId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
					DaHeight:            2,
					TransactionsCount:   3,
					MessageReceiptCount: 4,
					Height:              5,
					Time:                types.Tai64Timestamp{Time: time.Date(2024, time.April, 15, 2, 44, 2, 0, time.UTC)},
				},
			},
			Transaction: types.Transaction{
				Id: types.TransactionId{Hash: common.HexToHash("0x9999")},
			},
		},
	}
	text, err := json.MarshalIndent(status, "", "  ")
	assert.NoError(t, err)
	assert.Equal(t, `{
  "__typename": "SuccessStatus",
  "transactionId": "0x0000000000000000000000000000000000000000000000000000000000000123",
  "blockHeight": "5",
  "block": {
    "version": "V1",
    "id": "0x0000000000000000000000000000000000000000000000000000000000000001",
    "height": "5",
    "header": {
      "version": "V1",
      "id": "0x0000000000000000000000000000000000000000000000000000000000000001",
      "daHeight": "2",
      "consensusParametersVersion": "0",
      "stateTransitionBytecodeVersion": "0",
      "transactionsCount": "3",
      "messageReceiptCount": "4",
      "transactionsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "messageOutboxRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "eventInboxRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "height": "5",
      "prevRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "time": "4611686020140536983",
      "applicationHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
    },
    "consensus": null,
    "transactionIds": null,
    "transactions": null
  },
  "transaction": {
    "id": "0x0000000000000000000000000000000000000000000000000000000000009999",
    "inputAssetIds": null,
    "inputContracts": null,
    "inputContract": null,
    "policies": null,
    "scriptGasLimit": null,
    "maturity": null,
    "mintAmount": null,
    "mintAssetId": null,
    "mintGasPrice": null,
    "txPointer": null,
    "isScript": false,
    "isCreate": false,
    "isMint": false,
    "isUpgrade": false,
    "isUpload": false,
    "isBlob": false,
    "inputs": null,
    "outputs": null,
    "outputContract": null,
    "witnesses": null,
    "receiptsRoot": null,
    "status": null,
    "script": null,
    "scriptData": null,
    "bytecodeWitnessIndex": null,
    "blobId": null,
    "salt": null,
    "storageSlots": null,
    "bytecodeRoot": null,
    "subsectionIndex": null,
    "subsectionsNumber": null,
    "proofSet": null,
    "upgradePurpose": null,
    "rawPayload": "0x"
  },
  "time": "4611685956291791114",
  "programState": null,
  "receipts": null,
  "totalGas": "0",
  "totalFee": "0"
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

func Test_union_unmarshalJSON(t *testing.T) {
	var con types.Consensus
	text, err := json.Marshal(con)
	assert.NoError(t, err)
	assert.Equal(t, `null`, string(text))

	assert.NoError(t, json.Unmarshal([]byte("null"), &con))
	assert.Equal(t, types.Consensus{}, con)

	assert.NoError(t, json.Unmarshal([]byte(`{"__typename":""}`), &con))
	assert.Equal(t, types.Consensus{}, con)
}

func Test_marshalStructpb(t *testing.T) {
	txn := &txn_1067005_0
	exp := structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
		"id": structpb.NewStringValue("0x0ec0390a47eb248d579c74861d747259c2a2a3f4c5c4cdccf049f0670b9a4485"),
		"inputAssetIds": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStringValue("0xf8f8b6283d7fa5b672b530cbb84fcccb4ff8dc40f8176ef4544ddb1f1952ad07"),
		}}),
		"inputContracts": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStringValue("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850"),
		}}),
		"policies": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
			"maxFee": structpb.NewStringValue("59243"),
		}}),
		"scriptGasLimit": structpb.NewStringValue("100000"),
		"maturity":       structpb.NewStringValue("0"),
		"isScript":       structpb.NewBoolValue(true),
		"isCreate":       structpb.NewBoolValue(false),
		"isMint":         structpb.NewBoolValue(false),
		"isUpgrade":      structpb.NewBoolValue(false),
		"isUpload":       structpb.NewBoolValue(false),
		"isBlob":         structpb.NewBoolValue(false),
		"inputs": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"__typename":  structpb.NewStringValue("InputContract"),
				"utxoId":      structpb.NewStringValue("0x4ed2af7ccf2e111c376d0fe486396f594b360d93d694875361d56f43705c3fdf0000"),
				"balanceRoot": structpb.NewStringValue("0x313d948ce814d576ffe3dfdf36758dbd51356d9d2a3adcba12001939a4442b17"),
				"stateRoot":   structpb.NewStringValue("0x6021474cda1220f0dcae5824b3ea73beb5752b474e648701c85197ea530591ed"),
				"txPointer":   structpb.NewStringValue("0008b29d0000"),
				"contractId":  structpb.NewStringValue("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850"),
			}}),
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"__typename":       structpb.NewStringValue("InputCoin"),
				"utxoId":           structpb.NewStringValue("0x1799f59d5f3ee48479e18feb0e5e705aa4d53f677ba21630ec85200381e621e90000"),
				"owner":            structpb.NewStringValue("0xd3fe20c8ff68a4054d8587ac170c40db7d1e200208a575780542bd9a7e3eec08"),
				"amount":           structpb.NewStringValue("2000000"),
				"assetId":          structpb.NewStringValue("0xf8f8b6283d7fa5b672b530cbb84fcccb4ff8dc40f8176ef4544ddb1f1952ad07"),
				"txPointer":        structpb.NewStringValue("001046f00000"),
				"witnessIndex":     structpb.NewStringValue("0"),
				"predicateGasUsed": structpb.NewStringValue("0"),
				"predicate":        structpb.NewStringValue("0x"),
				"predicateData":    structpb.NewStringValue("0x"),
			}}),
		}}),
		"outputs": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"__typename":  structpb.NewStringValue("ContractOutput"),
				"inputIndex":  structpb.NewStringValue("0"),
				"balanceRoot": structpb.NewStringValue("0x313d948ce814d576ffe3dfdf36758dbd51356d9d2a3adcba12001939a4442b17"),
				"stateRoot":   structpb.NewStringValue("0xeb827e501eee0dff74b423e5173bcd9e3147d00fb2da246c2bfe8abc98660d33"),
			}}),
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"__typename": structpb.NewStringValue("ChangeOutput"),
				"to":         structpb.NewStringValue("0xd3fe20c8ff68a4054d8587ac170c40db7d1e200208a575780542bd9a7e3eec08"),
				"amount":     structpb.NewStringValue("1941078"),
				"assetId":    structpb.NewStringValue("0xf8f8b6283d7fa5b672b530cbb84fcccb4ff8dc40f8176ef4544ddb1f1952ad07"),
			}}),
		}}),
		"witnesses": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStringValue("0xf34f6666876b62e6930f7ed447c6ffd86344237325cd4db19f2ce1a0be2cc76b995e73e04f039f7ad60f614edf6b85667e876652499270d111cf2acc7ca41b5b"),
		}}),
		"receiptsRoot": structpb.NewStringValue("0x74ab41b67d3e1ccb4bda92dfac21bfb448f80d0500b3cd56480aaee2fa37cbf0"),
		"status": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
			"__typename":    structpb.NewStringValue("SuccessStatus"),
			"blockHeight":   structpb.NewStringValue("1067005"),
			"transactionId": structpb.NewStringValue("0x0ec0390a47eb248d579c74861d747259c2a2a3f4c5c4cdccf049f0670b9a4485"),
			"block": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"version": structpb.NewStringValue(""),
				"id":      structpb.NewStringValue("0x0000000000000000000000000000000000000000000000000000000000000000"),
				"height":  structpb.NewStringValue("0"),
				"header": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"version":                        structpb.NewStringValue(""),
					"id":                             structpb.NewStringValue("0x0000000000000000000000000000000000000000000000000000000000000000"),
					"daHeight":                       structpb.NewStringValue("0"),
					"consensusParametersVersion":     structpb.NewStringValue("0"),
					"stateTransitionBytecodeVersion": structpb.NewStringValue("0"),
					"transactionsCount":              structpb.NewStringValue("0"),
					"messageReceiptCount":            structpb.NewStringValue("0"),
					"transactionsRoot":               structpb.NewStringValue("0x0000000000000000000000000000000000000000000000000000000000000000"),
					"messageOutboxRoot":              structpb.NewStringValue("0x0000000000000000000000000000000000000000000000000000000000000000"),
					"eventInboxRoot":                 structpb.NewStringValue("0x0000000000000000000000000000000000000000000000000000000000000000"),
					"height":                         structpb.NewStringValue("0"),
					"prevRoot":                       structpb.NewStringValue("0x0000000000000000000000000000000000000000000000000000000000000000"),
					"time":                           structpb.NewStringValue("4611685956291791114"),
					"applicationHash":                structpb.NewStringValue("0x0000000000000000000000000000000000000000000000000000000000000000"),
				}}),
				"transactionIds": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
				"transactions":   structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
			}}),
			"transaction": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"id":             structpb.NewStringValue("0x0000000000000000000000000000000000000000000000000000000000000000"),
				"inputAssetIds":  structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
				"inputContracts": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
				"inputs":         structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
				"outputs":        structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
				"proofSet":       structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
				"storageSlots":   structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
				"witnesses":      structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
				"isScript":       structpb.NewBoolValue(false),
				"isCreate":       structpb.NewBoolValue(false),
				"isMint":         structpb.NewBoolValue(false),
				"isUpgrade":      structpb.NewBoolValue(false),
				"isUpload":       structpb.NewBoolValue(false),
				"isBlob":         structpb.NewBoolValue(false),
				"rawPayload":     structpb.NewStringValue("0x"),
			}}),
			"receipts": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
				structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"pc":          structpb.NewStringValue("11712"),
					"is":          structpb.NewStringValue("11712"),
					"to":          structpb.NewStringValue("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850"),
					"amount":      structpb.NewStringValue("0"),
					"assetId":     structpb.NewStringValue("0x0000000000000000000000000000000000000000000000000000000000000000"),
					"gas":         structpb.NewStringValue("80154"),
					"param1":      structpb.NewStringValue("10480"),
					"param2":      structpb.NewStringValue("10497"),
					"receiptType": structpb.NewStringValue("CALL"),
				}}),
				structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"id":          structpb.NewStringValue("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850"),
					"pc":          structpb.NewStringValue("16640"),
					"is":          structpb.NewStringValue("11712"),
					"ptr":         structpb.NewStringValue("67107840"),
					"digest":      structpb.NewStringValue("0x979f4ef2bcab47538ec8c1c92b8ec5c58ef61b7aa86a708fb9823324479144d6"),
					"ra":          structpb.NewStringValue("0"),
					"rb":          structpb.NewStringValue("14631882454972838106"),
					"len":         structpb.NewStringValue("25"),
					"receiptType": structpb.NewStringValue("LOG_DATA"),
					"data":        structpb.NewStringValue("0x63616c6c696e6720696e6372656d656e74206d6574686f642e"),
				}}),
				structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"id":          structpb.NewStringValue("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850"),
					"pc":          structpb.NewStringValue("16896"),
					"is":          structpb.NewStringValue("11712"),
					"ptr":         structpb.NewStringValue("67106816"),
					"digest":      structpb.NewStringValue("0x17eb70034b5b71092521d184c5e7b069d47de657e51aef2be11a00c115036943"),
					"ra":          structpb.NewStringValue("0"),
					"rb":          structpb.NewStringValue("15520703124961489725"),
					"len":         structpb.NewStringValue("4"),
					"receiptType": structpb.NewStringValue("LOG_DATA"),
					"data":        structpb.NewStringValue("0x00000008"),
				}}),
				structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"id":          structpb.NewStringValue("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850"),
					"pc":          structpb.NewStringValue("17132"),
					"is":          structpb.NewStringValue("11712"),
					"ptr":         structpb.NewStringValue("67105792"),
					"digest":      structpb.NewStringValue("0xcd2662154e6d76b2b2b92e70c0cac3ccf534f9b74eb5b89819ec509083d00a50"),
					"ra":          structpb.NewStringValue("0"),
					"rb":          structpb.NewStringValue("1515152261580153489"),
					"len":         structpb.NewStringValue("8"),
					"receiptType": structpb.NewStringValue("LOG_DATA"),
					"data":        structpb.NewStringValue("0x0000000000000001"),
				}}),
				structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"id":          structpb.NewStringValue("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850"),
					"pc":          structpb.NewStringValue("17380"),
					"is":          structpb.NewStringValue("11712"),
					"ptr":         structpb.NewStringValue("67104768"),
					"digest":      structpb.NewStringValue("0x02d449a31fbb267c8f352e9968a79e3e5fc95c1bbeaa502fd6454ebde5a4bedc"),
					"ra":          structpb.NewStringValue("0"),
					"rb":          structpb.NewStringValue("8961848586872524460"),
					"len":         structpb.NewStringValue("32"),
					"receiptType": structpb.NewStringValue("LOG_DATA"),
					"data":        structpb.NewStringValue("0x1111111111111111111111111111111111111111111111111111111111111111"),
				}}),
				structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"id":          structpb.NewStringValue("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850"),
					"pc":          structpb.NewStringValue("18408"),
					"is":          structpb.NewStringValue("11712"),
					"ptr":         structpb.NewStringValue("67103232"),
					"digest":      structpb.NewStringValue("0x5bc67471c189d78c76461dcab6141a733bdab3799d1d69e0c419119c92e82b3d"),
					"len":         structpb.NewStringValue("8"),
					"receiptType": structpb.NewStringValue("RETURN_DATA"),
					"data":        structpb.NewStringValue("0x0000000000000012"),
				}}),
				structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"pc":          structpb.NewStringValue("10388"),
					"is":          structpb.NewStringValue("10368"),
					"val":         structpb.NewStringValue("1"),
					"receiptType": structpb.NewStringValue("RETURN"),
				}}),
				structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"receiptType": structpb.NewStringValue("SCRIPT_RESULT"),
					"result":      structpb.NewStringValue("0"),
					"gasUsed":     structpb.NewStringValue("70408"),
				}}),
			}}),
			"time": structpb.NewStringValue("4611686020144334958"),
			"programState": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"returnType": structpb.NewStringValue("RETURN"),
				"data":       structpb.NewStringValue("0x0000000000000001"),
			}}),
			"totalGas": structpb.NewStringValue("5420740"),
			"totalFee": structpb.NewStringValue("58922"),
		}}),
		"script":       structpb.NewStringValue("0x724028c0724428985d451000724828a02d41148a24040000"),
		"scriptData":   structpb.NewStringValue("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000db0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e85000000000000028f000000000000029010000000000000009696e6372656d656e740000000800000000000000011111111111111111111111111111111111111111111111111111111111111111"),
		"storageSlots": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
		"proofSet":     structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
		"rawPayload":   structpb.NewStringValue("0x000000000000000000000000000186a074ab41b67d3e1ccb4bda92dfac21bfb448f80d0500b3cd56480aaee2fa37cbf0000000000000001800000000000000950000000000000008000000000000000200000000000000020000000000000001724028c0724428985d451000724828a02d41148a2404000000000000000000000000000000000000000000000000000000000000000000000000000000000000db0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e85000000000000028f000000000000029010000000000000009696e6372656d656e740000000800000000000000011111111111111111111111111111111111111111111111111111111111111111000000000000000000e76b00000000000000014ed2af7ccf2e111c376d0fe486396f594b360d93d694875361d56f43705c3fdf0000000000000000313d948ce814d576ffe3dfdf36758dbd51356d9d2a3adcba12001939a4442b176021474cda1220f0dcae5824b3ea73beb5752b474e648701c85197ea530591ed000000000008b29d0000000000000000db0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e85000000000000000001799f59d5f3ee48479e18feb0e5e705aa4d53f677ba21630ec85200381e621e90000000000000000d3fe20c8ff68a4054d8587ac170c40db7d1e200208a575780542bd9a7e3eec0800000000001e8480f8f8b6283d7fa5b672b530cbb84fcccb4ff8dc40f8176ef4544ddb1f1952ad0700000000001046f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000313d948ce814d576ffe3dfdf36758dbd51356d9d2a3adcba12001939a4442b17eb827e501eee0dff74b423e5173bcd9e3147d00fb2da246c2bfe8abc98660d330000000000000002d3fe20c8ff68a4054d8587ac170c40db7d1e200208a575780542bd9a7e3eec0800000000001d9e56f8f8b6283d7fa5b672b530cbb84fcccb4ff8dc40f8176ef4544ddb1f1952ad070000000000000040f34f6666876b62e6930f7ed447c6ffd86344237325cd4db19f2ce1a0be2cc76b995e73e04f039f7ad60f614edf6b85667e876652499270d111cf2acc7ca41b5b"),
	}})
	//tt1, _ := json.MarshalIndent(exp, "", "  ")
	//fmt.Printf("!!! ===1 %s\n", string(tt1))
	//tt2, _ := json.MarshalIndent(txn.MarshalStructpb(), "", "  ")
	//fmt.Printf("!!! ===2 %s\n", string(tt2))
	//assert.Equal(t, string(tt1), string(tt2))
	assert.Equal(t, exp, txn.MarshalStructpb())
}
