package executor

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/spf13/viper"

	"github.com/boltdb/bolt"
	lua "github.com/yuin/gopher-lua"
)

func mapToLTable(m map[string]interface{}) *lua.LTable {
	// Main table pointer
	resultTable := &lua.LTable{}

	// Loop map
	for key, element := range m {

		switch element.(type) {
		case float64:
			resultTable.RawSetString(key, lua.LNumber(element.(float64)))
		case int64:
			resultTable.RawSetString(key, lua.LNumber(element.(int64)))
		case string:
			resultTable.RawSetString(key, lua.LString(element.(string)))
		case bool:
			resultTable.RawSetString(key, lua.LBool(element.(bool)))
		case []byte:
			resultTable.RawSetString(key, lua.LString(string(element.([]byte))))
		case map[string]interface{}:

			// Get table from map
			tble := mapToLTable(element.(map[string]interface{}))

			resultTable.RawSetString(key, tble)

		case time.Time:
			resultTable.RawSetString(key, lua.LNumber(element.(time.Time).Unix()))

		case []map[string]interface{}:

			// Create slice table
			sliceTable := &lua.LTable{}

			// Loop element
			for _, s := range element.([]map[string]interface{}) {

				// Get table from map
				tble := mapToLTable(s)

				sliceTable.Append(tble)
			}

			// Set slice table
			resultTable.RawSetString(key, sliceTable)

		case []interface{}:

			// Create slice table
			sliceTable := &lua.LTable{}

			// Loop interface slice
			for _, s := range element.([]interface{}) {

				// Switch interface type
				switch s.(type) {
				case map[string]interface{}:

					// Convert map to table
					t := mapToLTable(s.(map[string]interface{}))

					// Append result
					sliceTable.Append(t)

				case float64:

					// Append result as number
					sliceTable.Append(lua.LNumber(s.(float64)))

				case string:

					// Append result as string
					sliceTable.Append(lua.LString(s.(string)))

				case bool:

					// Append result as bool
					sliceTable.Append(lua.LBool(s.(bool)))
				}
			}

			// Append to main table
			resultTable.RawSetString(key, sliceTable)
		}
	}

	return resultTable
}

func luaJSONDecode(L *lua.LState) int {
	content := L.ToString(1)
	var rawResult map[string]interface{}
	if err := json.Unmarshal([]byte(content), &rawResult); err != nil {
		log.Println("luaJSONDecode err: ", err)
		return 0
	}

	result := mapToLTable(rawResult)

	L.Push(result)

	return 1
}

func luaJSONEncode(L *lua.LState) int {
	lTable := L.ToTable(1)
	result := make(map[string]interface{})
	stack := make([]map[string]interface{}, 10)
	stack = append(stack, result)

	var prepare func(lua.LValue, lua.LValue)

	prepare = func(key lua.LValue, value lua.LValue) {
		sKey := string(key.(lua.LString))

		switch {
		case value.Type() == lua.LTNil:
			stack[len(stack)-1][sKey] = nil
		case value.Type() == lua.LTBool:
			stack[len(stack)-1][sKey] = lua.LVAsBool(value)
		case value.Type() == lua.LTNumber:
			stack[len(stack)-1][sKey] = float64(value.(lua.LNumber))
		case value.Type() == lua.LTString:
			stack[len(stack)-1][sKey] = string(value.(lua.LString))
		case value.Type() == lua.LTTable:
			stack[len(stack)-1][sKey] = make(map[string]interface{})
			stack = append(stack, stack[len(stack)-1][sKey].(map[string]interface{}))
			value.(*lua.LTable).ForEach(prepare)
			stack = stack[:len(stack)-1]
		}
	}

	lTable.ForEach(prepare)

	if r, err := json.Marshal(result); err == nil {
		L.Push(lua.LString(r))
	}

	return 1
}

func luaHTTPGet(L *lua.LState) int {
	url := L.ToString(1)
	response, err := http.Get(url)
	result := L.NewTable()

	if err != nil {
		result.RawSetString("status", lua.LNumber(-1))
		result.RawSetString("body", lua.LString(""))
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		body := buf.String()

		result.RawSetString("status", lua.LNumber(response.StatusCode))
		result.RawSetString("body", lua.LString(body))
	}

	L.Push(result)

	return 1
}

func luaKVSet(L *lua.LState) int {
	dbFilepath := filepath.Join(viper.GetString("bots_path"), viper.GetString("db_name"))
	db, err := bolt.Open(dbFilepath, 0600, nil)
	if err != nil {
		log.Println("luaKVSet err: ", err)
		return 0
	}
	defer db.Close()

	record := L.ToTable(1)

	rStore := record.RawGetString("store")
	store := []byte(rStore.(lua.LString))

	rKey := record.RawGetString("key")
	key := []byte(rKey.(lua.LString))

	rValue := record.RawGetString("value")
	value := []byte(rValue.(lua.LString))

	db.Update(func(tx *bolt.Tx) error {
		var err error

		b, err := tx.CreateBucketIfNotExists(store)
		if err != nil {
			return err
		}
		err = b.Put(key, value)
		return err
	})

	return 1
}

func luaKVGet(L *lua.LState) int {
	dbFilepath := filepath.Join(viper.GetString("bots_path"), viper.GetString("db_name"))
	db, err := bolt.Open(dbFilepath, 0600, nil)
	if err != nil {
		log.Println("luaKVGet err: ", err)
		return 0
	}
	defer db.Close()

	record := L.ToTable(1)

	rStore := record.RawGetString("store")
	store := []byte(rStore.(lua.LString))

	rKey := record.RawGetString("key")
	key := []byte(rKey.(lua.LString))

	var result string

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(store)
		v := b.Get(key)
		result = string(v)
		return nil
	})

	L.Push(lua.LString(result))

	return 1
}

func luaKVDel(L *lua.LState) int {
	dbFilepath := filepath.Join(viper.GetString("bots_path"), viper.GetString("db_name"))
	db, err := bolt.Open(dbFilepath, 0600, nil)
	if err != nil {
		log.Println("luaKVDel err: ", err)
		return 0
	}
	defer db.Close()

	record := L.ToTable(1)

	rStore := record.RawGetString("store")
	store := []byte(rStore.(lua.LString))

	rKey := record.RawGetString("key")
	key := []byte(rKey.(lua.LString))

	db.Update(func(tx *bolt.Tx) error {
		var err error

		b, err := tx.CreateBucketIfNotExists(store)
		if err != nil {
			return err
		}
		err = b.Delete(key)
		return err
	})

	return 1
}
