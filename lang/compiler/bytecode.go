package compiler

import (
	"encoding/binary"
)

const (
	op_iconst_0        = 0x03
	op_iconst_1        = 0x04
	op_load_constant   = 0x12
	op_load_constant_w = 0x13
	op_iand            = 0x7e
	op_ior             = 0x80
	op_ixor            = 0x82
	op_i2f             = 0x86
	op_fcmpl           = 0x95
	op_ifeq            = 0x99
	op_ifne            = 0x9a
	op_iflt            = 0x9b
	op_ifge            = 0x9c
	op_ifgt            = 0x9d
	op_ifle            = 0x9e
	op_if_icmpeq       = 0x9f
	op_if_icmpne       = 0xa0
	op_if_icmplt       = 0xa1
	op_if_icmpge       = 0xa2
	op_if_icmpgt       = 0xa3
	op_if_icmple       = 0xa4
	op_goto            = 0xa7
	op_return          = 0xb1
	op_invoke_virtual  = 0xb6
	op_invoke_static   = 0xb8
	op_invoke_dynamic  = 0xba
	op_iadd            = 0x60
	op_isub            = 0x64
	op_imul            = 0x68
	op_idiv            = 0x6c
	op_irem            = 0x70
	op_fadd            = 0x62
	op_fsub            = 0x66
	op_fmul            = 0x6a
	op_fdiv            = 0x6e
	op_frem            = 0x72
	op_fneg            = 0x76
	op_ineg            = 0x74
	op_bipush          = 0x10
	op_istore_0        = 0x3b
	op_istore_1        = 0x3c
	op_istore_2        = 0x3d
	op_istore_3        = 0x3e
	op_istore          = 0x36
	op_iload_0         = 0x1a
	op_iload_1         = 0x1b
	op_iload_2         = 0x1c
	op_iload_3         = 0x1d
	op_iload           = 0x15
	op_fstore_0        = 0x43
	op_fstore_1        = 0x44
	op_fstore_2        = 0x45
	op_fstore_3        = 0x46
	op_fstore          = 0x38
	op_fload_0         = 0x22
	op_fload_1         = 0x23
	op_fload_2         = 0x24
	op_fload_3         = 0x25
	op_fload           = 0x17
	op_astore_0        = 0x4b
	op_astore_1        = 0x4c
	op_astore_2        = 0x4d
	op_astore_3        = 0x4e
	op_astore          = 0x3a
	op_aload_0         = 0x2a
	op_aload_1         = 0x2b
	op_aload_2         = 0x2c
	op_aload_3         = 0x2d
	op_aload           = 0x19
	op_wide            = 0xc4
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

func (b *ByteCode) writeU24(v1 uint16, v2 uint32) {
	b.writeU2(v1)
	b.writeU4(v2)
}

func (b *ByteCode) dump(byteCode ByteCode) {
	b.byteCode = append(b.byteCode, byteCode.byteCode...)
}

func (b *ByteCode) dumpRaw(bytes []byte) {
	b.byteCode = append(b.byteCode, bytes...)
}

// helper functions

func (b *ByteCode) beginAttr() (pos int) {
	pos = len(b.byteCode)
	b.writeU2(0)
	return
}

func (b *ByteCode) endAttr(pos, count int) {
	binary.BigEndian.PutUint16(b.byteCode[pos:], uint16(count))
}

func (b *ByteCode) writeJumpPlaceholder(opcode byte) (pos int) {
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

func (b *ByteCode) patchJumpTo(pos int, target int) {
	offset := int16(target - pos + 1)
	binary.BigEndian.PutUint16(b.byteCode[pos:], uint16(offset))
}

// opcode helper functions

func (b *ByteCode) opBipush(value byte) {
	b.writeU1(op_bipush)
	b.writeU1(value)
}

func (b *ByteCode) opIstore(index uint16) {
	switch {
	case index == 0:
		b.writeU1(op_istore_0)
	case index == 1:
		b.writeU1(op_istore_1)
	case index == 2:
		b.writeU1(op_istore_2)
	case index == 3:
		b.writeU1(op_istore_3)
	case index >= 4 && index <= 255:
		b.writeU1(op_istore)
		b.writeU1(byte(index))
	default:
		b.writeU1(op_wide)
		b.writeU1(op_istore)
		b.writeU2(index)
	}
}

func (b *ByteCode) opIload(index uint16) {
	switch {
	case index == 0:
		b.writeU1(op_iload_0)
	case index == 1:
		b.writeU1(op_iload_1)
	case index == 2:
		b.writeU1(op_iload_2)
	case index == 3:
		b.writeU1(op_iload_3)
	case index >= 4 && index <= 255:
		b.writeU1(op_iload)
		b.writeU1(byte(index))
	default:
		b.writeU1(op_wide)
		b.writeU1(op_iload)
		b.writeU2(index)
	}
}

func (b *ByteCode) opFstore(index uint16) {
	switch {
	case index == 0:
		b.writeU1(op_fstore_0)
	case index == 1:
		b.writeU1(op_fstore_1)
	case index == 2:
		b.writeU1(op_fstore_2)
	case index == 3:
		b.writeU1(op_fstore_3)
	case index >= 4 && index <= 255:
		b.writeU1(op_fstore)
		b.writeU1(byte(index))
	default:
		b.writeU1(op_wide)
		b.writeU1(op_fstore)
		b.writeU2(index)
	}
}

func (b *ByteCode) opFload(index uint16) {
	switch {
	case index == 0:
		b.writeU1(op_fload_0)
	case index == 1:
		b.writeU1(op_fload_1)
	case index == 2:
		b.writeU1(op_fload_2)
	case index == 3:
		b.writeU1(op_fload_3)
	case index >= 4 && index <= 255:
		b.writeU1(op_fload)
		b.writeU1(byte(index))
	default:
		b.writeU1(op_wide)
		b.writeU1(op_fload)
		b.writeU2(index)
	}
}

func (b *ByteCode) opAstore(index uint16) {
	switch {
	case index == 0:
		b.writeU1(op_astore_0)
	case index == 1:
		b.writeU1(op_astore_1)
	case index == 2:
		b.writeU1(op_astore_2)
	case index == 3:
		b.writeU1(op_astore_3)
	case index >= 4 && index <= 255:
		b.writeU1(op_astore)
		b.writeU1(byte(index))
	default:
		b.writeU1(op_wide)
		b.writeU1(op_astore)
		b.writeU2(index)
	}
}

func (b *ByteCode) opAload(index uint16) {
	switch {
	case index == 0:
		b.writeU1(op_aload_0)
	case index == 1:
		b.writeU1(op_aload_1)
	case index == 2:
		b.writeU1(op_aload_2)
	case index == 3:
		b.writeU1(op_aload_3)
	case index >= 4 && index <= 255:
		b.writeU1(op_aload)
		b.writeU1(byte(index))
	default:
		b.writeU1(op_wide)
		b.writeU1(op_aload)
		b.writeU2(index)
	}
}

func (b *ByteCode) opLoadConstant(index uint16) {
	if index <= 255 {
		b.writeU1(op_load_constant)
		b.writeU1(uint8(index))
	} else {
		b.writeU1(op_load_constant_w)
		b.writeU2(index)
	}
}

func (b *ByteCode) opInvokeStatic(index uint16) {
	b.writeU1(op_invoke_static)
	b.writeU2(index)
}

func (b *ByteCode) opInvokeDynamic(index uint16) {
	b.writeU1(op_invoke_dynamic)
	b.writeU2(index)
	b.writeU1(0)
	b.writeU1(0)
}

func (b *ByteCode) opInvokeVirtual(index uint16) {
	b.writeU1(op_invoke_virtual)
	b.writeU2(index)
}
