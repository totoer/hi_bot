package executor

import (
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

func findDoc(L *lua.LState) int {
	return 1
}

func saveDoc(L *lua.LState) int {
	return 1
}

func deleteDoc(L *lua.LState) int {
	return 1
}

func ExecHandler(author, message, script string) (string, bool) {
	L := lua.NewState()
	defer L.Close()

	lua.OpenBase(L)

	L.SetGlobal("findDoc", L.NewFunction(findDoc))
	L.SetGlobal("saveDoc", L.NewFunction(saveDoc))
	L.SetGlobal("deleteDoc", L.NewFunction(deleteDoc))

	if err := L.DoFile(filepath.Join("./bots", script)); err != nil {
		panic(err)
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("main"), // name of Lua function
		NRet:    1,                   // number of returned values
		Protect: true,                // return err or panic
	}, lua.LString(author), lua.LString(message)); err != nil {
		panic(err)
	}

	if result, ok := L.Get(-1).(lua.LString); ok {
		return string(result), true
	}

	return "", false
}
