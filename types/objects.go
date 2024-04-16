package types

type Block struct {
	ID           BlockID
	Header       Header
	Consensus    Consensus
	Transactions []Transaction
}

type Header struct {
	ID                  BlockID
	DaHeight            U64
	TransactionsCount   U64
	MessageReceiptCount U64
	TransactionsRoot    Bytes32
	MessageReceiptRoot  Bytes32
	Height              U32
	PrevRoot            Bytes32
	Time                Tai64Timestamp
	ApplicationHash     Bytes32
}

type Consensus struct {
	*Genesis
	*PoAConsensus
}

type Genesis struct {
	ChainConfigHash Bytes32
	CoinsRoot       Bytes32
	ContractsRoot   Bytes32
	MessagesRoot    Bytes32
}

type PoAConsensus struct {
	Signature Signature
}

type Transaction struct {
	ID                   TransactionID
	InputAssetIDs        []AssetID
	InputContracts       []Contract
	InputContract        *InputContract
	Policies             *Policies
	GasPrice             *U64
	ScriptGasLimit       *U64
	Maturity             *U32
	MintAmount           *U64
	MintAssetID          *AssetID
	TxPointer            *TxPointer
	IsScript             bool
	IsCreate             bool
	IsMint               bool
	Inputs               []Input
	Outputs              []Output
	OutputContract       *ContractOutput
	Witnesses            []HexString
	ReceiptsRoot         *Bytes32
	Status               TransactionStatus
	Receipts             []Receipt
	Script               *HexString
	ScriptData           *HexString
	BytecodeWitnessIndex *Int
	BytecodeLength       *U64
	Salt                 *Salt
	StorageSlots         []HexString

	RawPayload HexString
}

type Contract struct {
	ID       ContractID
	Bytecode HexString
	Salt     Salt
}

type Policies struct {
	GasPrice     U64
	WitnessLimit U64
	Maturity     U32
	MaxFee       U64
}

type Input struct {
	*InputCoin
	*InputContract
	*InputMessage
}

type InputCoin struct {
	UtxoId           UtxoID
	Owner            Address
	Amount           U64
	AssetID          AssetID
	TxPointer        TxPointer
	WitnessIndex     Int
	Maturity         U32
	PredicateGasUsed U64
	Predicate        HexString
	PredicateData    HexString
}

type InputContract struct {
	UtxoId      UtxoID
	BalanceRoot Bytes32
	StateRoot   Bytes32
	TxPointer   TxPointer
	Contract    Contract
}

type InputMessage struct {
	Sender           Address
	Recipient        Address
	Amount           U64
	Nonce            Nonce
	WitnessIndex     Int
	PredicateGasUsed U64
	Data             HexString
	Predicate        HexString
	PredicateData    HexString
}

type Output struct {
	*CoinOutput
	*ContractOutput
	*ChangeOutput
	*VariableOutput
	*ContractCreated
}

type CoinOutput struct {
	To      Address
	Amount  U64
	AssetID AssetID
}

type ContractOutput struct {
	InputIndex  Int
	BalanceRoot Bytes32
	StateRoot   Bytes32
}

type ChangeOutput struct {
	To      Address
	Amount  U64
	AssetID AssetID
}

type VariableOutput struct {
	To      Address
	Amount  U64
	AssetID AssetID
}

type ContractCreated struct {
	Contract  Contract
	StateRoot Bytes32
}

type TransactionStatus struct {
	*SubmittedStatus
	*SuccessStatus
	*SqueezedOutStatus
	*FailureStatus
}

type SubmittedStatus struct {
	Time Tai64Timestamp
}

type SuccessStatus struct {
	TransactionID TransactionID
	Block         Block
	Time          Tai64Timestamp
	ProgramState  ProgramState
	Receipts      []Receipt
}

type SqueezedOutStatus struct {
	Reason string
}

type FailureStatus struct {
	TransactionID TransactionID
	Block         Block
	Time          Tai64Timestamp
	Reason        string
	ProgramState  ProgramState
	Receipts      []Receipt
}

type ProgramState struct {
	ReturnType ReturnType
	Data       HexString
}

type Receipt struct {
	Contract    *Contract
	Pc          *U64
	Is          *U64
	To          *Contract
	ToAddress   *Address
	Amount      *U64
	AssetID     *AssetID
	Gas         *U64
	Param1      *U64
	Param2      *U64
	Val         *U64
	Ptr         *U64
	Digest      *Bytes32
	Reason      *U64
	Ra          *U64
	Rb          *U64
	Rc          *U64
	Rd          *U64
	Len         *U64
	ReceiptType ReceiptType
	Result      *U64
	GasUsed     *U64
	Data        *HexString
	Sender      *Address
	Recipient   *Address
	Nonce       *Nonce
	ContractID  *ContractID
	SubID       *Bytes32
}
