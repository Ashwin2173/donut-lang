package compiler

type Class struct {
	AccessFlag uint16
	ThisClass  string
	SuperClass string
	Interfaces []string
	Fields     []string
	Methods    []Method
	Attributes string

	thisClassIndex  uint16
	superClassIndex uint16
}

func NewClass(context *Context, class Class) *Class {
	class.thisClassIndex = context.cp.AddClass(class.ThisClass)
	class.superClassIndex = context.cp.AddClass(class.SuperClass)
	return &class
}

func (c Class) GetByteCode(context *Context) (bc ByteCode) {
	bc.writeU2(c.AccessFlag)
	bc.writeU2(c.thisClassIndex)
	bc.writeU2(c.superClassIndex)
	bc.writeU2(uint16(len(c.Interfaces)))
	// add GetByteCode for Interface
	bc.writeU2(uint16(len(c.Fields)))
	// add GetByteCode for Fields
	bc.writeU2(uint16(len(c.Methods)))
	for _, m := range c.Methods {
		bc.dump(m.GetByteCode(context))
	}
	attrIndex := bc.beginAttr()
	attrCount := 0
	// SourceFile Attr
	bc.writeU24(context.cp.AddUtf8(sourceFile), 2)
	bc.writeU2(context.cp.AddUtf8(context.FileName))
	attrCount++
	bc.endAttr(attrIndex, attrCount)
	return
}
