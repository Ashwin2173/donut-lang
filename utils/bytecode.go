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

func (b *ByteCode) writeJumpPlaceHolder(opcode byte) (pos int) {
	b.writeU1(opcode)
	pos = len(b.byteCode)
	b.writeU1(0)
	b.writeU1(0)
	return
}

func (b *ByteCode) patchJump(pos int) {
	offset := int16(len(b.byteCode) - (pos - 1))
	binary.BigEndian.PutUint16(b.byteCode[pos:], uint16(offset))
}

// opcode functions

func (b *ByteCode) opReturn() {
	b.writeU1(OP_RETURN)
}

func (b *ByteCode) opBipush(value byte) {
	b.writeU1(OP_BIPUSH)
	b.writeU1(value)
}

func (b *ByteCode) opIstore(index uint16) {
	switch {
	case index == 0:
		b.writeU1(OP_ISTORE_0)
	case index == 1:
		b.writeU1(OP_ISTORE_1)
	case index == 2:
		b.writeU1(OP_ISTORE_2)
	case index == 3:
		b.writeU1(OP_ISTORE_3)
	case index >= 4 && index <= 255:
		b.writeU1(OP_ISTORE)
		b.writeU1(byte(index))
	default:
		b.writeU1(OP_WIDE)
		b.writeU1(OP_ISTORE)
		b.writeU2(index)
	}
}

func (b *ByteCode) opIload(index uint16) {
	switch {
	case index == 0:
		b.writeU1(OP_ILOAD_0)
	case index == 1:
		b.writeU1(OP_ILOAD_1)
	case index == 2:
		b.writeU1(OP_ILOAD_2)
	case index == 3:
		b.writeU1(OP_ILOAD_3)
	case index >= 4 && index <= 255:
		b.writeU1(OP_ILOAD)
		b.writeU1(byte(index))
	default:
		b.writeU1(OP_WIDE)
		b.writeU1(OP_ILOAD)
		b.writeU2(index)
	}
}

func (b *ByteCode) opFstore(index uint16) {
	switch {
	case index == 0:
		b.writeU1(OP_FSTORE_0)
	case index == 1:
		b.writeU1(OP_FSTORE_1)
	case index == 2:
		b.writeU1(OP_FSTORE_2)
	case index == 3:
		b.writeU1(OP_FSTORE_3)
	case index >= 4 && index <= 255:
		b.writeU1(OP_FSTORE)
		b.writeU1(byte(index))
	default:
		b.writeU1(OP_WIDE)
		b.writeU1(OP_FSTORE)
		b.writeU2(index)
	}
}

func (b *ByteCode) opFload(index uint16) {
	switch {
	case index == 0:
		b.writeU1(OP_FLOAD_0)
	case index == 1:
		b.writeU1(OP_FLOAD_1)
	case index == 2:
		b.writeU1(OP_FLOAD_2)
	case index == 3:
		b.writeU1(OP_FLOAD_3)
	case index >= 4 && index <= 255:
		b.writeU1(OP_FLOAD)
		b.writeU1(byte(index))
	default:
		b.writeU1(OP_WIDE)
		b.writeU1(OP_FLOAD)
		b.writeU2(index)
	}
}

func (b *ByteCode) opAstore(index uint16) {
	switch {
	case index == 0:
		b.writeU1(OP_ASTORE_0)
	case index == 1:
		b.writeU1(OP_ASTORE_1)
	case index == 2:
		b.writeU1(OP_ASTORE_2)
	case index == 3:
		b.writeU1(OP_ASTORE_3)
	case index >= 4 && index <= 255:
		b.writeU1(OP_ASTORE)
		b.writeU1(byte(index))
	default:
		b.writeU1(OP_WIDE)
		b.writeU1(OP_ASTORE)
		b.writeU2(index)
	}
}

func (b *ByteCode) opAload(index uint16) {
	switch {
	case index == 0:
		b.writeU1(OP_ALOAD_0)
	case index == 1:
		b.writeU1(OP_ALOAD_1)
	case index == 2:
		b.writeU1(OP_ALOAD_2)
	case index == 3:
		b.writeU1(OP_ALOAD_3)
	case index >= 4 && index <= 255:
		b.writeU1(OP_ALOAD)
		b.writeU1(byte(index))
	default:
		b.writeU1(OP_WIDE)
		b.writeU1(OP_ALOAD)
		b.writeU2(index)
	}
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

func (b *ByteCode) opInvokeDynamic(index uint16) {
	b.writeU1(OP_INVOKE_DYNAMIC)
	b.writeU2(index)
	b.writeU1(0)
	b.writeU1(0)
}

func (b *ByteCode) opIAdd() {
	b.writeU1(OP_IADD)
}

func (b *ByteCode) opISub() {
	b.writeU1(OP_ISUB)
}

func (b *ByteCode) opIMul() {
	b.writeU1(OP_IMUL)
}

func (b *ByteCode) opIDiv() {
	b.writeU1(OP_IDIV)
}

func (b *ByteCode) opIRem() {
	b.writeU1(OP_IREM)
}

func (b *ByteCode) opFAdd() {
	b.writeU1(OP_FADD)
}

func (b *ByteCode) opFSub() {
	b.writeU1(OP_FSUB)
}

func (b *ByteCode) opFMul() {
	b.writeU1(OP_FMUL)
}

func (b *ByteCode) opFDiv() {
	b.writeU1(OP_FDIV)
}

func (b *ByteCode) opFRem() {
	b.writeU1(OP_FREM)
}

func (b *ByteCode) opI2F() {
	b.writeU1(OP_I2F)
}

func (b *ByteCode) opIconst0() {
	b.writeU1(OP_ICONST_0)
}

func (b *ByteCode) opIconst1() {
	b.writeU1(OP_ICONST_1)
}

func (b *ByteCode) opFcmpl() {
	b.writeU1(OP_FCMPL)
}

func (b *ByteCode) opIAnd() {
	b.writeU1(OP_IAND)
}

func (b *ByteCode) opIOr() {
	b.writeU1(OP_IOR)
}

func (b *ByteCode) opIXor() {
	b.writeU1(OP_IXOR)
}

func (b *ByteCode) opInvokeVirtual(index uint16) {
	b.writeU1(OP_INVOKE_VIRTUAL)
	b.writeU2(index)
}
