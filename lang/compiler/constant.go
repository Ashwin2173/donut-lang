package compiler

const (
	sourceFile = "SourceFile"
	code       = "Code"
)

const (
	AccPublic       = 0x0001 // class, method
	AccFinal        = 0x0010 // class, method
	AccSuper        = 0x0020 // class
	AccInterface    = 0x0200 // class
	AccAbstract     = 0x0400 // class, method
	AccSynthetic    = 0x1000 // class, method
	AccAnnotation   = 0x2000 // class
	AccEnum         = 0x4000 // class
	AccPrivate      = 0x0002 // method
	AccProtected    = 0x0004 // method
	AccStatic       = 0x0008 // method
	AccSynchronized = 0x0020 // method
	AccBridge       = 0x0040 // method
	AccVarargs      = 0x0080 // method
	AccNative       = 0x0100 // method
	AccStrict       = 0x0800 // method
)
