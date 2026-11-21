package compiler

type Method struct {
	AccessFlag uint16
	Name       string
	Descriptor string
	Code       Code
	Attributes []string

	nameIndex       uint16
	descriptorIndex uint16
}

func NewMethod(context *Context, method Method) *Method {
	method.nameIndex = context.cp.AddUtf8(method.Name)
	method.descriptorIndex = context.cp.AddUtf8(method.Descriptor)
	return &method
}

func (m Method) GetByteCode(context *Context) (bc ByteCode) {
	bc.writeU2(m.AccessFlag)
	bc.writeU2(m.nameIndex)
	bc.writeU2(m.descriptorIndex)
	attrIndex := bc.beginAttr()
	attrCount := 0
	// Code
	codeDump := m.Code.GetByteCode(context)
	bc.writeU24(context.cp.AddUtf8(code), uint32(codeDump.length()))
	bc.dump(codeDump)
	attrCount++
	bc.endAttr(attrIndex, attrCount)
	return
}
