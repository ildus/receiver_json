import json
import sys
from utils import sqlgen, consts as c

RULES = {
    "tables": {
        "public.table1": {
            "first_name": {"type": c.FIELD_USING_PK, "prefix": "n"},
            "passport_seria": {"type": c.FIELD_USING_PK, "lpad": (4, '0')},
            "full_name": {"type": c.FIELD_USING_PK, "prefix": "n", "postfix": "d"},
        },
        "public.table2": {
            "__where": "idx != 5",
            "number": {"type": c.FIELD_USING_PK, "lpad": (6, ' '), "pk": "idx"},
            "flat": {"type": c.FIELD_REWRITE, "value": None},
        },
        "public.tmp1": c.TABLE_SKIP,
        "public.tmp2": c.TABLE_TRUNCATE,
    }
}

if '--json' in sys.argv:
    print json.dumps(RULES)
elif '--sql' in sys.argv:
    sqlgen.generate_sql(RULES['tables'])
