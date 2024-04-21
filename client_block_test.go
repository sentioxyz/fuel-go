package fuel

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sentioxyz/fuel-go/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_GetBlock0(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	block, err := cli.GetBlock(context.Background(), QueryBlockParams{
		Height: util.GetPointer(U32(9758550)),
	}, GetBlockOption{WithTransactions: false, OnlyTransactionID: false})
	assert.NoError(t, err)
	assert.Equal(t, &Block{
		Id: BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
		Header: Header{
			Id:                  BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
			DaHeight:            5700482,
			TransactionsCount:   3,
			MessageReceiptCount: 0,
			TransactionsRoot:    Bytes32{Hash: common.HexToHash("0x6acba90c0da2a5946cde70bc5d211ca06f1903b0fe7318bf7653ad4de3caf004")},
			MessageReceiptRoot:  Bytes32{Hash: common.HexToHash("0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")},
			Height:              9758550,
			PrevRoot:            Bytes32{Hash: common.HexToHash("0xe14198c9e1fbc499df5a9dacdb1219135a2d4915011962b5ac379c54b9499b83")},
			Time:                Tai64Timestamp{Time: time.Date(2024, time.April, 15, 2, 44, 2, 0, time.UTC)},
			ApplicationHash:     Bytes32{Hash: common.HexToHash("0xe0c1360865782cc46da4f65787896aa264e4e8812b6fdb7864cdf7ef6bf42437")},
		},
		Consensus: Consensus{
			TypeName_: "PoAConsensus",
			PoAConsensus: &PoAConsensus{
				Signature: Signature{
					Bytes: common.FromHex("0x765a5f984189fde36733774e7d76bd9ffcdaa17850a87bd86430addc5c0855923f527871c7cd83617006a31c42d1ca3d0438dc3681b03d6362a90f4498e1db40"),
				},
			},
		},
	}, block)
}

func Test_GetBlock1(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	block, err := cli.GetBlock(context.Background(), QueryBlockParams{
		Height: util.GetPointer(U32(9758550)),
	}, GetBlockOption{WithTransactions: true, OnlyTransactionID: true})
	assert.NoError(t, err)
	assert.Equal(t, &Block{
		Id: BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
		Header: Header{
			Id:                  BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
			DaHeight:            5700482,
			TransactionsCount:   3,
			MessageReceiptCount: 0,
			TransactionsRoot:    Bytes32{Hash: common.HexToHash("0x6acba90c0da2a5946cde70bc5d211ca06f1903b0fe7318bf7653ad4de3caf004")},
			MessageReceiptRoot:  Bytes32{Hash: common.HexToHash("0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")},
			Height:              9758550,
			PrevRoot:            Bytes32{Hash: common.HexToHash("0xe14198c9e1fbc499df5a9dacdb1219135a2d4915011962b5ac379c54b9499b83")},
			Time:                Tai64Timestamp{Time: time.Date(2024, time.April, 15, 2, 44, 2, 0, time.UTC)},
			ApplicationHash:     Bytes32{Hash: common.HexToHash("0xe0c1360865782cc46da4f65787896aa264e4e8812b6fdb7864cdf7ef6bf42437")},
		},
		Consensus: Consensus{
			TypeName_: "PoAConsensus",
			PoAConsensus: &PoAConsensus{
				Signature: Signature{
					Bytes: common.FromHex("0x765a5f984189fde36733774e7d76bd9ffcdaa17850a87bd86430addc5c0855923f527871c7cd83617006a31c42d1ca3d0438dc3681b03d6362a90f4498e1db40"),
				},
			},
		},
		Transactions: []Transaction{{
			Id: TransactionId{Hash: common.HexToHash("0x9b7a9353faacd4ce91c47707d66c81ec7e4d547905168a592312a94a5c67b69f")},
		}, {
			Id: TransactionId{Hash: common.HexToHash("0x4928a04ef03e8146c530249c6a2e97e389a3dd3c00deb6efcb652de0ea62da47")},
		}, {
			Id: TransactionId{Hash: common.HexToHash("0x1a978dcf45d87d2793d7da58d45244d68241aa6363d6a435a38c5fdfeafff178")},
		}},
	}, block)
}

func Test_GetBlock2(t *testing.T) {
	cli := NewClient(testnetEndpoint)
	block, err := cli.GetBlock(context.Background(), QueryBlockParams{
		Id: &BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
	}, GetBlockOption{WithTransactions: true})
	assert.NoError(t, err)
	assert.Equal(t, &Block{
		Id: BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
		Header: Header{
			Id:                  BlockId{Hash: common.HexToHash("0x5d7f48fc777144b21ea760525936db069329dee2ccce509550c1478c1c0b5b2c")},
			DaHeight:            5700482,
			TransactionsCount:   3,
			MessageReceiptCount: 0,
			TransactionsRoot:    Bytes32{Hash: common.HexToHash("0x6acba90c0da2a5946cde70bc5d211ca06f1903b0fe7318bf7653ad4de3caf004")},
			MessageReceiptRoot:  Bytes32{Hash: common.HexToHash("0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")},
			Height:              9758550,
			PrevRoot:            Bytes32{Hash: common.HexToHash("0xe14198c9e1fbc499df5a9dacdb1219135a2d4915011962b5ac379c54b9499b83")},
			Time:                Tai64Timestamp{Time: time.Date(2024, time.April, 15, 2, 44, 2, 0, time.UTC)},
			ApplicationHash:     Bytes32{Hash: common.HexToHash("0xe0c1360865782cc46da4f65787896aa264e4e8812b6fdb7864cdf7ef6bf42437")},
		},
		Consensus: Consensus{
			TypeName_: "PoAConsensus",
			PoAConsensus: &PoAConsensus{
				Signature: Signature{
					Bytes: common.FromHex("0x765a5f984189fde36733774e7d76bd9ffcdaa17850a87bd86430addc5c0855923f527871c7cd83617006a31c42d1ca3d0438dc3681b03d6362a90f4498e1db40"),
				},
			},
		},
		Transactions: []Transaction{{
			Id:            TransactionId{Hash: common.HexToHash("0x9b7a9353faacd4ce91c47707d66c81ec7e4d547905168a592312a94a5c67b69f")},
			InputAssetIds: []AssetId{{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")}},
			InputContracts: []Contract{{
				Id:   ContractId{Hash: common.HexToHash("0xd2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f")},
				Salt: "0x572e0502c9ca4347b88a0faf5b4a36bbbbf3c4c62d4f77ea893f4be7541b42e6",
			}},
			Policies: &Policies{
				GasPrice: util.GetPointer[U64](1),
			},
			GasPrice:       util.GetPointer[U64](1),
			ScriptGasLimit: util.GetPointer[U64](800000),
			Maturity:       util.GetPointer[U32](0),
			IsScript:       true,
			IsCreate:       false,
			IsMint:         false,
			Inputs: []Input{{
				TypeName_: "InputContract",
				InputContract: &InputContract{
					UtxoId:      UtxoId{Bytes: common.FromHex("0x16d77c9e9f146e5523a14fd428aa57a839ab775245ac9d5b662ee6f5e99fed2300")},
					BalanceRoot: Bytes32{Hash: common.HexToHash("0x5cc28e489493724ead2b4b771c276109f3995bf527b1c038b7cc14dbfe92758c")},
					StateRoot:   Bytes32{Hash: common.HexToHash("0x8f36f4ef87d3260fcbbb8b7d047bae772b12265d9c45bb11814855d57fdacee3")},
					TxPointer:   "0094e7550005",
					Contract: Contract{
						Id:   ContractId{Hash: common.HexToHash("0xd2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f")},
						Salt: "0x572e0502c9ca4347b88a0faf5b4a36bbbbf3c4c62d4f77ea893f4be7541b42e6",
					},
				},
			}, {
				TypeName_: "InputCoin",
				InputCoin: &InputCoin{
					UtxoId:           UtxoId{Bytes: common.FromHex("0x88dd1c739b7539af1d82ee96d26f528ab7e0ea2485f9c0febb4c2bdb884c19f101")},
					Owner:            Address{Hash: common.HexToHash("0xe173edec3aad8af6d0735165fc013527e93316af30d035dc465dbb6f37cb6972")},
					Amount:           1974461,
					AssetId:          AssetId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
					TxPointer:        "0094e74d0001",
					WitnessIndex:     0,
					Maturity:         0,
					PredicateGasUsed: 0,
					Predicate:        HexString{Bytes: common.FromHex("0x")},
					PredicateData:    HexString{Bytes: common.FromHex("0x")},
				},
			}},
			Outputs: []Output{{
				TypeName_: "ContractOutput",
				ContractOutput: &ContractOutput{
					InputIndex:  0,
					BalanceRoot: Bytes32{Hash: common.HexToHash("0x5cc28e489493724ead2b4b771c276109f3995bf527b1c038b7cc14dbfe92758c")},
					StateRoot:   Bytes32{Hash: common.HexToHash("0xacfcb3e44140f37fea2c50ba5fafe00920a3b56ffffbe1e233880b5ee4abcf08")},
				},
			}, {
				TypeName_: "ChangeOutput",
				ChangeOutput: &ChangeOutput{
					To:      Address{Hash: common.HexToHash("0xe173edec3aad8af6d0735165fc013527e93316af30d035dc465dbb6f37cb6972")},
					Amount:  1968969,
					AssetId: AssetId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
				},
			}},
			Witnesses: []HexString{{
				Bytes: common.FromHex("0x5d9106ea5d5ca649cf1383ebb687b04f9acac9a1abc3ad119b165bccdb7792792fc88b094abad7204c4732515ca290b6026f699146e73dcb65defaad2fa06eaf"),
			}},
			ReceiptsRoot: &Bytes32{Hash: common.HexToHash("0x387ae9a320be556f6b79aa5738f59f9381d1919d30b9daacb0da4d412bf7eba3")},
			Status: &TransactionStatus{
				TypeName_:     "SuccessStatus",
				SuccessStatus: &SuccessStatus{},
			},
			Receipts: []Receipt{{
				Pc: util.GetPointer[U64](11640),
				Is: util.GetPointer[U64](11640),
				To: &Contract{
					Id:   ContractId{Hash: common.HexToHash("0xd2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f")},
					Salt: "0x572e0502c9ca4347b88a0faf5b4a36bbbbf3c4c62d4f77ea893f4be7541b42e6",
				},
				Amount:      util.GetPointer[U64](0),
				AssetId:     &AssetId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
				Gas:         util.GetPointer[U64](765817),
				Param1:      util.GetPointer[U64](3918102790),
				Param2:      util.GetPointer[U64](10448),
				ReceiptType: "CALL",
			}, {
				Contract: &Contract{
					Id:   ContractId{Hash: common.HexToHash("0xd2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f")},
					Salt: "0x572e0502c9ca4347b88a0faf5b4a36bbbbf3c4c62d4f77ea893f4be7541b42e6",
				},
				Pc:          util.GetPointer[U64](44000),
				Is:          util.GetPointer[U64](11640),
				Val:         util.GetPointer[U64](0),
				ReceiptType: "RETURN",
			}, {
				Pc:          util.GetPointer[U64](10356),
				Is:          util.GetPointer[U64](10336),
				Val:         util.GetPointer[U64](1),
				ReceiptType: "RETURN",
			}, {
				ReceiptType: "SCRIPT_RESULT",
				Result:      util.GetPointer[U64](0),
				GasUsed:     util.GetPointer[U64](406656),
			}},
			Script:     &HexString{Bytes: common.FromHex("0x724028a0724428785d451000724828802d41148a24040000")},
			ScriptData: &HexString{Bytes: common.FromHex("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000d2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f00000000e989810600000000000028d000000000000000000000000000000003")},
			RawPayload: HexString{Bytes: common.FromHex("0x000000000000000000000000000c3500000000000000001800000000000000680000000000000001000000000000000200000000000000020000000000000001387ae9a320be556f6b79aa5738f59f9381d1919d30b9daacb0da4d412bf7eba3724028a0724428785d451000724828802d41148a2404000000000000000000000000000000000000000000000000000000000000000000000000000000000000d2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f00000000e989810600000000000028d0000000000000000000000000000000030000000000000001000000000000000116d77c9e9f146e5523a14fd428aa57a839ab775245ac9d5b662ee6f5e99fed2300000000000000005cc28e489493724ead2b4b771c276109f3995bf527b1c038b7cc14dbfe92758c8f36f4ef87d3260fcbbb8b7d047bae772b12265d9c45bb11814855d57fdacee3000000000094e7550000000000000005d2a93abef5c3f45f48bb9f0736ccfda4c3f32c9c57fc307ab9363ef7712f305f000000000000000088dd1c739b7539af1d82ee96d26f528ab7e0ea2485f9c0febb4c2bdb884c19f10000000000000001e173edec3aad8af6d0735165fc013527e93316af30d035dc465dbb6f37cb697200000000001e20bd0000000000000000000000000000000000000000000000000000000000000000000000000094e74d000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000005cc28e489493724ead2b4b771c276109f3995bf527b1c038b7cc14dbfe92758cacfcb3e44140f37fea2c50ba5fafe00920a3b56ffffbe1e233880b5ee4abcf080000000000000002e173edec3aad8af6d0735165fc013527e93316af30d035dc465dbb6f37cb697200000000001e0b49000000000000000000000000000000000000000000000000000000000000000000000000000000405d9106ea5d5ca649cf1383ebb687b04f9acac9a1abc3ad119b165bccdb7792792fc88b094abad7204c4732515ca290b6026f699146e73dcb65defaad2fa06eaf")},
		}, {
			Id:             TransactionId{Hash: common.HexToHash("0x4928a04ef03e8146c530249c6a2e97e389a3dd3c00deb6efcb652de0ea62da47")},
			InputAssetIds:  []AssetId{{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")}},
			InputContracts: []Contract{},
			Policies: &Policies{
				GasPrice: util.GetPointer[U64](1),
			},
			GasPrice:       util.GetPointer[U64](1),
			ScriptGasLimit: util.GetPointer[U64](1000),
			Maturity:       util.GetPointer[U32](0),
			IsScript:       true,
			IsCreate:       false,
			IsMint:         false,
			Inputs: []Input{{
				TypeName_: "InputCoin",
				InputCoin: &InputCoin{
					UtxoId:           UtxoId{Bytes: common.FromHex("0xfa4e65145a7daf41ef2d731d06f243609c58bee540e4ed65682727156641756b01")},
					Owner:            Address{Hash: common.HexToHash("0xe7f16524b53c8a8fdf84c916941351666715ae76d45ecc82ad81c13d481682db")},
					Amount:           1972805,
					AssetId:          AssetId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
					TxPointer:        "008fe9f70001",
					WitnessIndex:     0,
					Maturity:         0,
					PredicateGasUsed: 0,
					Predicate:        HexString{Bytes: common.FromHex("0x")},
					PredicateData:    HexString{Bytes: common.FromHex("0x")},
				},
			}},
			Outputs: []Output{{
				TypeName_: "CoinOutput",
				CoinOutput: &CoinOutput{
					To:      Address{Hash: common.HexToHash("0x447bd836060bf59574f601c78948c271efad5955880a09429979013884dad4f9")},
					Amount:  22300,
					AssetId: AssetId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
				},
			}, {
				TypeName_: "ChangeOutput",
				ChangeOutput: &ChangeOutput{
					To:      Address{Hash: common.HexToHash("0xe7f16524b53c8a8fdf84c916941351666715ae76d45ecc82ad81c13d481682db")},
					Amount:  1949618,
					AssetId: AssetId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
				},
			}},
			Witnesses: []HexString{{
				Bytes: common.FromHex("0xb5b054daef45956a2439a702af9b4bfaf0ca6614a49fba244c7d0875ef8e5d8071c592f856a17835d8d48ed73eafb8b62abb8487c406dfc3b5ce6829450285e3"),
			}},
			ReceiptsRoot: &Bytes32{Hash: common.HexToHash("0xe7f678a2e8df7da272cf303aff96023da2ab1968b74d86bb92f5b558d38ed6bd")},
			Status: &TransactionStatus{
				TypeName_:     "SuccessStatus",
				SuccessStatus: &SuccessStatus{},
			},
			Receipts: []Receipt{{
				Pc:          util.GetPointer[U64](10336),
				Is:          util.GetPointer[U64](10336),
				Val:         util.GetPointer[U64](0),
				ReceiptType: "RETURN",
			}, {
				ReceiptType: "SCRIPT_RESULT",
				Result:      util.GetPointer[U64](0),
				GasUsed:     util.GetPointer[U64](733),
			}},
			Script:     &HexString{Bytes: common.FromHex("0x24000000")},
			ScriptData: &HexString{Bytes: common.FromHex("0x")},
			RawPayload: HexString{Bytes: common.FromHex("0x000000000000000000000000000003e8000000000000000400000000000000000000000000000001000000000000000100000000000000020000000000000001e7f678a2e8df7da272cf303aff96023da2ab1968b74d86bb92f5b558d38ed6bd240000000000000000000000000000010000000000000000fa4e65145a7daf41ef2d731d06f243609c58bee540e4ed65682727156641756b0000000000000001e7f16524b53c8a8fdf84c916941351666715ae76d45ecc82ad81c13d481682db00000000001e1a45000000000000000000000000000000000000000000000000000000000000000000000000008fe9f70000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000447bd836060bf59574f601c78948c271efad5955880a09429979013884dad4f9000000000000571c00000000000000000000000000000000000000000000000000000000000000000000000000000002e7f16524b53c8a8fdf84c916941351666715ae76d45ecc82ad81c13d481682db00000000001dbfb200000000000000000000000000000000000000000000000000000000000000000000000000000040b5b054daef45956a2439a702af9b4bfaf0ca6614a49fba244c7d0875ef8e5d8071c592f856a17835d8d48ed73eafb8b62abb8487c406dfc3b5ce6829450285e3")},
		}, {
			Id: TransactionId{Hash: common.HexToHash("0x1a978dcf45d87d2793d7da58d45244d68241aa6363d6a435a38c5fdfeafff178")},
			InputContracts: []Contract{{
				Id:   ContractId{Hash: common.HexToHash("0x7777777777777777777777777777777777777777777777777777777777777777")},
				Salt: "0x1bfd51cb31b8d0bc7d93d38f97ab771267d8786ab87073e0c2b8f9ddc44b274e",
			}},
			InputContract: &InputContract{
				UtxoId:      UtxoId{Bytes: common.FromHex("0xae426ee0c79cac25ffe515ca4148e27086669bee7043b23cd380dce443213eff00")},
				BalanceRoot: Bytes32{Hash: common.HexToHash("0x2d19f8c34395032b25ae83bf88cd618a8598c6f2459f137c931879b323a41e04")},
				StateRoot:   Bytes32{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
				TxPointer:   "0094e7550006",
				Contract: Contract{
					Id:   ContractId{Hash: common.HexToHash("0x7777777777777777777777777777777777777777777777777777777777777777")},
					Salt: "0x1bfd51cb31b8d0bc7d93d38f97ab771267d8786ab87073e0c2b8f9ddc44b274e",
				},
			},
			MintAmount:  util.GetPointer[U64](6379),
			MintAssetId: &AssetId{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
			TxPointer:   util.GetPointer[TxPointer]("0094e7560002"),
			IsScript:    false,
			IsCreate:    false,
			IsMint:      true,
			Outputs:     []Output{},
			OutputContract: &ContractOutput{
				InputIndex:  0,
				BalanceRoot: Bytes32{Hash: common.HexToHash("0xa24223e950599e9577bd2782388411d5f92c35793dc72fc25355041bb9e2f869")},
				StateRoot:   Bytes32{Hash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
			},
			Status: &TransactionStatus{
				TypeName_:     "SuccessStatus",
				SuccessStatus: &SuccessStatus{},
			},
			RawPayload: HexString{Bytes: common.FromHex("0x0000000000000002000000000094e7560000000000000002ae426ee0c79cac25ffe515ca4148e27086669bee7043b23cd380dce443213eff00000000000000002d19f8c34395032b25ae83bf88cd618a8598c6f2459f137c931879b323a41e040000000000000000000000000000000000000000000000000000000000000000000000000094e755000000000000000677777777777777777777777777777777777777777777777777777777777777770000000000000000a24223e950599e9577bd2782388411d5f92c35793dc72fc25355041bb9e2f869000000000000000000000000000000000000000000000000000000000000000000000000000018eb0000000000000000000000000000000000000000000000000000000000000000")},
		}},
	}, block)
}
