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
    "consensus": null,
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
	txn := &types.Transaction{
		Id:            types.TransactionId{Hash: common.HexToHash("0x9b7a9353faacd4ce91c47707d66c81ec7e4d547905168a592312a94a5c67b69f")},
		InputAssetIds: []types.AssetId{{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")}},
		InputContracts: []types.Contract{{
			Id: types.ContractId{Hash: common.HexToHash("0xd2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f")},
		}},
		Policies: &types.Policies{
			GasPrice: util.GetPointer[types.U64](1),
		},
		GasPrice:       util.GetPointer[types.U64](1),
		ScriptGasLimit: util.GetPointer[types.U64](800000),
		Maturity:       util.GetPointer[types.U32](0),
		IsScript:       true,
		IsCreate:       false,
		IsMint:         false,
		Inputs: []types.Input{{
			TypeName_: "InputContract",
			InputContract: &types.InputContract{
				UtxoId:      types.UtxoId{Bytes: common.FromHex("0x16d77c9e9f146e5523a14fd428aa57a839ab775245ac9d5b662ee6f5e99fed2300")},
				BalanceRoot: types.Bytes32{Hash: common.HexToHash("0x5cc28e489493724ead2b4b771c276109f3995bf527b1c038b7cc14dbfe92758c")},
				StateRoot:   types.Bytes32{Hash: common.HexToHash("0x8f36f4ef87d3260fcbbb8b7d047bae772b12265d9c45bb11814855d57fdacee3")},
				TxPointer:   "0094e7550005",
				Contract: types.Contract{
					Id: types.ContractId{Hash: common.HexToHash("0xd2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f")},
				},
			},
		}, {
			TypeName_: "InputCoin",
			InputCoin: &types.InputCoin{
				UtxoId:           types.UtxoId{Bytes: common.FromHex("0x88dd1c739b7539af1d82ee96d26f528ab7e0ea2485f9c0febb4c2bdb884c19f101")},
				Owner:            types.Address{Hash: common.HexToHash("0xe173edec3aad8af6d0735165fc013527e93316af30d035dc465dbb6f37cb6972")},
				Amount:           1974461,
				AssetId:          types.AssetId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
				TxPointer:        "0094e74d0001",
				WitnessIndex:     0,
				Maturity:         0,
				PredicateGasUsed: 0,
				Predicate:        types.HexString{Bytes: common.FromHex("0x")},
				PredicateData:    types.HexString{Bytes: common.FromHex("0x")},
			},
		}},
		Outputs: []types.Output{{
			TypeName_: "ContractOutput",
			ContractOutput: &types.ContractOutput{
				InputIndex:  0,
				BalanceRoot: types.Bytes32{Hash: common.HexToHash("0x5cc28e489493724ead2b4b771c276109f3995bf527b1c038b7cc14dbfe92758c")},
				StateRoot:   types.Bytes32{Hash: common.HexToHash("0xacfcb3e44140f37fea2c50ba5fafe00920a3b56ffffbe1e233880b5ee4abcf08")},
			},
		}, {
			TypeName_: "ChangeOutput",
			ChangeOutput: &types.ChangeOutput{
				To:      types.Address{Hash: common.HexToHash("0xe173edec3aad8af6d0735165fc013527e93316af30d035dc465dbb6f37cb6972")},
				Amount:  1968969,
				AssetId: types.AssetId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
			},
		}},
		Witnesses: []types.HexString{{
			Bytes: common.FromHex("0x5d9106ea5d5ca649cf1383ebb687b04f9acac9a1abc3ad119b165bccdb7792792fc88b094abad7204c4732515ca290b6026f699146e73dcb65defaad2fa06eaf"),
		}},
		ReceiptsRoot: &types.Bytes32{Hash: common.HexToHash("0x387ae9a320be556f6b79aa5738f59f9381d1919d30b9daacb0da4d412bf7eba3")},
		Status: &types.TransactionStatus{
			TypeName_: "SuccessStatus",
			SuccessStatus: &types.SuccessStatus{
				TransactionId: types.TransactionId{Hash: common.HexToHash("0x9b7a9353faacd4ce91c47707d66c81ec7e4d547905168a592312a94a5c67b69f")},
				Block: types.Block{
					Id: types.BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
					Header: types.Header{
						Id:                  types.BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
						DaHeight:            5700482,
						TransactionsCount:   3,
						MessageReceiptCount: 0,
						TransactionsRoot:    types.Bytes32{Hash: common.HexToHash("0x6acba90c0da2a5946cde70bc5d211ca06f1903b0fe7318bf7653ad4de3caf004")},
						MessageReceiptRoot:  types.Bytes32{Hash: common.HexToHash("0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")},
						Height:              9758550,
						PrevRoot:            types.Bytes32{Hash: common.HexToHash("0xe14198c9e1fbc499df5a9dacdb1219135a2d4915011962b5ac379c54b9499b83")},
						Time:                types.Tai64Timestamp{Time: time.Date(2024, time.April, 15, 2, 44, 2, 0, time.UTC)},
						ApplicationHash:     types.Bytes32{Hash: common.HexToHash("0xe0c1360865782cc46da4f65787896aa264e4e8812b6fdb7864cdf7ef6bf42437")},
					},
					Consensus: types.Consensus{
						TypeName_: "PoAConsensus",
						PoAConsensus: &types.PoAConsensus{
							Signature: types.Signature{
								Bytes: common.FromHex("0x765a5f984189fde36733774e7d76bd9ffcdaa17850a87bd86430addc5c0855923f527871c7cd83617006a31c42d1ca3d0438dc3681b03d6362a90f4498e1db40"),
							},
						},
					},
				},
				Time: types.Tai64Timestamp{Time: time.Date(2024, time.April, 15, 2, 44, 2, 0, time.UTC)},
				ProgramState: &types.ProgramState{
					ReturnType: "RETURN",
					Data:       types.HexString{Bytes: common.FromHex("0x0000000000000001")},
				},
			},
		},
		Receipts: []types.Receipt{{
			Pc: util.GetPointer[types.U64](11640),
			Is: util.GetPointer[types.U64](11640),
			To: &types.Contract{
				Id: types.ContractId{Hash: common.HexToHash("0xd2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f")},
			},
			Amount:      util.GetPointer[types.U64](0),
			AssetId:     &types.AssetId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
			Gas:         util.GetPointer[types.U64](765817),
			Param1:      util.GetPointer[types.U64](3918102790),
			Param2:      util.GetPointer[types.U64](10448),
			ReceiptType: "CALL",
		}, {
			Contract: &types.Contract{
				Id: types.ContractId{Hash: common.HexToHash("0xd2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f")},
			},
			Pc:          util.GetPointer[types.U64](44000),
			Is:          util.GetPointer[types.U64](11640),
			Val:         util.GetPointer[types.U64](0),
			ReceiptType: "RETURN",
		}, {
			Pc:          util.GetPointer[types.U64](10356),
			Is:          util.GetPointer[types.U64](10336),
			Val:         util.GetPointer[types.U64](1),
			ReceiptType: "RETURN",
		}, {
			ReceiptType: "SCRIPT_RESULT",
			Result:      util.GetPointer[types.U64](0),
			GasUsed:     util.GetPointer[types.U64](406656),
		}},
		Script:     &types.HexString{Bytes: common.FromHex("0x724028a0724428785d451000724828802d41148a24040000")},
		ScriptData: &types.HexString{Bytes: common.FromHex("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000d2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f00000000e989810600000000000028d000000000000000000000000000000003")},
		RawPayload: types.HexString{Bytes: common.FromHex("0x000000000000000000000000000c3500000000000000001800000000000000680000000000000001000000000000000200000000000000020000000000000001387ae9a320be556f6b79aa5738f59f9381d1919d30b9daacb0da4d412bf7eba3724028a0724428785d451000724828802d41148a2404000000000000000000000000000000000000000000000000000000000000000000000000000000000000d2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f00000000e989810600000000000028d0000000000000000000000000000000030000000000000001000000000000000116d77c9e9f146e5523a14fd428aa57a839ab775245ac9d5b662ee6f5e99fed2300000000000000005cc28e489493724ead2b4b771c276109f3995bf527b1c038b7cc14dbfe92758c8f36f4ef87d3260fcbbb8b7d047bae772b12265d9c45bb11814855d57fdacee3000000000094e7550000000000000005d2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f000000000000000088dd1c739b7539af1d82ee96d26f528ab7e0ea2485f9c0febb4c2bdb884c19f10000000000000001e173edec3aad8af6d0735165fc013527e93316af30d035dc465dbb6f37cb697200000000001e20bd0000000000000000000000000000000000000000000000000000000000000000000000000094e74d000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000005cc28e489493724ead2b4b771c276109f3995bf527b1c038b7cc14dbfe92758cacfcb3e44140f37fea2c50ba5fafe00920a3b56ffffbe1e233880b5ee4abcf080000000000000002e173edec3aad8af6d0735165fc013527e93316af30d035dc465dbb6f37cb697200000000001e0b49000000000000000000000000000000000000000000000000000000000000000000000000000000405d9106ea5d5ca649cf1383ebb687b04f9acac9a1abc3ad119b165bccdb7792792fc88b094abad7204c4732515ca290b6026f699146e73dcb65defaad2fa06eaf")},
	}
	exp := structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
		"id": structpb.NewStringValue("0x9b7a9353faacd4ce91c47707d66c81ec7e4d547905168a592312a94a5c67b69f"),
		"inputAssetIds": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStringValue("0x0000000000000000000000000000000000000000000000000000000000000000"),
		}}),
		"inputContracts": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"id":       structpb.NewStringValue("0xd2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f"),
				"bytecode": structpb.NewStringValue("0x"),
				"salt":     structpb.NewStringValue(""),
			}}),
		}}),
		"policies": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
			"gasPrice": structpb.NewStringValue("1"),
		}}),
		"gasPrice":       structpb.NewStringValue("1"),
		"scriptGasLimit": structpb.NewStringValue("800000"),
		"maturity":       structpb.NewStringValue("0"),
		"isScript":       structpb.NewBoolValue(true),
		"isCreate":       structpb.NewBoolValue(false),
		"isMint":         structpb.NewBoolValue(false),
		"inputs": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"__typename":  structpb.NewStringValue("InputContract"),
				"utxoId":      structpb.NewStringValue("0x16d77c9e9f146e5523a14fd428aa57a839ab775245ac9d5b662ee6f5e99fed2300"),
				"balanceRoot": structpb.NewStringValue("0x5cc28e489493724ead2b4b771c276109f3995bf527b1c038b7cc14dbfe92758c"),
				"stateRoot":   structpb.NewStringValue("0x8f36f4ef87d3260fcbbb8b7d047bae772b12265d9c45bb11814855d57fdacee3"),
				"txPointer":   structpb.NewStringValue("0094e7550005"),
				"contract": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"id":       structpb.NewStringValue("0xd2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f"),
					"bytecode": structpb.NewStringValue("0x"),
					"salt":     structpb.NewStringValue(""),
				}}),
			}}),
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"__typename":       structpb.NewStringValue("InputCoin"),
				"utxoId":           structpb.NewStringValue("0x88dd1c739b7539af1d82ee96d26f528ab7e0ea2485f9c0febb4c2bdb884c19f101"),
				"owner":            structpb.NewStringValue("0xe173edec3aad8af6d0735165fc013527e93316af30d035dc465dbb6f37cb6972"),
				"amount":           structpb.NewStringValue("1974461"),
				"assetId":          structpb.NewStringValue("0x0000000000000000000000000000000000000000000000000000000000000000"),
				"txPointer":        structpb.NewStringValue("0094e74d0001"),
				"witnessIndex":     structpb.NewStringValue("0"),
				"maturity":         structpb.NewStringValue("0"),
				"predicateGasUsed": structpb.NewStringValue("0"),
				"predicate":        structpb.NewStringValue("0x"),
				"predicateData":    structpb.NewStringValue("0x"),
			}}),
		}}),
		"outputs": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"__typename":  structpb.NewStringValue("ContractOutput"),
				"inputIndex":  structpb.NewStringValue("0"),
				"balanceRoot": structpb.NewStringValue("0x5cc28e489493724ead2b4b771c276109f3995bf527b1c038b7cc14dbfe92758c"),
				"stateRoot":   structpb.NewStringValue("0xacfcb3e44140f37fea2c50ba5fafe00920a3b56ffffbe1e233880b5ee4abcf08"),
			}}),
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"__typename": structpb.NewStringValue("ChangeOutput"),
				"to":         structpb.NewStringValue("0xe173edec3aad8af6d0735165fc013527e93316af30d035dc465dbb6f37cb6972"),
				"amount":     structpb.NewStringValue("1968969"),
				"assetId":    structpb.NewStringValue("0x0000000000000000000000000000000000000000000000000000000000000000"),
			}}),
		}}),
		"witnesses": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStringValue("0x5d9106ea5d5ca649cf1383ebb687b04f9acac9a1abc3ad119b165bccdb7792792fc88b094abad7204c4732515ca290b6026f699146e73dcb65defaad2fa06eaf"),
		}}),
		"receiptsRoot": structpb.NewStringValue("0x387ae9a320be556f6b79aa5738f59f9381d1919d30b9daacb0da4d412bf7eba3"),
		"status": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
			"__typename":    structpb.NewStringValue("SuccessStatus"),
			"transactionId": structpb.NewStringValue("0x9b7a9353faacd4ce91c47707d66c81ec7e4d547905168a592312a94a5c67b69f"),
			"block": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"id": structpb.NewStringValue("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c"),
				"header": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"id":                  structpb.NewStringValue("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c"),
					"daHeight":            structpb.NewStringValue("5700482"),
					"transactionsCount":   structpb.NewStringValue("3"),
					"messageReceiptCount": structpb.NewStringValue("0"),
					"transactionsRoot":    structpb.NewStringValue("0x6acba90c0da2a5946cde70bc5d211ca06f1903b0fe7318bf7653ad4de3caf004"),
					"messageReceiptRoot":  structpb.NewStringValue("0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"),
					"height":              structpb.NewStringValue("9758550"),
					"prevRoot":            structpb.NewStringValue("0xe14198c9e1fbc499df5a9dacdb1219135a2d4915011962b5ac379c54b9499b83"),
					"time":                structpb.NewStringValue("4611686020140536983"),
					"applicationHash":     structpb.NewStringValue("0xe0c1360865782cc46da4f65787896aa264e4e8812b6fdb7864cdf7ef6bf42437"),
				}}),
				"consensus": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"__typename": structpb.NewStringValue("PoAConsensus"),
					"signature":  structpb.NewStringValue("0x765a5f984189fde36733774e7d76bd9ffcdaa17850a87bd86430addc5c0855923f527871c7cd83617006a31c42d1ca3d0438dc3681b03d6362a90f4498e1db40"),
				}}),
				"transactions": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
			}}),
			"receipts": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
			"time":     structpb.NewStringValue("4611686020140536983"),
			"programState": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"returnType": structpb.NewStringValue("RETURN"),
				"data":       structpb.NewStringValue("0x0000000000000001"),
			}}),
		}}),
		"receipts": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"pc": structpb.NewStringValue("11640"),
				"is": structpb.NewStringValue("11640"),
				"to": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"id":       structpb.NewStringValue("0xd2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f"),
					"bytecode": structpb.NewStringValue("0x"),
					"salt":     structpb.NewStringValue(""),
				}}),
				"amount":      structpb.NewStringValue("0"),
				"assetId":     structpb.NewStringValue("0x0000000000000000000000000000000000000000000000000000000000000000"),
				"gas":         structpb.NewStringValue("765817"),
				"param1":      structpb.NewStringValue("3918102790"),
				"param2":      structpb.NewStringValue("10448"),
				"receiptType": structpb.NewStringValue("CALL"),
			}}),
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"contract": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
					"id":       structpb.NewStringValue("0xd2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f"),
					"bytecode": structpb.NewStringValue("0x"),
					"salt":     structpb.NewStringValue(""),
				}}),
				"pc":          structpb.NewStringValue("44000"),
				"is":          structpb.NewStringValue("11640"),
				"val":         structpb.NewStringValue("0"),
				"receiptType": structpb.NewStringValue("RETURN"),
			}}),
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"pc":          structpb.NewStringValue("10356"),
				"is":          structpb.NewStringValue("10336"),
				"val":         structpb.NewStringValue("1"),
				"receiptType": structpb.NewStringValue("RETURN"),
			}}),
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"receiptType": structpb.NewStringValue("SCRIPT_RESULT"),
				"result":      structpb.NewStringValue("0"),
				"gasUsed":     structpb.NewStringValue("406656"),
			}}),
		}}),
		"script":       structpb.NewStringValue("0x724028a0724428785d451000724828802d41148a24040000"),
		"scriptData":   structpb.NewStringValue("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000d2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f00000000e989810600000000000028d000000000000000000000000000000003"),
		"storageSlots": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{}}),
		"rawPayload":   structpb.NewStringValue("0x000000000000000000000000000c3500000000000000001800000000000000680000000000000001000000000000000200000000000000020000000000000001387ae9a320be556f6b79aa5738f59f9381d1919d30b9daacb0da4d412bf7eba3724028a0724428785d451000724828802d41148a2404000000000000000000000000000000000000000000000000000000000000000000000000000000000000d2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f00000000e989810600000000000028d0000000000000000000000000000000030000000000000001000000000000000116d77c9e9f146e5523a14fd428aa57a839ab775245ac9d5b662ee6f5e99fed2300000000000000005cc28e489493724ead2b4b771c276109f3995bf527b1c038b7cc14dbfe92758c8f36f4ef87d3260fcbbb8b7d047bae772b12265d9c45bb11814855d57fdacee3000000000094e7550000000000000005d2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f000000000000000088dd1c739b7539af1d82ee96d26f528ab7e0ea2485f9c0febb4c2bdb884c19f10000000000000001e173edec3aad8af6d0735165fc013527e93316af30d035dc465dbb6f37cb697200000000001e20bd0000000000000000000000000000000000000000000000000000000000000000000000000094e74d000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000005cc28e489493724ead2b4b771c276109f3995bf527b1c038b7cc14dbfe92758cacfcb3e44140f37fea2c50ba5fafe00920a3b56ffffbe1e233880b5ee4abcf080000000000000002e173edec3aad8af6d0735165fc013527e93316af30d035dc465dbb6f37cb697200000000001e0b49000000000000000000000000000000000000000000000000000000000000000000000000000000405d9106ea5d5ca649cf1383ebb687b04f9acac9a1abc3ad119b165bccdb7792792fc88b094abad7204c4732515ca290b6026f699146e73dcb65defaad2fa06eaf"),
	}})
	//tt1, _ := json.MarshalIndent(exp, "", "  ")
	//fmt.Printf("!!! ===1 %s\n", string(tt1))
	//tt2, _ := json.MarshalIndent(txn.MarshalStructpb(), "", "  ")
	//fmt.Printf("!!! ===2 %s\n", string(tt2))
	//assert.Equal(t, string(tt1), string(tt2))
	assert.Equal(t, exp, txn.MarshalStructpb())
}
