package executor

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/viper"

	python "github.com/sbinet/go-python"
	lua "github.com/yuin/gopher-lua"
)

func execLuaScript(author, message, script string) (string, bool) {
	L := lua.NewState()
	defer L.Close()

	lua.OpenBase(L)

	L.SetGlobal("httpGet", L.NewFunction(luaHTTPGet))
	L.SetGlobal("jsonDecode", L.NewFunction(luaJSONDecode))
	L.SetGlobal("jsonEncode", L.NewFunction(luaJSONEncode))
	L.SetGlobal("kvSet", L.NewFunction(luaKVSet))
	L.SetGlobal("kvGet", L.NewFunction(luaKVGet))
	L.SetGlobal("kvDel", L.NewFunction(luaKVDel))

	botFilepath := filepath.Join(viper.GetString("bots_lua_path"), script)
	if err := L.DoFile(botFilepath); err != nil {
		log.Println("Bot script not opening: ", err)
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("main"), // name of Lua function
		NRet:    1,                   // number of returned values
		Protect: true,                // return err or panic
	}, lua.LString(author), lua.LString(message)); err != nil {
		log.Println("Bot script error: ", err)
	}

	if result, ok := L.Get(-1).(lua.LString); ok {
		return string(result), true
	}

	return "", false
}

func execPythonModule(author, message, module string) (string, bool) {
	python.Initialize()
	defer python.Finalize()

	python.PySys_SetPath(viper.GetString("bots_python_path"))

	pyBotModule := python.PyImport_ImportModule(module)
	if pyBotModule == nil {
		log.Println("Error importing python module: ", module)
		return "", false
	}

	pyBotMain := pyBotModule.GetAttrString("main")
	if pyBotMain == nil {
		log.Println("Error importing python function main in module", module)
		return "", false
	}

	pyKeyAuthor := python.PyString_FromString("author")
	pyAuthor := python.PyString_FromString(author)
	pyKeyMessage := python.PyString_FromString("message")
	pyMesssage := python.PyString_FromString(message)

	pyKwargs := python.PyDict_New()
	python.PyDict_SetItem(pyKwargs, pyKeyAuthor, pyAuthor)
	python.PyDict_SetItem(pyKwargs, pyKeyMessage, pyMesssage)

	// The Python function takes no params but when using the C api
	// we're required to send (empty) *args and **kwargs anyways.
	result := pyBotMain.Call(python.PyTuple_New(0), pyKwargs)

	if python.PyString_Check(result) {
		message := python.PyString_AS_STRING(result)
		return message, true
	}

	log.Println("Function main in returned not string in module: ", module)
	return "", false
}

func ExecHandler(author, message, script string) (string, bool) {
	if matched, _ := regexp.Match(`.+\.py`, []byte(script)); matched {
		module := strings.Split(script, ".")[0]
		message, ok := execPythonModule(author, message, module)

		return message, ok
	}

	if matched, _ := regexp.Match(`.+\.lua`, []byte(script)); matched {
		message, ok := execLuaScript(author, message, script)

		return message, ok
	}

	log.Println("Not found script: ", script)
	return "", false
}
