package utils

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

func (c ConstantUtf8) isConstant()  {}
func (c ConstantClass) isConstant() {}

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
