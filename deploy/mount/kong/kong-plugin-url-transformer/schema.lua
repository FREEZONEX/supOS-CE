local typedefs = require "kong.db.schema.typedefs"

return {
    name = "url-redirect",
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