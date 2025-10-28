package state

func (self *luaState) PC() int {
	return self.pc
}

func (self *luaState) AddPC(n int) {
	self.pc += n
}

func (self *luaState) Fetch() uint32 {
	i := self.proto.Code[self.pc]
	self.pc++
	return i
}

func (self *luaState) GetConst(idx int) {
	c := self.proto.Constants[idx]
	self.stack.push(c)
}

func (self *luaState) GetRK(rk int) { //rk is OpArgK
	if rk > 0xFF { //constant ,first bit is 1
		self.GetConst(rk & 0xFF)
	} else { //register
		self.PushValue(rk + 1) //Lua Api'stack index begin in 1
	}
}
