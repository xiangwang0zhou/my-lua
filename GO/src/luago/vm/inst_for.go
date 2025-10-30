package vm

import . "LuaGo/src/luago/api"

func forPrep(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1

	if vm.Type(a) == LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a))
		vm.Replace(a)
	}
	if vm.Type(a+1) == LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 1))
		vm.Replace(a + 1)
	}
	if vm.Type(a+2) == LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 2))
		vm.Replace(a + 2)
	}
	//R(a)-=R(a+2)
	vm.PushValue(a)
	vm.PushValue(a + 2)
	vm.Arith(LUA_OPSUB)
	vm.Replace(a)
	//pc+=sBx
	vm.AddPC(sBx)
}

func forLoop(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1
	//R(A)+=R(A+2)
	vm.PushValue(a + 2)
	vm.PushValue(a)
	vm.Arith(LUA_OPADD)
	vm.Replace(a)

	//R(A)<?=R(A+1)
	isPositiveStep := vm.ToNumber(a+2) >= 0
	if isPositiveStep && vm.Compare(a, a+1, LUA_OPLE) ||
		!isPositiveStep && vm.Compare(a+1, a, LUA_OPLE) {
		vm.AddPC(sBx)   //pc+=sBx
		vm.Copy(a, a+3) //R(A+3)=R(A)
	}
}
