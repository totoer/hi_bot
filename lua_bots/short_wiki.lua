--! !HiBot (факт!)

function main(author, message)
    result = httpGet("https://ru.wikipedia.org/w/api.php?format=json&action=query&generator=random&prop=description")
    if result.status == 200 then
        body = jsonDecode(result.body)
        for key, value in pairs(body.query.pages) do
            if value.description ~= nil then
                return value.title .. " - " .. value.description
            else
                return main(author, message)
            end
        end
    end
end