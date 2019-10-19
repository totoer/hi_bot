--! !HiBot (test .+)

function main(author, message)
    if string.find(message, "test jsonDecode") ~= nil then
        print("Test jsonDecode")
        result = jsonDecode('{"test": "Hello World"}')
        return result.test
    end

    if string.find(message, "test jsonEncode") ~= nil then
        print("Test jsonEncode")
        result = jsonEncode({test = "Hello World", key = {foo = 1}})
        return result
    end

    if string.find(message, "test kvSet") ~= nil then
        print("Test kvSet")
        kvSet({store = "Test", key = "Test", value = "Hello World"})
        return "Ok!"
    end

    if string.find(message, "test kvGet") ~= nil then
        print("Test kvGet")
        value = kvGet({store = "Test", key = "Test"})
        return value
    end

    if string.find(message, "test kvDel") ~= nil then
        print("Test kvDel")
        kvDel({store = "Test", key = "Test"})
        return "Ok!"
    end
end
