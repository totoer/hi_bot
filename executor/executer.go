package executor

import (
	"path/filepath"

	"github.com/spf13/viper"
	lua "github.com/yuin/gopher-lua"
)

func ExecHandler(author, message, script string) (string, bool) {
	L := lua.NewState()
	defer L.Close()

	lua.OpenBase(L)

	L.SetGlobal("httpGet", L.NewFunction(luaHTTPGet))
	L.SetGlobal("jsonDecode", L.NewFunction(luaJSONDecode))
	L.SetGlobal("jsonEncode", L.NewFunction(luaJSONEncode))
	L.SetGlobal("kvSet", L.NewFunction(luaKVSet))
	L.SetGlobal("kvGet", L.NewFunction(luaKVGet))
	L.SetGlobal("kvDel", L.NewFunction(luaKVDel))

	botFilepath := filepath.Join(viper.GetString("bots_path"), script)
	if err := L.DoFile(botFilepath); err != nil {
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
