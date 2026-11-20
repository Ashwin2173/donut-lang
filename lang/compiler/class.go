package compiler

const ( // todo: move this to a new file
	sourceFile = "SourceFile"
	code       = "Code"
)

const (
	AccPublic     = 0x0001
	AccFinal      = 0x0010
	AccSuper      = 0x0020
	AccInterface  = 0x0200
	AccAbstract   = 0x0400
	AccSynthetic  = 0x1000
	AccAnnotation = 0x2000
	AccEnum       = 0x4000
)

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
