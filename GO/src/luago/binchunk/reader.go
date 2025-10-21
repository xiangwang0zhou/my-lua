package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte
}

func (self *reader) readByte() byte {
	b := self.data[0]         //read the first byte
	self.data = self.data[1:] //self.data+=byte
	return b
}

func (self *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(self.data)
	self.data = self.data[4:] //push a c_int
	return i
}

func (self *reader) readUint64() uint64 { //read a size_t
	i := binary.LittleEndian.Uint64(self.data)
	self.data = self.data[8:]
	return i
}

func (self *reader) readLuaInteger() int64 { //read a Lua integer the same as go int64
	return int64(self.readUint64())
}

func (self *reader) readLuaNumber() float64 {
	return math.Float64frombits(self.readUint64())
}

func (self *reader) readString() string { //read a string
	size := uint(self.readByte())
	if size == 0 { //Null string
		return ""
	}
	if size == 0xFF { //a long string
		size = uint(self.readUint64())
	}
	bytes := self.readBytes(size - 1)
	return string(bytes)
}

func (self *reader) readBytes(n uint) []byte { //read a char
	bytes := self.data[:n]
	self.data = self.data[n:]
	return bytes
}

func (self *reader) checkHeader() {
	if string(self.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk!")
	} else if self.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	} else if self.readByte() != LUAC_FORMAT {
		panic("corrupted!")
	} else if self.readByte() != CINT_SIZE {
		panic("int size mismatched!")
	} else if self.readByte() != CSIZET_SIZE {
		panic("size_t size mismatched!")
	} else if self.readByte() != INSTRUCTION_SIZE {
		panic("instruction siez mismatch!")
	} else if self.readByte() != LUA_INTEGER_SIZE {
		panic("lua_integer size mismatch!")
	} else if self.readByte() != LUA_NUMBER_SIZE {
		panic("lua_number size missmatch!")
	} else if self.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	} else if self.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}

func (self *reader) readProto(parentSource string) *Prototype {
	source := self.readString()
	if source == "" {
		source = parentSource
	}
	return &Prototype{
		Source:         source,
		LineDefined:    self.readUint32(),
		LastLineDfined: self.readUint32(),
		NumberParams:   self.readByte(),
		IsVararg:       self.readByte(),
		MaxStackSize:   self.readByte(),
		Code:           self.readCode(),
		Constants:      self.readConstants(),
		Upvalues:       self.readUpvalue(),
		Protos:         self.readProtos(source),
		LineIofo:       self.readLineInfo(),
		LocVars:        self.readLocVars(),
		UpvalueNames:   self.readUpvalueNames(),
	}
}

func (self *reader) readCode() []uint32 {
	code := make([]uint32, self.readUint32())
	for i := range code {
		code[i] = self.readUint32()
	}
	return code
}

func (self *reader) readConstants() []interface{} {
	constants := make([]interface{}, self.readUint32())
	for i := range constants {
		constants[i] = self.readConstant()
	}
	return constants
}

func (self *reader) readConstant() interface{} {
	switch self.readByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEN:
		return self.readByte() != 0
	case TAG_INTERGER:
		return self.readLuaInteger()
	case TAG_NUMBER:
		return self.readLuaNumber()
	case TAG_SHORT_STR:
		return self.readString()
	case TAG_LONG_STR:
		return self.readString()
	default:
		panic("corrupted!")
	}
}

func (self *reader) readUpvalue() []Upvalue {
	upvalues := make([]Upvalue, self.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: self.readByte(),
			Idx:     self.readByte(),
		}
	}
	return upvalues
}

func (self *reader) readProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, self.readUint32())
	for i := range protos {
		protos[i] = self.readProto(parentSource)
	}
	return protos
}

func (self *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, self.readUint32())
	for i := range lineInfo {
		lineInfo[i] = self.readUint32()
	}
	return lineInfo
}

func (self *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, self.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: self.readString(),
			StartPc: self.readUint32(),
			EndPc:   self.readUint32(),
		}
	}
	return locVars
}

func (self *reader) readUpvalueNames() []string {
	names := make([]string, self.readUint32())
	for i := range names {
		names[i] = self.readString()
	}
	return names
}
