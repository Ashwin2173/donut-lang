package utils

import (
	"encoding/binary"
	"math"
)

type Constant interface{ isConstant() }

type ConstantPool struct {
	bc           ByteCode
	constantPool map[Constant]uint16
}

func NewConstantPool() *ConstantPool {
	return &ConstantPool{
		bc:           NewByteCode(),
		constantPool: map[Constant]uint16{},
	}
}

func (c *ConstantPool) length() uint16 {
	return uint16(len(c.constantPool))
}

type ConstantUtf8 struct {
	value string
}

type ConstantClass struct {
	index uint16
}

type ConstantString struct {
	index uint16
}

type ConstantInteger struct {
	value int32
}

type ConstantFloat struct {
	value float32
}

type ConstantNameAndType struct {
	nameIndex       uint16
	descriptorIndex uint16
}

type ConstantMethodRef struct {
	classIndex       uint16
	nameAndTypeIndex uint16
}

func (c ConstantUtf8) isConstant()        {}
func (c ConstantClass) isConstant()       {}
func (c ConstantFloat) isConstant()       {}
func (c ConstantString) isConstant()      {}
func (c ConstantInteger) isConstant()     {}
func (c ConstantMethodRef) isConstant()   {}
func (c ConstantNameAndType) isConstant() {}

func (cp *ConstantPool) AddUtf8(value string) (index uint16) {
	newConstant := ConstantUtf8{
		value: value,
	}
	index, exists := cp.constantPool[newConstant]
	if !exists {
		cp.bc.writeU1(0x1)                // type
		cp.bc.writeU2(uint16(len(value))) // length
		cp.bc.dumpRaw([]byte(value))
		index = uint16(len(cp.constantPool)) + 1
		cp.constantPool[newConstant] = index
	}
	return
}

func (cp *ConstantPool) AddString(value string) (index uint16) {
	stringIndex := cp.AddUtf8(value)
	newConstant := ConstantString{
		index: stringIndex,
	}
	index, exists := cp.constantPool[newConstant]
	if !exists {
		cp.bc.writeU1(0x8)
		cp.bc.writeU2(stringIndex)
		index = uint16(len(cp.constantPool)) + 1
		cp.constantPool[newConstant] = index
	}
	return
}

func (cp *ConstantPool) AddInteger(value int32) (index uint16) {
	newConstant := ConstantInteger{
		value: value,
	}
	index, exists := cp.constantPool[newConstant]
	if !exists {
		cp.bc.writeU1(0x3)
		cp.bc.writeU4(uint32(value))
		index = uint16(len(cp.constantPool)) + 1
		cp.constantPool[newConstant] = index
	}
	return
}

func (cp *ConstantPool) AddFloat(value float32) (index uint16) {
	newConstant := ConstantFloat{
		value: value,
	}
	index, exists := cp.constantPool[newConstant]
	if !exists {
		data := make([]byte, 4)
		binary.BigEndian.PutUint32(data, math.Float32bits(value))
		cp.bc.writeU1(0x4)
		cp.bc.writeU4(binary.BigEndian.Uint32(data))
		index = uint16(len(cp.constantPool)) + 1
		cp.constantPool[newConstant] = index
	}
	return
}

func (cp *ConstantPool) AddClass(name string) (index uint16) {
	nameIndex := cp.AddUtf8(name)
	newConstant := ConstantClass{
		index: nameIndex,
	}
	index, exists := cp.constantPool[newConstant]
	if !exists {
		cp.bc.writeU1(0x7)       //type
		cp.bc.writeU2(nameIndex) // nameIndex
		index = uint16(len(cp.constantPool)) + 1
		cp.constantPool[newConstant] = index
	}
	return
}

func (cp *ConstantPool) AddNameAndType(name, descriptor string) (index uint16) {
	nameIndex := cp.AddUtf8(name)
	descriptorIndex := cp.AddUtf8(descriptor)
	newConstant := ConstantNameAndType{
		nameIndex:       nameIndex,
		descriptorIndex: descriptorIndex,
	}
	index, exists := cp.constantPool[newConstant]
	if !exists {
		cp.bc.writeU1(0xC)
		cp.bc.writeU2(nameIndex)
		cp.bc.writeU2(descriptorIndex)
		index = uint16(len(cp.constantPool)) + 1
		cp.constantPool[newConstant] = index
	}
	return
}

func (cp *ConstantPool) AddMethodRef(class, name, descriptor string) (index uint16) {
	classIndex := cp.AddClass(class)
	nameAndTypeIndex := cp.AddNameAndType(name, descriptor)
	newConstant := ConstantMethodRef{
		classIndex:       classIndex,
		nameAndTypeIndex: nameAndTypeIndex,
	}
	index, exists := cp.constantPool[newConstant]
	if !exists {
		cp.bc.writeU1(0xA)
		cp.bc.writeU2(classIndex)
		cp.bc.writeU2(nameAndTypeIndex)
		index = uint16(len(cp.constantPool)) + 1
		cp.constantPool[newConstant] = index
	}
	return
}
