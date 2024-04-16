package types

type Enum string

type ReturnType Enum

var ReturnTypeValues = []string{
	"RETURN",
	"RETURN_DATA",
	"REVERT",
}

type ReceiptType Enum

var ReceiptTypeValues = []string{
	"CALL",
	"RETURN",
	"RETURN_DATA",
	"PANIC",
	"REVERT",
	"LOG",
	"LOG_DATA",
	"TRANSFER",
	"TRANSFER_OUT",
	"SCRIPT_RESULT",
	"MESSAGE_OUT",
	"MINT",
	"BURN",
}
