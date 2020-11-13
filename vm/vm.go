package vm

import (
	"bytes"
	"github.com/tybc/crypto"
	"github.com/tybc/errors"
	"github.com/tybc/vm/vmcommon"
)

var (
	vmErr = errors.New("virualMachine run error")
)

type SignFunc func(pk []byte, sig []byte) bool

type VM struct {
	Script []byte
	Stack  [][]byte
	SignTx SignFunc
}

func NewVirtualMachine(scriptSig []byte, scriptPK []byte, signTx SignFunc) *VM {
	return &VM{
		Script: bytes.Join([][]byte{
			scriptSig, scriptPK,
		}, []byte{}),
		SignTx: signTx,
	}
}

func (v *VM) Run() error {
	scriptLen := len(v.Script)

	push := func(data []byte) {
		v.Stack = append(v.Stack, data)
	}

	pop := func() []byte {
		top := v.Stack[len(v.Stack)-1]
		v.Stack = v.Stack[:len(v.Stack)-1]
		return top
	}

	pointer := 0
	for {
		if pointer >= scriptLen {
			break
		}
		op := v.Script[pointer : pointer+1]
		pointer++
		switch {
		case bytes.Equal(op, []byte{byte(vmcommon.OpPushData64)}):
			push(v.Script[pointer : pointer+64])
			pointer += 64 - 1
		case bytes.Equal(op, []byte{byte(vmcommon.OpPushData32)}):
			push(v.Script[pointer : pointer+32])
			pointer += 32 - 1
		case bytes.Equal(op, []byte{byte(vmcommon.OpDup)}):
			d := v.Stack[len(v.Stack)-1]
			v.Stack = append(v.Stack, d)
		case bytes.Equal(op, []byte{byte(vmcommon.OpHash256)}):
			push(crypto.Sha256(pop()))
		case bytes.Equal(op, []byte{byte(vmcommon.OpEqualVerify)}):
			a := pop()
			b := pop()
			if !bytes.Equal(a, b) {
				return errors.Wrap(vmErr, "OP_EQUAL_VERIFY failed")
			}
		case bytes.Equal(op, []byte{byte(vmcommon.OpCheckSig)}):
			if ok := v.SignTx(pop(), pop()); !ok {
				return errors.Wrap(vmErr, "OP_CHECK_SIG signature failed")
			}
		}
	}
	return nil

}
