# coding: utf-8

import consts as c


def drop_table(table):
    return 'DROP TABLE %s IF EXISTS CASCADE;' % table


def truncate_table(table):
    return 'TRUNCATE TABLE %s;' % table


def update_str(field_name, options):
    modify_type = options['type']
    new_val = None
    if modify_type == c.FIELD_USING_PK:
        prefix = options.get('prefix')
        postfix = options.get('postfix')
        lpad = options.get('lpad')
        pk_field = options.get('pk') or 'id'

        new_val = ''
        if prefix:
            new_val = "'%s' || " % prefix

        cast_str = 'cast(%s as text)' % pk_field
        if lpad is not None:
            count, ch = lpad
            cast_str = "lpad(%s, %s, '%s')" % (cast_str, count, ch)

        new_val += cast_str
        if postfix:
            new_val += " || '%s'" % postfix

    elif modify_type == c.FIELD_REWRITE:
        new_val = options.get("value")
        if isinstance(new_val, basestring):
            new_val = "'%s'" % new_val
        elif new_val is None:
            new_val = 'null'

    if new_val is None:
        raise Exception("Cannot construct UPDATE str")

    return "\n  %s=%s" % (field_name, new_val)


def modify_fields(table, fields):
    base = "UPDATE %s \nSET %s"
    if '__where' in fields:
        base += '\nWHERE ' + fields.pop('__where')

    equals = (update_str(field, options) for field, options in fields.iteritems())
    return base % (table, ', '.join(equals)) + ';'


def generate_sql(tables):
    for table, data in tables.iteritems():
        # print "--- %s" % table
        if isinstance(data, int):
            if data == c.TABLE_SKIP:
                print drop_table(table)
            elif data == c.TABLE_TRUNCATE:
                print truncate_table(table)
        elif isinstance(data, dict):
            print modify_fields(table, data)
        print ""
