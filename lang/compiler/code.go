package compiler

type Code struct {
	stackSize int16
	maxLocals int16
}

func NewCode(Code) *Code {
	return &Code
}

func (c Code) GetByteCode(context *Context) (bc ByteCode) {
	return
}