package fuel

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_jsonMarshalScalar(t *testing.T) {
	type typ struct {
		A1 HexString
		A2 *HexString
		A3 *HexString
		B1 Boolean
		B2 *Boolean
		B3 *Boolean
		C1 Address
		C2 *Address
		C3 *Address
		D1 U8
		D2 *U8
		D3 *U8
		E1 String
		E2 *String
		E3 *String
	}

	{
		var bv Boolean
		var u8v U8 = 101
		var sv String = "def"
		a := typ{
			A1: HexString{Bytes: []byte{1, 2, 3, 4}},
			A2: &HexString{Bytes: []byte{2, 3, 4}},
			B1: true,
			B2: &bv,
			C1: Address{Hash: [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}},
			C2: &Address{Hash: [32]byte{2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3}},
			D1: 100,
			D2: &u8v,
			E1: "abc",
			E2: &sv,
		}
		b, err := json.Marshal(&a)
		assert.NoError(t, err)
		assert.Equal(t, `{"A1":"0x01020304","A2":"0x020304","A3":null,"B1":true,"B2":false,"B3":null,"C1":"0x0102030405060708090001020304050607080900010203040506070809000102","C2":"0x0203040506070809000102030405060708090001020304050607080900010203","C3":null,"D1":100,"D2":101,"D3":null,"E1":"abc","E2":"def","E3":null}`, string(b))
	}

	{
		var a typ
		err := json.Unmarshal([]byte(`{"A1":"0x01020304","A2":"0x030201","B3":true,"C1":"0x0304050607080900010203040506070809000102030405060708090001020304","C2":"0x0405060708090001020304050607080900010203040506070809000102030405","D1":50,"D2":51,"E1":"xxx","E2":"yyy"}`), &a)
		assert.NoError(t, err)
		assert.Equal(t, []byte{1, 2, 3, 4}, []byte(a.A1.Bytes))
		assert.Equal(t, []byte{3, 2, 1}, []byte(a.A2.Bytes))
		assert.Nil(t, a.A3)
		assert.False(t, bool(a.B1))
		assert.Nil(t, a.B2)
		assert.True(t, bool(*a.B3))
		assert.Equal(t, [32]byte{3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4}, [32]byte(a.C1.Hash))
		assert.Equal(t, [32]byte{4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5}, [32]byte(a.C2.Hash))
		assert.Nil(t, a.C3)
		assert.Equal(t, uint8(50), uint8(a.D1))
		assert.Equal(t, uint8(51), uint8(*a.D2))
		assert.Nil(t, a.D3)
		assert.Equal(t, "xxx", string(a.E1))
		assert.Equal(t, "yyy", string(*a.E2))
		assert.Nil(t, a.E3)
	}

}

func Test_jsonMarshalEnum(t *testing.T) {
	type typ struct {
		RT ReceiptType
	}
	{
		a := typ{
			RT: "CALL",
		}
		b, err := json.Marshal(a)
		assert.NoError(t, err)
		assert.Equal(t, `{"RT":"CALL"}`, string(b))
	}
	{
		var a typ
		err := json.Unmarshal([]byte(`{"RT":"CALL"}`), &a)
		assert.NoError(t, err)
		assert.Equal(t, "CALL", string(a.RT))
	}
	{
		var a typ
		err := json.Unmarshal([]byte(`{"RT":"CAL"}`), &a)
		assert.EqualError(t, err, "invalid value \"CAL\" for enum type ReceiptType")
	}
}
