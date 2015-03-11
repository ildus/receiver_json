// field constants
var FIELD = {}
FIELD.PK_WITH_PREFIX = 0

//table constants
var TABLE = {};
TABLE.SKIP = 0
TABLE.TRUNCATE = 1

var RULES = {
    "constants": {
        "table.skip": TABLE.SKIP,
        "field.pk_with_prefix": FIELD.PK_WITH_PREFIX,
    },
    "tables": {
        "public.main_contragent": {
            "first_name": {"type": FIELD.PK_WITH_PREFIX, "prefix": "i"},
            "passport_seria": {"type": FIELD.PK_WITH_PREFIX, "zerofill": 4},
            "passport_number": {"type": FIELD.PK_WITH_PREFIX, "zerofill": 6},
        },
        "public.misc_phone_borrowers": TABLE.SKIP,
    }
}

console.log(JSON.stringify(RULES, null, 4));