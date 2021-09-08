---
title: "Specification"
date: 2019-05-28T14:28:03-07:00
draft: false
weight: 1
---

![](/images/sqlcommenter_logo.png)

- [Introduction](#introduction)
- [Format](#format)
- [Comment escaping](#comment-escaping)
- [Meta characters](#meta-characters)
    - [Algorithm](#meta-characters-algorithm)
- [Key-Value serialization](#key-value-serialization)
    - [Key serialization](#key-serialization)
        - [Algorithm](#key-serialization-algorithm)
    - [Value serialization](#key-serialization)
        - [Algorithm](#value-serialization-algorithm)
- [Sorting](#sorting)
    - [Algorithm](#sorting-algorithm)
    - [Exhibit](#sorting-exhibit)
- [Concatenation](#contentation)
    - [Separator](#separator)
    - [Algorithm](#concatenation-algorithm)
    - [Exhibit](#concatenation-exhibit)
- [Affix comment](#affix-comment)
    - [Algorithm](#affix-comment-algorithm)
    - [Exhibit](#affix-comment-exhibit)
- [SQL commenter](#sql-commenter)
    - [Exhibit](#sql-commenter-exhibit)
- [Parsing](#parsing)
    - [Algorithm](#parsing-algorithm)
    - [Exhibit](#parsing-exhibit)
- [References](#references)

### Introduction

This section defines the [SQL commenter algorithm](/) which augments a SQL statement with a comment containing serialized
key value pairs that are retrieved from the various ORMs and frameworks in your programming language and environment of choice.

A preview of the result can be seen as per [exhibit](#sql-commenter-exhibit)
```python
SELECT * FROM FOO /*action='%2Fparam*d',controller='index,'framework='spring',
traceparent='00-5bd66ef5095369c7b0d1f8f4bd33716a-c532cb4098ac3dd2-01',
tracestate='congo%3Dt61rcWkgMzE%2Crojo%3D00f067aa0ba902b7'*/
```

Read along to see how you can conform to the specification and produce similar output.

### Format

The final comment SHOULD be affixed to the final SQL statement in the format

```shell
<SQL STATEMENT> /*<ATTRIBUTE_KEY_VALUE_PAIRS>*/
```

### Comment escaping
[Comments within SQL comments](https://docs.oracle.com/cd/B12037_01/server.101/b10759/sql_elements006.htm) are of the format

* Following ` -- ` e.g. `SELECT * from FOO -- this is the comment`
* Contained within `/*` and `*/` e.g. `SELECT * from FOO /* this is the comment */`

If a comment already exists within a SQL statement, we MUST NOT mutate that statement.

### Separator
Each key value pair MUST be separated by a comma "," so for example, given

Values: `[<FIELD_1>, <FIELD_2>, <FIELD_3>, ...]`

Expected concatenation result: `<FIELD_1>,<FIELD_2>,<FIELD_3>,<FIELD_N...>`

### Meta characters

Meta characters such as `'` should be escaped with a slash `\`. That creates the following algorithm:

#### <a name="meta-characters-algorithm"></a> Algorithm

```python
algorithm(value):
    escaped := value.escape_with_slash_any_of(')

    return escaped
```

### Key serialization

1. [URL encode](#url-encode) the key e.g. given `route parameter`, that'll become `route%20parameter`

Which produces the following algorithm:

#### <a name="key-serialization-algorithm"></a> Algorithm

```python
key_serialization(key):
    encoded := url_encode(key)
    meta_escaped := escape_meta_characters(encoded)

    return meta_escaped
```

### Value serialization

1. [URL encode](#url-encode) the value e.g. given `/param first`, that SHOULD become `%2Fparam%20first`
2. Escape meta-characters within the raw value; a single quote `'` becomes `\'`
3. SQL escape the value by placing it within two single quotes e.g.

    `DROP` should become `'DROP'`

    `FOO 'BAR` should become `'FOO%20\'BAR'`

And when generalized into an algorithm:

#### <a name="value-serialization-algorithm"></a> Algorithm
```python
value_serialization(value):
    encoded := url_encode(value)
    meta_escaped := escape_meta_characters(encoded)
    final := sql_escape_with_single_quotes(meta_escaped)

    return final
```

and running the algorithm on the following table will produce

value|url_encode(value)|sql_escape_with_single_quotes
---|---|---
`DROP TABLE FOO`|`DROP%20TABLE%20FOO`|`'DROP%20TABLE%20FOO'`
`/param first`|`%2Fparam%20first`|`'%2Fparam%20first'`
`1234`|`1234`|`'1234'`


### Key Value format
Given a key value pair (key, value):

1. Run the [Key serialization algorithm](#key-serialization-algorithm) on `key`
2. Run the [Value serialization algorithm](#value-serialization-algorithm) on `value`
3. Using an equals sign `=`, concatenate the result from 1. and 2. to give

    `<SERIALIZED_KEY>=<SERIALIZED_VALUE>`
gotten from:

    `serialize_key(key)=serialize_value(value)`

Thus given for example the following key value pairs

key value pair|serialized_key|serialized_value|Final
---|---|---|---
`route=/polls 1000`|`route`|`'%2Fpolls%201000'`|`route='%2Fpolls%201000'`
`name='DROP TABLE FOO'`|`route`|`'%2Fpolls%201000'`|`route='%2Fpolls%201000'`
`name''="DROP TABLE USERS'"`|name=\'\'|DROP%20TABLE%20USERS\'|name=\'\'='DROP%20TABLE%20USERS\''

### Sorting

With a list of serialized `key=value` pairs, sort them by lexicographic order.

#### <a name="sorting-algorithm"></a> Algorithm

```python
sort(key_value_pairs):
    sorted = lexicographically_sort(key_value_pairs)

    return sorted
```

#### <a name="sorting-exhibit"></a> Exhibit

Thus

```python
    sort([
        traceparent='00-5bd66ef5095369c7b0d1f8f4bd33716a-c532cb4098ac3dd2-01',
        tracestate='congo%3Dt61rcWkgMzE%2Crojo%3D00f067aa0ba902b7',
        route='%2Fparam*d',
        controller='index',
    ])
```

produces

```python
 [
        controller='index',
        route='%2Fparam*d',
        traceparent='00-5bd66ef5095369c7b0d1f8f4bd33716a-c532cb4098ac3dd2-01',
        tracestate='congo%3Dt61rcWkgMzE%2Crojo%3D00f067aa0ba902b7',
 ]
```

### Concatenation

After all the keys and values have been serialized and sorted, they MUST be joined by a comma `,`.

If no values are present, `concatenate` MUST return the empty value `''`

#### <a name="concatenation-algorithm"></a> Algorithm
```python
concatenate(key_value_pairs):

    if len(key_value_pairs) == 0:
        return ''

    return ','.join(key_value_pairs)
```

#### <a name="concatenation-exhibit"></a> Exhibit
Therefore

```python
    concatenate([
        controller='index',
        route='%2Fparam*d',
        traceparent='00-5bd66ef5095369c7b0d1f8f4bd33716a-c532cb4098ac3dd2-01',
        tracestate='congo%3Dt61rcWkgMzE%2Crojo%3D00f067aa0ba902b7',
    ])
```

produces

```python
controller='index',route='%2Fparam*d',traceparent='00-5bd66ef5095369c7b0d1f8f4bd33716a-c532cb4098ac3dd2-01',tracestate='congo%3Dt61rcWkgMzE%2Crojo%3D00f067aa0ba902b7'
```

### Affix comment

After [serialization](#serialization), [sorting](#sorting), [concatenation](#concatenation), the final form MUST be placed between `/*` and `*/`

#### <a name="affix-comment"></a> Algorithm
```python
affix_comment(sql, concatenated):
    if is_empty(concatenated):
        return sql // Do NOT modify the SQL if concatenated is blank.

    affixed := sql + '/*' + concatenated + '*/'

    return affixed
```

#### <a name="affix-comment-exhibit"></a> Exhibit
for example given

```python
affix_comment('SELECT * from FOO', '')
```

produces
```python
SELECT * from FOO
```

```python
affix_comment('SELECT * from FOO', "route='%2Fparam*d'")
```

produces
```python
SELECT * from FOO /*route='%2Fparam*d'*/
```

### SQL commenter

Wrapping all the steps together, we thus have the following algorithm

```python
sql_commenter(sql, attributes):
    if contains_sql_comment(sql):
        return sql # DO NOT mutate a statement with an already present comment.

    serialized_key_value_pairs := []

    for each attribute in attributes:
        serialized := serialize_key_value_pair(attribute)
        if serialized:
            serialized_key_value_pairs.append(serialized)

    sorted := sort(serialized_key_value_pairs)
    concatenated := concatenate(sorted)
    final := affix_comment(sql, concatenated)

    return final
```

#### <a name="sql-commenter-exhibit"></a> Exhibit

Running [sql_commenter](#sql-commenter) on an ORM integration that extracts the respective attributes:

```python
sql_commenter('SELECT * FROM FOO', [
        tracestate='congo%3Dt61rcWkgMzE%2Crojo%3D00f067aa0ba902b7',
        traceparent='00-5bd66ef5095369c7b0d1f8f4bd33716a-c532cb4098ac3dd2-01',
        framework='spring',
        action='%2Fparam*d',
        controller='index',
])
```

finally produces

```python
SELECT * FROM FOO /*action='%2Fparam*d',controller='index,'framework='spring',
traceparent='00-5bd66ef5095369c7b0d1f8f4bd33716a-c532cb4098ac3dd2-01',
tracestate='congo%3Dt61rcWkgMzE%2Crojo%3D00f067aa0ba902b7'*/
```

### Parsing

Parsing is the step to reverse sql-commenter and extract the key value attributes.

It'll follow the following steps:

1. Find the last comment so search for and strip out `/*` and `*/`
2. Split the comment by comma `,`
3. Split each `key='value'` pair so extract `key` and `'value'`
    - 3.1. For `key`, `unescape_meta_characters` then `url_decode`
    - 3.2. For `value`, sql_unescape/trim the `'` at the beginning and end of `'value'` -> `value`
        - 3.2.1. Unescape the meta characters in `value`
        - 3.2.2. URL Decode the value

#### <a name="parsing-algorithm"></a> Algorithm
```go
parse(sql_with_comment):
    if !contains_sql_comment(sql_with_comment):
        return sql_with_comment, null

    // Since we now have a SQL comment, let's extract the serialized attributes.
    sql_stmt, serialized_attrs := extract_sql_commenter(sql_with_comment)

    if is_empty(serialized_attrs):
        return sql_stmt, null

    attrs := {}
    kv_splits := split_by_comma(serialized_attrs)
    for kv in kv_splits:
        e_key, e_value := split_by_equals(kv)
        key := decode_key(e_key)
        value := decode_value(e_value)

        attrs[key] = value

    // Some attributes such as traceparent, tracestate, sampled
    // might need need some grouping and reconstruction.
    final := deconstruct_and_group_attributes(attrs)

    return sql_stmt, final
```

#### <a name="parsing-exhibit"></a> Exhibit

Given the value from [SQLCommenter exhibit](#sql-commenter-exhibit)

```python
SELECT * FROM FOO /*action='%2Fparam*d',controller='index,'framework='spring',
traceparent='00-5bd66ef5095369c7b0d1f8f4bd33716a-c532cb4098ac3dd2-01',
tracestate='congo%3Dt61rcWkgMzE%2Crojo%3D00f067aa0ba902b7'*/
```

Running `parse` on the value

```go
sql, attributes = parse(`SELECT * FROM FOO
/*action='%2Fparam*d',controller='index,'framework='spring',
traceparent='00-5bd66ef5095369c7b0d1f8f4bd33716a-c532cb4098ac3dd2-01',
tracestate='congo%3Dt61rcWkgMzE%2Crojo%3D00f067aa0ba902b7'*/`)
```

produces

```python
sql: SELECT * FROM FOO
attributes: {
    controller: 'index',
    framework: 'spring',
    action: '/param*d',
    trace: {
        sampled: true,
        span_id: 'c532cb4098ac3dd2',
        trace_id: '5bd66ef5095369c7b0d1f8f4bd33716a',
        trace_state: [{'congo': 't61rcWkgMzE'}, {'rojo': '00f067aa0ba902b7'}],
    },
}
```


### References

Resource|URL
---|---
URL Encoding|https://en.wikipedia.org/wiki/Percent-encoding
Comments within SQL comments|https://docs.oracle.com/cd/B12037_01/server.101/b10759/sql_elements006.htm
