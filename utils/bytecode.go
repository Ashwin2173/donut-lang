package utils

import (
	"encoding/binary"
)

type ByteCode struct {
	byteCode []byte
}

func NewByteCode() ByteCode {
	return ByteCode{
		byteCode: []byte{},
	}
}

func (b *ByteCode) length() int {
	return len(b.byteCode)
}

func (b *ByteCode) writeU1(value byte) {
	b.byteCode = append(b.byteCode, value)
}

func (b *ByteCode) writeU2(value uint16) {
	b.byteCode = binary.BigEndian.AppendUint16(b.byteCode, value)
}

func (b *ByteCode) writeU4(value uint32) {
	b.byteCode = binary.BigEndian.AppendUint32(b.byteCode, value)
}

func (b *ByteCode) dump(byteCode ByteCode) {
	b.byteCode = append(b.byteCode, byteCode.byteCode...)
}

func (b *ByteCode) dumpRaw(bytes []byte) {
	b.byteCode = append(b.byteCode, bytes...)
}

// opcode functions

func (b *ByteCode) opReturn() {
	b.writeU1(OP_RETURN)
}

func (b *ByteCode) opLoadConstant(index uint16) {
	if index <= 255 {
		b.writeU1(OP_LOAD_CONSTANT)
		b.writeU1(uint8(index))
	} else {
		b.writeU1(OP_LOAD_CONSTANT_W)
		b.writeU2(index)
	}
}

func (b *ByteCode) opInvokeStatic(index uint16) {
	b.writeU1(OP_INVOKE_STATIC)
	b.writeU2(index)
}

func (b *ByteCode) opIAdd() {
	b.writeU1(OP_IADD)
}

func (b *ByteCode) opISub() {
	b.writeU1(OP_ISUB)
}

func (b *ByteCode) opFAdd() {
	b.writeU1(OP_FADD)
}

func (b *ByteCode) opFSub() {
	b.writeU1(OP_FSUB)
}

func (b *ByteCode) opI2F() {
	b.writeU1(OP_I2F)
}
