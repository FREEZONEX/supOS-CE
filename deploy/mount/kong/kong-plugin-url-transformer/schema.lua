local typedefs = require "kong.db.schema.typedefs"

return {
    name = "supos-url-transformer",
    fields = {
        { config = {
            type = "record",
            fields = {
                { home_url = { type = "string" , default = "/home"} },
            },
        },
        },
    },
}