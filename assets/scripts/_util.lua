---highlight some value
---@param val any
function highlight(val)
    return text_underline(text_bold("[" .. tostring(val) .. "]"))
end
