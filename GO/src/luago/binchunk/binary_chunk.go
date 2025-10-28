package binchunk

type binaryChunk struct {
	header
	sizeUpvalue byte
	mainFunc    *Prototype
}

type header struct {
	signature       [4]byte //binary file will start with a Magic Number(like .class start with 0xCAFEBABE).lua'chunk start with Magic Number (in other words Signature) is 0x1B4C7561 (is ESC,L,u,a in ASCII byte),in GOlang is "\x1bLua"
	version         byte    //in Lua 'version' have tree parts :Major version ,Minor version,Release version,in lua_chunk these version = Major version * 16 + Minor version
	format          byte    //the version for chunk file ,in lua offical is 00
	luaData         [6]byte //a const byte
	cintSize        byte    // the size of
	sizetSize       byte    // the size of
	instructionSize byte    // the size of
	luaIntegerSize  byte    // the size of
	luaNumberSize   byte    // the size of
	luacInt         int64   //the size of
	luacNum         float64 //the format of float ,In common is IEEE754
}

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

type Prototype struct { //the proto of function
	Source         string
	LineDefined    uint32
	LastLineDfined uint32 //the position
	NumberParams   byte
	IsVararg       byte
	MaxStackSize   byte
	Code           []uint32
	Constants      []interface{} //used as the union in C
	Upvalues       []Upvalue
	Protos         []*Prototype
	LineIofo       []uint32
	LocVars        []LocVar
	UpvalueNames   []string
}

const (
	TAG_NIL       = 0x00
	TAG_BOOLEN    = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTERGER  = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type Upvalue struct {
	Instack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPc uint32
	EndPc   uint32
}

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()        //check the Header
	reader.readByte()           //jump acrose the cont of Upvalue
	return reader.readProto("") //read the function proto
}
