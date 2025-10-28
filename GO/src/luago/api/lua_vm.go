package api

type LuaVM interface {
	LuaState
	PC() int          //return present Pc
	AddPC(n int)      //change Pc
	Fetch() uint32    //get present instruction ,let Pc point next instruction
	GetConst(idx int) //push ordered Constant to stack top
	GetRK(rk int)     //push ordered constant or stack value to stack top
}
