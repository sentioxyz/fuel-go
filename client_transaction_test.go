package fuel

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sentioxyz/fuel-go/types"
	"github.com/sentioxyz/fuel-go/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_GetTransaction0(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	txn, err := cli.GetTransaction(context.Background(), types.QueryTransactionParams{
		Id: types.TransactionId{Hash: common.HexToHash("0x9b7a9353faacd4ce91c47707d66c81ec7e4d547905168a592312a94a5c67b690")},
	}, GetTransactionOption{
		WithReceipts: true,
		WithStatus:   true,
	})
	assert.NoError(t, err)
	assert.Nil(t, txn)
}

var (
	txn_1067005_0 = types.Transaction{
		Id:             types.TransactionId{Hash: common.HexToHash("0x0ec0390a47eb248d579c74861d747259c2a2a3f4c5c4cdccf049f0670b9a4485")},
		InputAssetIds:  []types.AssetId{{Hash: common.HexToHash("0xf8f8b6283d7fa5b672b530cbb84fcccb4ff8dc40f8176ef4544ddb1f1952ad07")}},
		InputContracts: []types.ContractId{{Hash: common.HexToHash("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850")}},
		InputContract:  nil,
		Policies: &types.Policies{
			Tip:          nil,
			WitnessLimit: nil,
			Maturity:     nil,
			MaxFee:       util.GetPointer[types.U64](59243),
		},
		ScriptGasLimit: util.GetPointer[types.U64](100000),
		Maturity:       util.GetPointer[types.U32](0),
		MintAmount:     nil,
		MintAssetId:    nil,
		MintGasPrice:   nil,
		TxPointer:      nil,
		IsScript:       true,
		IsCreate:       false,
		IsMint:         false,
		IsUpgrade:      false,
		IsUpload:       false,
		IsBlob:         false,
		Inputs: []types.Input{{
			TypeName_: "InputContract",
			InputContract: &types.InputContract{
				UtxoId:      types.UtxoId{Bytes: common.FromHex("0x4ed2af7ccf2e111c376d0fe486396f594b360d93d694875361d56f43705c3fdf0000")},
				BalanceRoot: types.Bytes32{Hash: common.HexToHash("0x313d948ce814d576ffe3dfdf36758dbd51356d9d2a3adcba12001939a4442b17")},
				StateRoot:   types.Bytes32{Hash: common.HexToHash("0x6021474cda1220f0dcae5824b3ea73beb5752b474e648701c85197ea530591ed")},
				TxPointer:   "0008b29d0000",
				ContractId:  types.ContractId{Hash: common.HexToHash("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850")},
			},
		}, {
			TypeName_: "InputCoin",
			InputCoin: &types.InputCoin{
				UtxoId:           types.UtxoId{Bytes: common.FromHex("0x1799f59d5f3ee48479e18feb0e5e705aa4d53f677ba21630ec85200381e621e90000")},
				Owner:            types.Address{Hash: common.HexToHash("0xd3fe20c8ff68a4054d8587ac170c40db7d1e200208a575780542bd9a7e3eec08")},
				Amount:           2000000,
				AssetId:          types.AssetId{Hash: common.HexToHash("0xf8f8b6283d7fa5b672b530cbb84fcccb4ff8dc40f8176ef4544ddb1f1952ad07")},
				TxPointer:        "001046f00000",
				WitnessIndex:     0,
				PredicateGasUsed: 0,
				Predicate:        types.HexString{Bytes: common.FromHex("0x")},
				PredicateData:    types.HexString{Bytes: common.FromHex("0x")},
			},
		}},
		Outputs: []types.Output{{
			TypeName_: "ContractOutput",
			ContractOutput: &types.ContractOutput{
				InputIndex:  0,
				BalanceRoot: types.Bytes32{Hash: common.HexToHash("0x313d948ce814d576ffe3dfdf36758dbd51356d9d2a3adcba12001939a4442b17")},
				StateRoot:   types.Bytes32{Hash: common.HexToHash("0xeb827e501eee0dff74b423e5173bcd9e3147d00fb2da246c2bfe8abc98660d33")},
			},
		}, {
			TypeName_: "ChangeOutput",
			ChangeOutput: &types.ChangeOutput{
				To:      types.Address{Hash: common.HexToHash("0xd3fe20c8ff68a4054d8587ac170c40db7d1e200208a575780542bd9a7e3eec08")},
				Amount:  1941078,
				AssetId: types.AssetId{Hash: common.HexToHash("0xf8f8b6283d7fa5b672b530cbb84fcccb4ff8dc40f8176ef4544ddb1f1952ad07")},
			},
		}},
		OutputContract: nil,
		Witnesses: []types.HexString{{
			Bytes: common.FromHex("0xf34f6666876b62e6930f7ed447c6ffd86344237325cd4db19f2ce1a0be2cc76b995e73e04f039f7ad60f614edf6b85667e876652499270d111cf2acc7ca41b5b"),
		}},
		ReceiptsRoot: &types.Bytes32{Hash: common.HexToHash("0x74ab41b67d3e1ccb4bda92dfac21bfb448f80d0500b3cd56480aaee2fa37cbf0")},
		Status: &types.TransactionStatus{
			TypeName_: "SuccessStatus",
			SuccessStatus: &types.SuccessStatus{
				TransactionId: types.TransactionId{Hash: common.HexToHash("0x0ec0390a47eb248d579c74861d747259c2a2a3f4c5c4cdccf049f0670b9a4485")},
				BlockHeight:   1067005,
				Time:          types.Tai64Timestamp{Time: time.Date(2024, time.May, 29, 1, 43, 37, 0, time.UTC)},
				ProgramState: &types.ProgramState{
					ReturnType: "RETURN",
					Data:       types.HexString{Bytes: common.FromHex("0x0000000000000001")},
				},
				Receipts: []types.Receipt{{
					Pc:          util.GetPointer[types.U64](11712),
					Is:          util.GetPointer[types.U64](11712),
					To:          &types.ContractId{Hash: common.HexToHash("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850")},
					Amount:      util.GetPointer[types.U64](0),
					AssetId:     &types.AssetId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
					Gas:         util.GetPointer[types.U64](80154),
					Param1:      util.GetPointer[types.U64](10480),
					Param2:      util.GetPointer[types.U64](10497),
					ReceiptType: "CALL",
				}, {
					Id:          &types.ContractId{Hash: common.HexToHash("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850")},
					Pc:          util.GetPointer[types.U64](16640),
					Is:          util.GetPointer[types.U64](11712),
					Ptr:         util.GetPointer[types.U64](67107840),
					Digest:      &types.Bytes32{Hash: common.HexToHash("0x979f4ef2bcab47538ec8c1c92b8ec5c58ef61b7aa86a708fb9823324479144d6")},
					Ra:          util.GetPointer[types.U64](0),
					Rb:          util.GetPointer[types.U64](14631882454972838106),
					Len:         util.GetPointer[types.U64](25),
					ReceiptType: "LOG_DATA",
					Data:        &types.HexString{Bytes: common.FromHex("0x63616c6c696e6720696e6372656d656e74206d6574686f642e")},
				}, {
					Id:          &types.ContractId{Hash: common.HexToHash("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850")},
					Pc:          util.GetPointer[types.U64](16896),
					Is:          util.GetPointer[types.U64](11712),
					Ptr:         util.GetPointer[types.U64](67106816),
					Digest:      &types.Bytes32{Hash: common.HexToHash("0x17eb70034b5b71092521d184c5e7b069d47de657e51aef2be11a00c115036943")},
					Ra:          util.GetPointer[types.U64](0),
					Rb:          util.GetPointer[types.U64](15520703124961489725),
					Len:         util.GetPointer[types.U64](4),
					ReceiptType: "LOG_DATA",
					Data:        &types.HexString{Bytes: common.FromHex("0x00000008")},
				}, {
					Id:          &types.ContractId{Hash: common.HexToHash("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850")},
					Pc:          util.GetPointer[types.U64](17132),
					Is:          util.GetPointer[types.U64](11712),
					Ptr:         util.GetPointer[types.U64](67105792),
					Digest:      &types.Bytes32{Hash: common.HexToHash("0xcd2662154e6d76b2b2b92e70c0cac3ccf534f9b74eb5b89819ec509083d00a50")},
					Ra:          util.GetPointer[types.U64](0),
					Rb:          util.GetPointer[types.U64](1515152261580153489),
					Len:         util.GetPointer[types.U64](8),
					ReceiptType: "LOG_DATA",
					Data:        &types.HexString{Bytes: common.FromHex("0x0000000000000001")},
				}, {
					Id:          &types.ContractId{Hash: common.HexToHash("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850")},
					Pc:          util.GetPointer[types.U64](17380),
					Is:          util.GetPointer[types.U64](11712),
					Ptr:         util.GetPointer[types.U64](67104768),
					Digest:      &types.Bytes32{Hash: common.HexToHash("0x02d449a31fbb267c8f352e9968a79e3e5fc95c1bbeaa502fd6454ebde5a4bedc")},
					Ra:          util.GetPointer[types.U64](0),
					Rb:          util.GetPointer[types.U64](8961848586872524460),
					Len:         util.GetPointer[types.U64](32),
					ReceiptType: "LOG_DATA",
					Data:        &types.HexString{Bytes: common.FromHex("0x1111111111111111111111111111111111111111111111111111111111111111")},
				}, {
					Id:          &types.ContractId{Hash: common.HexToHash("0xdb0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e850")},
					Pc:          util.GetPointer[types.U64](18408),
					Is:          util.GetPointer[types.U64](11712),
					Ptr:         util.GetPointer[types.U64](67103232),
					Digest:      &types.Bytes32{Hash: common.HexToHash("0x5bc67471c189d78c76461dcab6141a733bdab3799d1d69e0c419119c92e82b3d")},
					Len:         util.GetPointer[types.U64](8),
					ReceiptType: "RETURN_DATA",
					Data:        &types.HexString{Bytes: common.FromHex("0x0000000000000012")},
				}, {
					Pc:          util.GetPointer[types.U64](10388),
					Is:          util.GetPointer[types.U64](10368),
					Val:         util.GetPointer[types.U64](1),
					ReceiptType: "RETURN",
				}, {
					ReceiptType: "SCRIPT_RESULT",
					Result:      util.GetPointer[types.U64](0),
					GasUsed:     util.GetPointer[types.U64](70408),
				}},
				TotalGas: 5420740,
				TotalFee: 58922,
			},
		},
		Script:               &types.HexString{Bytes: common.FromHex("0x724028c0724428985d451000724828a02d41148a24040000")},
		ScriptData:           &types.HexString{Bytes: common.FromHex("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000db0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e85000000000000028f000000000000029010000000000000009696e6372656d656e740000000800000000000000011111111111111111111111111111111111111111111111111111111111111111")},
		BytecodeWitnessIndex: nil,
		BlobId:               nil,
		Salt:                 nil,
		StorageSlots:         nil,
		BytecodeRoot:         nil,
		SubsectionIndex:      nil,
		SubsectionsNumber:    nil,
		ProofSet:             nil,
		UpgradePurpose:       nil,
		RawPayload:           types.HexString{Bytes: common.FromHex("0x000000000000000000000000000186a074ab41b67d3e1ccb4bda92dfac21bfb448f80d0500b3cd56480aaee2fa37cbf0000000000000001800000000000000950000000000000008000000000000000200000000000000020000000000000001724028c0724428985d451000724828a02d41148a2404000000000000000000000000000000000000000000000000000000000000000000000000000000000000db0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e85000000000000028f000000000000029010000000000000009696e6372656d656e740000000800000000000000011111111111111111111111111111111111111111111111111111111111111111000000000000000000e76b00000000000000014ed2af7ccf2e111c376d0fe486396f594b360d93d694875361d56f43705c3fdf0000000000000000313d948ce814d576ffe3dfdf36758dbd51356d9d2a3adcba12001939a4442b176021474cda1220f0dcae5824b3ea73beb5752b474e648701c85197ea530591ed000000000008b29d0000000000000000db0d550935d601c45791ba18664f0a821c11745b1f938e87f10a79e21988e85000000000000000001799f59d5f3ee48479e18feb0e5e705aa4d53f677ba21630ec85200381e621e90000000000000000d3fe20c8ff68a4054d8587ac170c40db7d1e200208a575780542bd9a7e3eec0800000000001e8480f8f8b6283d7fa5b672b530cbb84fcccb4ff8dc40f8176ef4544ddb1f1952ad0700000000001046f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000313d948ce814d576ffe3dfdf36758dbd51356d9d2a3adcba12001939a4442b17eb827e501eee0dff74b423e5173bcd9e3147d00fb2da246c2bfe8abc98660d330000000000000002d3fe20c8ff68a4054d8587ac170c40db7d1e200208a575780542bd9a7e3eec0800000000001d9e56f8f8b6283d7fa5b672b530cbb84fcccb4ff8dc40f8176ef4544ddb1f1952ad070000000000000040f34f6666876b62e6930f7ed447c6ffd86344237325cd4db19f2ce1a0be2cc76b995e73e04f039f7ad60f614edf6b85667e876652499270d111cf2acc7ca41b5b")},
	}

	txn_1067005_1 = types.Transaction{
		Id:            types.TransactionId{Hash: common.HexToHash("0xb0b4aaafa1df52c844ea4b970d40397ef2880087f8bd8d45619e7cac95b1c0d8")},
		InputAssetIds: nil,
		InputContracts: []types.ContractId{
			{Hash: common.HexToHash("0x7777777777777777777777777777777777777777777777777777777777777777")},
		},
		InputContract: &types.InputContract{
			UtxoId:      types.UtxoId{Bytes: common.FromHex("0xe3567832f02269cc7e4656f764df9f220b89e38057b92a8b5fc72afd661e7b4d0000")},
			BalanceRoot: types.Bytes32{Hash: common.HexToHash("0x116f5d05fe944663faaf956488db0147e8bfcd7b471b918b023969dda788a103")},
			StateRoot:   types.Bytes32{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
			TxPointer:   "001047fc0000",
			ContractId:  types.ContractId{Hash: common.HexToHash("0x7777777777777777777777777777777777777777777777777777777777777777")},
		},
		Policies:       nil,
		ScriptGasLimit: nil,
		Maturity:       nil,
		MintAmount:     util.GetPointer[types.U64](58922),
		MintAssetId:    &types.AssetId{Hash: common.HexToHash("0xf8f8b6283d7fa5b672b530cbb84fcccb4ff8dc40f8176ef4544ddb1f1952ad07")},
		MintGasPrice:   util.GetPointer[types.U64](1),
		TxPointer:      util.GetPointer[types.TxPointer]("001047fd0001"),
		IsScript:       false,
		IsCreate:       false,
		IsMint:         true,
		IsUpgrade:      false,
		IsUpload:       false,
		IsBlob:         false,
		Inputs:         nil,
		Outputs:        []types.Output{},
		OutputContract: &types.ContractOutput{
			InputIndex:  0,
			BalanceRoot: types.Bytes32{Hash: common.HexToHash("0xd167e412c527c48fe0896e33bf8e2555ef2384a40283c842cc7e44163851f7a2")},
			StateRoot:   types.Bytes32{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
		},
		Witnesses:    nil,
		ReceiptsRoot: nil,
		Status: &types.TransactionStatus{
			TypeName_: "SuccessStatus",
			SuccessStatus: &types.SuccessStatus{
				TransactionId: types.TransactionId{Hash: common.HexToHash("0xb0b4aaafa1df52c844ea4b970d40397ef2880087f8bd8d45619e7cac95b1c0d8")},
				BlockHeight:   1067005,
				Time:          types.Tai64Timestamp{Time: time.Date(2024, time.May, 29, 1, 43, 37, 0, time.UTC)},
				ProgramState:  nil,
				Receipts:      []types.Receipt{},
				TotalGas:      0,
				TotalFee:      0,
			},
		},
		Script:               nil,
		ScriptData:           nil,
		BytecodeWitnessIndex: nil,
		BlobId:               nil,
		Salt:                 nil,
		StorageSlots:         nil,
		BytecodeRoot:         nil,
		SubsectionIndex:      nil,
		SubsectionsNumber:    nil,
		ProofSet:             nil,
		UpgradePurpose:       nil,
		RawPayload:           types.HexString{Bytes: common.FromHex("0x000000000000000200000000001047fd0000000000000001e3567832f02269cc7e4656f764df9f220b89e38057b92a8b5fc72afd661e7b4d0000000000000000116f5d05fe944663faaf956488db0147e8bfcd7b471b918b023969dda788a103000000000000000000000000000000000000000000000000000000000000000000000000001047fc000000000000000077777777777777777777777777777777777777777777777777777777777777770000000000000000d167e412c527c48fe0896e33bf8e2555ef2384a40283c842cc7e44163851f7a20000000000000000000000000000000000000000000000000000000000000000000000000000e62af8f8b6283d7fa5b672b530cbb84fcccb4ff8dc40f8176ef4544ddb1f1952ad070000000000000001")},
	}
)

func Test_GetTransaction1(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	txn, err := cli.GetTransaction(context.Background(), types.QueryTransactionParams{
		Id: types.TransactionId{Hash: common.HexToHash("0x0ec0390a47eb248d579c74861d747259c2a2a3f4c5c4cdccf049f0670b9a4485")},
	}, GetTransactionOption{
		WithReceipts: false,
		WithStatus:   false,
	})
	assert.NoError(t, err)
	exp := txn_1067005_0
	exp.Status = nil
	assert.Equal(t, &exp, txn)
}

func Test_GetTransaction2(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	txn, err := cli.GetTransaction(context.Background(), types.QueryTransactionParams{
		Id: types.TransactionId{Hash: common.HexToHash("0x0ec0390a47eb248d579c74861d747259c2a2a3f4c5c4cdccf049f0670b9a4485")},
	}, GetTransactionOption{
		WithReceipts: false,
		WithStatus:   true,
	})
	assert.NoError(t, err)
	succ := *(txn_1067005_0.Status.SuccessStatus)
	succ.Receipts = nil
	succ.Block = types.Block{
		Version:        "V1",
		Id:             types.BlockId{Hash: common.HexToHash("0x4e02668366cbdc2ea9197fa3a84e57e723028de5fe4f574ccaa7b6b744ced495")},
		Height:         1067005,
		Header:         header_1067005,
		Consensus:      consensus_1067005,
		TransactionIds: txnIdList_1067005,
	}
	status := types.TransactionStatus{
		TypeName_:     "SuccessStatus",
		SuccessStatus: &succ,
	}
	exp := txn_1067005_0
	exp.Status = &status
	assert.Equal(t, &exp, txn)
}

func Test_GetTransaction3(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	txn, err := cli.GetTransaction(context.Background(), types.QueryTransactionParams{
		Id: types.TransactionId{Hash: common.HexToHash("0x0ec0390a47eb248d579c74861d747259c2a2a3f4c5c4cdccf049f0670b9a4485")},
	}, GetTransactionOption{
		WithReceipts: true,
		WithStatus:   true,
	})
	assert.NoError(t, err)
	succ := *(txn_1067005_0.Status.SuccessStatus)
	succ.Block = types.Block{
		Version:        "V1",
		Id:             types.BlockId{Hash: common.HexToHash("0x4e02668366cbdc2ea9197fa3a84e57e723028de5fe4f574ccaa7b6b744ced495")},
		Height:         1067005,
		Header:         header_1067005,
		Consensus:      consensus_1067005,
		TransactionIds: txnIdList_1067005,
	}
	status := types.TransactionStatus{
		TypeName_:     "SuccessStatus",
		SuccessStatus: &succ,
	}
	exp := txn_1067005_0
	exp.Status = &status
	assert.Equal(t, &exp, txn)
}
