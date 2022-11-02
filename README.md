# 📅 dq: jq for date / time ⏲

## Examples

```
$ dq  # => shows current date / time in JSON format
{
  "am": false,
  "day": 23,
  "dayOfYear": 296,
  "hour": 22,
  "hour12": 10,
  "microsecond": 479854,
  "millisecond": 479,
  "minute": 59,
  "month": 10,
  "nanosecond": 479854622,
  "second": 42,
  "timezone": {
    "offsetSeconds": 32400,
    "short": "JST"
  },
  "unix": 1666533582,
  "unixString": "1666533582",
  "unixMicro": 1666533582479854,
  "unixMicroString": "1666533582479854",
  "unixMilli": 1666533582479,
  "unixMilliString": "1666533582479",
  "unixNano": 1666533582479854622,
  "unixNanoString": "1666533582479854622",
  "weekday": {
    "name": "Sunday"
  },
  "year": 2022
}
```

```
$ dq .year  # => shows current year
2022
```

```
$ dq 'fromunix(1666533582)'  # => shows info for the specified date / time. Note: single quotation marks `' ... '` around the filter expression are required to prevent parenthesis from being interpreted by the shell
{
  "am": false,
  "day": 23,
  "dayOfYear": 296,
  "hour": 22,
  "hour12": 10,
  "microsecond": 0,
  "millisecond": 0,
  "minute": 59,
  "month": 10,
  "nanosecond": 0,
  "second": 42,
  "timezone": {
    "offsetSeconds": 32400,
    "short": "JST"
  },
  "unix": 1666533582,
  "unixString": "1666533582",
  "unixMicro": 1666533582000000,
  "unixMicroString": "1666533582000000",
  "unixMilli": 1666533582000,
  "unixMilliString": "1666533582000",
  "unixNano": 1666533582000000000,
  "unixNanoString": "1666533582000000000",
  "weekday": {
    "name": "Sunday"
  },
  "year": 2022
}
```

```
$ echo 1666533582 | dq fromunix  # => you can provide the input from another process via the pipe
(same above)
```

```
$ dq 'fromrfc3339("2022-10-23T23:03:01+09:00")'
{
  "am": false,
  "day": 23,
  "dayOfYear": 296,
  "hour": 23,
  "hour12": 11,
  "microsecond": 0,
  "millisecond": 0,
  "minute": 3,
  "month": 10,
  "nanosecond": 0,
  "second": 1,
  "timezone": {
    "offsetSeconds": 32400,
    "short": "JST"
  },
  "unix": 1666533781,
  "unixString": "1666533781",
  "unixMicro": 1666533781000000,
  "unixMicroString": "1666533781000000",
  "unixMilli": 1666533781000,
  "unixMilliString": "1666533781000",
  "unixNano": 1666533781000000000,
  "unixNanoString": "1666533781000000000",
  "weekday": {
    "name": "Sunday"
  },
  "year": 2022
}
```

```
$ dq 'fromrfc3339("2022-10-23T23:03:01+09:00") | add_day(1) | .weekday.name'
# => shows the name of weekday of the next day
"Monday"
```

```
$ ./dq -r 'fromrfc3339("2022-10-23T23:03:01+09:00") | add_day(1) | .weekday.name'
# => raw output (double quotation marks will be removed when the result is only a string)
Monday
```

## Install

TBD


## Reference

- `fromunix` (`from_unix`)

    $in:integer \vert string \rightarrow t:time $

    in: Unix time represented in integer or string. e.g. `1666533582` or `"1666533582"`

    $in$ can be provided from input stream or the first item of the arguments. i.e. both of the following are supported:

    - `echo '1666533582' | dq fromunix`
    - `dq fromunix(1666533582)`


- `fromunixmilli` (`from_unixmilli`)

    $in:integer \vert string \rightarrow t:time $

    in: Unix time in milliseconds represented in integer or string. e.g. `1666533582694` or `"1666533582694"`

    $in$ can be provided from input stream or the first item of the arguments. i.e. both of the following are supported:

    - `echo '1666533582694' | dq fromunixmilli`
    - `dq fromunixmilli(1666533582694)`


- `fromunixmicro` (`from_unixmicro`)

    $in:integer \vert string \rightarrow t:time $

    in: Unix time in microseconds represented in integer or string. e.g. `1666533582694357` or `"1666533582694357"`

    $in$ can be provided from input stream or the first item of the arguments. i.e. both of the following are supported:

    - `echo '1666533582694357' | dq fromunixmicro`
    - `dq fromunixmicro(1666533582694357)`


- `fromunixnano` (`from_unixnano`)

    $in:integer \vert string \rightarrow t:time $

    in: Unix time in nanoseconds represented in integer or string. e.g. `1666533582694357016` or `"1666533582694357016"`

    $in$ can be provided from input stream or the first item of the arguments. i.e. both of the following are supported:

    - `echo '1666533582694357016' | dq fromunixnano`
    - `dq fromunixnano(1666533582694357016)`
