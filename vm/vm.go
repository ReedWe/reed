// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package vm

import (
	"bytes"
	"github.com/reed/crypto"
	"github.com/reed/errors"
	"github.com/reed/vm/vmcommon"
)

var (
	vmErr = errors.New("virtualMachine run error")
)

type signFunc func(pk []byte, sig []byte) bool

type VM struct {
	script []byte
	stack  [][]byte
	signTx signFunc
}

func NewVirtualMachine(scriptSig []byte, scriptPK []byte, signTx signFunc) *VM {
	return &VM{
		script: bytes.Join([][]byte{
			scriptSig, scriptPK,
		}, []byte{}),
		signTx: signTx,
	}
}

func (v *VM) Run() error {
	scriptLen := len(v.script)

	push := func(data []byte) {
		v.stack = append(v.stack, data)
	}

	pop := func() []byte {
		top := v.stack[len(v.stack)-1]
		v.stack = v.stack[:len(v.stack)-1]
		return top
	}

	pointer := 0
	for {
		if pointer >= scriptLen {
			break
		}
		op := v.script[pointer : pointer+1]
		pointer++
		switch {
		case bytes.Equal(op, []byte{byte(vmcommon.OpPushData64)}):
			push(v.script[pointer : pointer+64])
			pointer += 64 - 1
		case bytes.Equal(op, []byte{byte(vmcommon.OpPushData32)}):
			push(v.script[pointer : pointer+32])
			pointer += 32 - 1
		case bytes.Equal(op, []byte{byte(vmcommon.OpDup)}):
			d := v.stack[len(v.stack)-1]
			v.stack = append(v.stack, d)
		case bytes.Equal(op, []byte{byte(vmcommon.OpHash256)}):
			push(crypto.Sha256(pop()))
		case bytes.Equal(op, []byte{byte(vmcommon.OpEqualVerify)}):
			a := pop()
			b := pop()
			if !bytes.Equal(a, b) {
				return errors.Wrap(vmErr, "OP_EQUAL_VERIFY failed")
			}
		case bytes.Equal(op, []byte{byte(vmcommon.OpCheckSig)}):
			if ok := v.signTx(pop(), pop()); !ok {
				return errors.Wrap(vmErr, "OP_CHECK_SIG signature failed")
			}
		}
	}
	return nil

}
