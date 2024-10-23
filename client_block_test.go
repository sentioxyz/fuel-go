package fuel

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sentioxyz/fuel-go/types"
	"github.com/sentioxyz/fuel-go/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_GetBlock0(t *testing.T) {
	cli := NewClient(testnetEndpoint)

	block, err := cli.GetBlock(context.Background(), types.QueryBlockParams{
		Height: util.GetPointer(types.U32(1067005)),
	}, GetBlockOption{})
	assert.NoError(t, err)
	assert.Equal(t, &types.Block{
		Version: "V1",
		Id:      types.BlockId{Hash: common.HexToHash("0x4e02668366cbdc2ea9197fa3a84e57e723028de5fe4f574ccaa7b6b744ced495")},
		Height:  1067005,
	}, block)
}

var (
	header_1067005 = types.Header{
		Version:                        "V1",
		Id:                             types.BlockId{Hash: common.HexToHash("0x4e02668366cbdc2ea9197fa3a84e57e723028de5fe4f574ccaa7b6b744ced495")},
		DaHeight:                       5997348,
		ConsensusParametersVersion:     1,
		StateTransitionBytecodeVersion: 0,
		TransactionsCount:              2,
		MessageReceiptCount:            0,
		TransactionsRoot:               types.Bytes32{Hash: common.HexToHash("0x27dffa568c95bb9007fe7f41c8c3acaf2beba08563415d975386fda526ef4b88")},
		MessageOutboxRoot:              types.Bytes32{Hash: common.HexToHash("0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")},
		EventInboxRoot:                 types.Bytes32{Hash: common.HexToHash("0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")},
		Height:                         1067005,
		PrevRoot:                       types.Bytes32{Hash: common.HexToHash("0x184ff0a0c46958b4d5594330765ce7400bb1d7c66a34873860a519c5e765cf37")},
		Time:                           types.Tai64Timestamp{Time: time.Date(2024, time.May, 29, 1, 43, 37, 0, time.UTC)},
		ApplicationHash:                types.Bytes32{Hash: common.HexToHash("0x4abb34eacf7f6e70ffe0f94d0bdcd9833fd37980505a717cab7b2c27cd16ee77")},
	}
	consensus_1067005 = types.Consensus{
		TypeName_: "PoAConsensus",
		PoAConsensus: &types.PoAConsensus{
			Signature: types.Signature{
				Bytes: common.FromHex("0x63d99146a049c0e124b4698842b00c04c54f51533fdf6b7cfd51f2ed9b3f8a074391d6dc1947cde05eea62c4ab23e058ec2373e3a66cba6d45a9237b7a26a2a0"),
			},
		},
	}
	txnIdList_1067005 = []types.TransactionId{
		{Hash: common.HexToHash("0x0ec0390a47eb248d579c74861d747259c2a2a3f4c5c4cdccf049f0670b9a4485")},
		{Hash: common.HexToHash("0xb0b4aaafa1df52c844ea4b970d40397ef2880087f8bd8d45619e7cac95b1c0d8")},
	}
)

func Test_GetBlock1(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	block, err := cli.GetBlock(context.Background(), types.QueryBlockParams{
		Height: util.GetPointer(types.U32(1067005)),
	}, GetBlockOption{WithHeader: true, WithConsensus: true})
	assert.NoError(t, err)
	assert.Equal(t, &types.Block{
		Version:   "V1",
		Id:        types.BlockId{Hash: common.HexToHash("0x4e02668366cbdc2ea9197fa3a84e57e723028de5fe4f574ccaa7b6b744ced495")},
		Height:    1067005,
		Header:    header_1067005,
		Consensus: consensus_1067005,
	}, block)
}

func Test_GetBlock2(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	block, err := cli.GetBlock(context.Background(), types.QueryBlockParams{
		Height: util.GetPointer(types.U32(1067005)),
	}, GetBlockOption{WithTransactions: true})
	assert.NoError(t, err)
	assert.Equal(t, &types.Block{
		Version:        "V1",
		Id:             types.BlockId{Hash: common.HexToHash("0x4e02668366cbdc2ea9197fa3a84e57e723028de5fe4f574ccaa7b6b744ced495")},
		Height:         1067005,
		TransactionIds: txnIdList_1067005,
		Transactions: []types.Transaction{{
			Id: types.TransactionId{Hash: common.HexToHash("0x0ec0390a47eb248d579c74861d747259c2a2a3f4c5c4cdccf049f0670b9a4485")},
		}, {
			Id: types.TransactionId{Hash: common.HexToHash("0xb0b4aaafa1df52c844ea4b970d40397ef2880087f8bd8d45619e7cac95b1c0d8")},
		}},
	}, block)
}

func Test_GetBlock3(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	block, err := cli.GetBlock(context.Background(), types.QueryBlockParams{
		Id: &types.BlockId{Hash: common.HexToHash("0x4e02668366cbdc2ea9197fa3a84e57e723028de5fe4f574ccaa7b6b744ced495")},
	}, GetBlockOption{WithTransactions: true, WithTransactionDetail: true, WithTransactionReceipts: true})
	exp := &types.Block{
		Version:        "V1",
		Id:             types.BlockId{Hash: common.HexToHash("0x4e02668366cbdc2ea9197fa3a84e57e723028de5fe4f574ccaa7b6b744ced495")},
		Height:         1067005,
		TransactionIds: txnIdList_1067005,
		Transactions:   []types.Transaction{txn_1067005_0, txn_1067005_1},
	}
	assert.NoError(t, err)
	assert.Equal(t, exp, block)
}

func Test_GetBlock4(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	block, err := cli.GetBlock(context.Background(), types.QueryBlockParams{
		Id: &types.BlockId{Hash: common.HexToHash("0x4e02668366cbdc2ea9197fa3a84e57e723028de5fe4f574ccaa7b6b744ced495")},
	}, GetBlockOption{})
	assert.NoError(t, err)
	assert.Equal(t, &types.Block{
		Version: "V1",
		Id:      types.BlockId{Hash: common.HexToHash("0x4e02668366cbdc2ea9197fa3a84e57e723028de5fe4f574ccaa7b6b744ced495")},
		Height:  1067005,
	}, block)
}

func Test_GetBlock_marshalJSON(t *testing.T) {
	block := types.Block{
		Version: "V1",
		Id:      types.BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
		Height:  9758550,
		Header: types.Header{
			Version: "V1",
			Id:      types.BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
			Height:  9758550,
			Time:    types.Tai64Timestamp{Time: time.Date(2024, time.April, 5, 1, 2, 3, 0, time.UTC)},
		},
		Consensus: types.Consensus{
			TypeName_: "PoAConsensus",
			PoAConsensus: &types.PoAConsensus{
				Signature: types.Signature{Bytes: common.FromHex("0x123456")},
			},
		},
	}
	text, err := json.MarshalIndent(block, "", "  ")
	assert.NoError(t, err)
	assert.Equal(t, `{
  "version": "V1",
  "id": "0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c",
  "height": "9758550",
  "header": {
    "version": "V1",
    "id": "0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c",
    "daHeight": "0",
    "consensusParametersVersion": "0",
    "stateTransitionBytecodeVersion": "0",
    "transactionsCount": "0",
    "messageReceiptCount": "0",
    "transactionsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "messageOutboxRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "eventInboxRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "height": "9758550",
    "prevRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "time": "4611686020139666864",
    "applicationHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
  },
  "consensus": {
    "__typename": "PoAConsensus",
    "signature": "0x123456"
  },
  "transactionIds": null,
  "transactions": null
}`, string(text))

	var block2 types.Block
	assert.NoError(t, json.Unmarshal(text, &block2))
	assert.Equal(t, block, block2)
}

func Test_GetBlocks(t *testing.T) {
	cli := NewClientWithLogger(testnetEndpoint, SimpleLogger)
	blocks, err := cli.GetBlocks(context.Background(), []types.QueryBlockParams{{
		Id: &types.BlockId{Hash: common.HexToHash("0x4e02668366cbdc2ea9197fa3a84e57e723028de5fe4f574ccaa7b6b744ced495")},
	}, {
		Id: &types.BlockId{Hash: common.HexToHash("0xae197dbe746cc36784aaff0cf38e1678dbc6ee6b85b15d6bb2c3e932c2ad156d")},
	}}, GetBlockOption{})
	assert.NoError(t, err)
	assert.Equal(t, []*types.Block{{
		Version: "V1",
		Id:      types.BlockId{Hash: common.HexToHash("0x4e02668366cbdc2ea9197fa3a84e57e723028de5fe4f574ccaa7b6b744ced495")},
		Height:  1067005,
	}, {
		Version: "V1",
		Id:      types.BlockId{Hash: common.HexToHash("0xae197dbe746cc36784aaff0cf38e1678dbc6ee6b85b15d6bb2c3e932c2ad156d")},
		Height:  1067006,
	}}, blocks)
}
