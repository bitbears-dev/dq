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

### Types

- $time$

  | Field name        | Type       | Description                                                             |
  | -------------     | ---------- | ----------------------------------------------------------------------- |
  | `am`              | bool       | `true`: AM, `false`: PM                                                 |
  | `day`             | integer    | Day of the month                                                        |
  | `dayOfYear`       | integer    | Day of the year                                                         |
  | `hour`            | integer    | Hour within the day, 24-hour format i.e. in range [0, 23]               |
  | `hour12`          | integer    | Hour within the day, 12-hour format i.e. in range [0, 12]               |
  | `microsecond`     | integer    | Microsecond offset within the second, in range [0, 999999]              |
  | `millisecond`     | integer    | Millisecond offset within the second, in range [0, 999]                 |
  | `minute`          | integer    | Minute offset within the hour, in range [0, 59]                         |
  | `month`           | integer    | Month of the year                                                       |
  | `nanosecond`      | integer    | Nanosecond offset within the second, in range [0, 999999999]            |
  | `second`          | integer    | Second offset within the minute, in range [0, 59]                       |
  | `timezone`        | $timezone$ | Timezone object to describe the timezone                                |
  | `unix`            | integer    | Unix time, the number of seconds elapsed since January 1, 1970 UTC      |
  | `unixMicro`       | integer    | Unix time, the number of microseconds elapsed since January 1, 1970 UTC |
  | `unixMicroString` | string     | String representation of `unixMicro`                                    |
  | `unixMilli`       | integer    | Unix time, the number of milliseconds elapsed since January 1, 1970 UTC |
  | `unixMilliString` | string     | String representation of `unixMilli`                                    |
  | `unixNano`        | integer    | Unix time, the number of nanoseconds elapsed since January 1, 1970 UTC <br/> Note: Some JSON implementations cannot handle integer values of `unixNano` correctly because the values exceed the number of significant digits. Use `unixNanoString` instead to avoid the issue. |
  | `unixNanoString`  | string     | String representation of `unixNano`                                     |
  | `unixString`      | string     | String representation of `unix`                                         |
  | `weekday`         | $weekday$  | Weekday object to describe the day of the week                          |
  | `year`            | integer    | Year                                                                    |

- $timezone$

  | Field name      | Type    | Description                  |
  | --------------- | ------- | ---------------------------- |
  | `offsetSeconds` | integer | Offset in seconds            |
  | `short`         | string  | Abbreviated name of the zone |

- $weekday$

  | Field name      | Type    | Description             |
  | --------------- | ------- | ----------------------- |
  | `name`          | string  | English name of the day |


### Functions

- Format conversion
  <details>
  <summary><code>fromunix</code> (<code>from_unix</code>)</summary>

    Generate $time$ object from Unix time.

    $in:integer \vert float \vert string \rightarrow t:time $

    - $in$: Unix time represented in integer, floating point number or string. e.g. `1666533582`, `166533582.694357016` or `"1666533582"`

      $in$ can be provided from input stream or the first item of the arguments. i.e. both of the following are supported:

      - `echo '1666533582' | dq fromunix`
      - `dq fromunix(1666533582)`

    - $t$: $time$ object representing local time.

  </details>

  <details>
  <summary><code>fromunixmilli</code> (<code>from_unixmilli</code>)</summary>

    Generate $time$ object from Unix time (in milliseconds).

    $in:integer \vert float \vert string \rightarrow t:time $

    - $in$: Unix time in milliseconds represented in integer or string. e.g. `1666533582694`, `166533582.694357016` or `"1666533582694"`

      $in$ can be provided from input stream or the first item of the arguments. i.e. both of the following are supported:

      - `echo '1666533582694' | dq fromunixmilli`
      - `dq fromunixmilli(1666533582694)`

    - $t$: $time$ object representing local time.

  </details>

  <details>
  <summary><code>fromunixmicro</code> (<code>from_unixmicro</code>)</summary>

    Generate $time$ object from Unix time (in microseconds).

    $in:integer \vert float \vert string \rightarrow t:time $

    - $in$: Unix time in microseconds represented in integer or string. e.g. `1666533582694357`, `166533582.694357016` or `"1666533582694357"`

      $in$ can be provided from input stream or the first item of the arguments. i.e. both of the following are supported:

      - `echo '1666533582694357' | dq fromunixmicro`
      - `dq fromunixmicro(1666533582694357)`

    - $t$: $time$ object representing local time.

  </details>


  <details>
  <summary><code>fromunixnano</code> (<code>from_unixnano</code>)</summary>

    Generate $time$ object from Unix time (in nanoseconds).

    $in:integer \vert float \vert string \rightarrow t:time $

    - $in$: Unix time in nanoseconds represented in integer or string. e.g. `1666533582694357016`, `166533582.694357016` or `"1666533582694357016"`

      $in$ can be provided from input stream or the first item of the arguments. i.e. both of the following are supported:

      - `echo '1666533582694357016' | dq fromunixnano`
      - `dq fromunixnano(1666533582694357016)`

    - $t$: $time$ object representing local time.

  </details>

  <details>
  <summary><code>fromansic</code> (<code>from_ansic</code>)</summary>

    Generate $time$ object from an ANSI C style string.

    $in: string \rightarrow t:time$

    - $in$: ANSI C style string. e.g. "Fri Oct 28 05:59:07 2022"

      $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

      - `echo '"Fri Oct 28 05:59:07 2022"' | dq fromansic`
      - `echo 'Fri Oct 28 05:59:07 2022' | dq -R fromansic`
      - `dq fromansic("Fri Oct 28 05:59:07 2022")`

    - $t$: $time$ object representing universal time.

      Note: `fromansic` always parses the input string as UTC. Timezones are not supported.

  </details>

  <details>
  <summary><code>toansic</code> (<code>to_ansic</code>)</summary>

    Generate ANSI C style string represents $time$ object.

    $t: time \rightarrow out: string$

    - $t$: $time$ object

    - $out$: ANSI C style string represents the time specified by the $time$ object.

      e.g.)
      ```
      echo '1666936747' | dq 'fromunix | utc | toansic'
      #=> "Fri Oct 28 05:59:07 2022"
      ```

  </details>

  <details>
  <summary><code>fromunixdate</code> (<code>from_unixdate</code>)</summary>

    Generate $time$ object from a Unix date style string. "Unix date style" means the output format of `date` command with `LC_TIME=C`.

    $in: string \rightarrow t:time$

    - $in$: Unix date style string. e.g. "Fri Oct 28 05:59:07 JST 2022"

      $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

      - `echo '"Fri Oct 28 05:59:07 JST 2022"' | dq fromunixdate`
      - `echo 'Fri Oct 28 05:59:07 JST 2022' | dq -R fromunixdate`
      - `dq fromunixdate("Fri Oct 28 05:59:07 JST 2022")`

    - $t$: $time$ object representing the specified time.

  </details>

  <details>
  <summary><code>tounixdate</code> (<code>to_unixdate</code>)</summary>

    Generate Unix date style string represents $time$ object. "Unix date style" means the output format of `date` command with `LC_TIME=C`

    $t: time \rightarrow out: string$

    - $t$: $time$ object

    - $out$: Unix date style string represents the time specified by the $time$ object.

      e.g.)
      ```
      echo '1666936747' | TZ=Asia/Tokyo dq 'fromunix | tounixdate'
      #=> "Fri Oct 28 14:59:07 JST 2022"
      ```

  </details>

  <details>
  <summary><code>fromrubydate</code> (<code>from_rubydate</code>)</summary>

    Generate $time$ object from a Ruby Date style string.

    $in: string \rightarrow t:time$

    - $in$: Ruby Date style string. e.g. "Fri Oct 28 05:59:07 +0900 2022"

      $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

      - `echo '"Fri Oct 28 05:59:07 +0900 2022"' | dq fromrubydate`
      - `echo 'Fri Oct 28 05:59:07 +0900 2022' | dq -R fromrubydate`
      - `dq fromrubydate("Fri Oct 28 05:59:07 +0900 2022")`

    - $t$: $time$ object representing the specified time.

  </details>

  <details>
  <summary><code>torubydate</code> (<code>to_rubydate</code>)</summary>

    Generate Ruby Date style string represents $time$ object.

    $t: time \rightarrow out: string$

    - $t$: $time$ object

    - $out$: Ruby Date style string represents the time specified by the $time$ object.

      e.g.)
      ```
      echo '1666903217' | TZ=Asia/Tokyo dq 'fromunix | torubydate'
      #=> "Fri Oct 28 05:40:17 +0900 2022"
      ```

  </details>

  <details>
  <summary><code>fromrfc822</code> (<code>from_rfc822</code>)</summary>

    Generate $time$ object from a RFC822 string.

    $in: string \rightarrow t:time$

    - $in$: RFC822 string. e.g. "28 Oct 22 05:59:07 JST"

      $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

      - `echo '"28 Oct 22 05:59:07 JST"' | dq fromrfc822`
      - `echo '28 Oct 22 05:59:07 JST' | dq -R fromrfc822`
      - `dq fromrfc822("28 Oct 22 05:59:07 JST")`

    - $t$: $time$ object representing the specified time.

  </details>

  <details>
  <summary><code>torfc822</code> (<code>to_rfc822</code>)</summary>

    Generate RFC822 style string represents $time$ object.

    $t: time \rightarrow out: string$

    - $t$: $time$ object

    - $out$: RFC822 string represents the time specified by the $time$ object.

      e.g.)
      ```
      echo '1666903217' | TZ=Asia/Tokyo dq 'fromunix | torfc822'
      #=> "28 Oct 22 05:40:17 JST"
      ```

  </details>
