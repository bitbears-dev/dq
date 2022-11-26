# 📅 dq: jq for date / time ⏲

## Examples

```
$ dq  # => shows current date / time in JSON format
{
  "am": false,
  "day": 23,
  "dayOfYear": 296,
  "daysInMonth": 31,
  "hour": 22,
  "hour12": 10,
  "leapYear": false,
  "microsecond": 479854,
  "millisecond": 479,
  "minute": 59,
  "month": 10,
  "nanosecond": 479854622,
  "rfc3339": "2022-10-23T22:59:42+09:00",
  "second": 42,
  "timezone": {
    "offsetSeconds": 32400,
    "short": "JST"
  },
  "unix": 1666533582,
  "unixMicro": 1666533582479854,
  "unixMicroString": "1666533582479854",
  "unixMilli": 1666533582479,
  "unixMilliString": "1666533582479",
  "unixNano": 1666533582479854622,
  "unixNanoString": "1666533582479854622",
  "unixString": "1666533582",
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
$ dq 'fromunix(1666533582)'  # => shows info for the specified date / time.
                             # Note: single quotation marks `' ... '` around the filter expression
                             # are required to prevent parenthesis from being interpreted by the shell
{
  "am": false,
  "day": 23,
  "dayOfYear": 296,
  "daysInMonth": 31,
  "hour": 22,
  "hour12": 10,
  "leapYear": false,
  "microsecond": 0,
  "millisecond": 0,
  "minute": 59,
  "month": 10,
  "nanosecond": 0,
  "rfc3339": "2022-10-23T22:59:42+09:00",
  "second": 42,
  "timezone": {
    "offsetSeconds": 32400,
    "short": "JST"
  },
  "unix": 1666533582,
  "unixMicro": 1666533582000000,
  "unixMicroString": "1666533582000000",
  "unixMilli": 1666533582000,
  "unixMilliString": "1666533582000",
  "unixNano": 1666533582000000000,
  "unixNanoString": "1666533582000000000",
  "unixString": "1666533582",
  "weekday": {
    "name": "Sunday"
  },
  "year": 2022
}
```

```
$ echo 1666533582 | dq fromunix  # => you can provide the input from another process via the pipe
# => (the result is the same above)
```

```
$ dq 'fromrfc3339("2022-10-23T23:03:01+09:00")'
{
  "am": false,
  "day": 23,
  "dayOfYear": 296,
  "daysInMonth": 31,
  "hour": 23,
  "hour12": 11,
  "leapYear": false,
  "microsecond": 0,
  "millisecond": 0,
  "minute": 3,
  "month": 10,
  "nanosecond": 0,
  "rfc3339": "2022-10-23T23:03:01+09:00",
  "second": 1,
  "timezone": {
    "offsetSeconds": 32400,
    "short": "JST"
  },
  "unix": 1666533781,
  "unixMicro": 1666533781000000,
  "unixMicroString": "1666533781000000",
  "unixMilli": 1666533781000,
  "unixMilliString": "1666533781000",
  "unixNano": 1666533781000000000,
  "unixNanoString": "1666533781000000000",
  "unixString": "1666533781",
  "weekday": {
    "name": "Sunday"
  },
  "year": 2022
}
```

```
$ dq 'fromrfc3339("2022-10-23T23:03:01+09:00") | add_date(0;0;1) | .weekday.name'
# => shows the name of weekday of the next day
"Monday"
```

```
$ dq -r 'fromrfc3339("2022-10-23T23:03:01+09:00") | add_date(0;0;1) | .weekday.name'
# => raw output (double quotation marks will be removed when the result is only a string)
Monday
```

## Install

TBD


## Reference

### Types

<details>
<summary><code>time</code></summary>

  | Field name        | Type       | Description                                                             |
  | -------------     | ---------- | ----------------------------------------------------------------------- |
  | `am`              | bool       | `true`: AM, `false`: PM                                                 |
  | `day`             | integer    | Day of the month                                                        |
  | `dayOfYear`       | integer    | Day of the year                                                         |
  | `daysInMonth`     | integer    | Number of days in the month                                             |
  | `hour`            | integer    | Hour within the day, 24-hour format i.e. in range [0, 23]               |
  | `hour12`          | integer    | Hour within the day, 12-hour format i.e. in range [0, 12]               |
  | `leapYear`        | bool       | Whether the year is a leap year or not                                  |
  | `microsecond`     | integer    | Microsecond offset within the second, in range [0, 999999]              |
  | `millisecond`     | integer    | Millisecond offset within the second, in range [0, 999]                 |
  | `minute`          | integer    | Minute offset within the hour, in range [0, 59]                         |
  | `month`           | integer    | Month of the year                                                       |
  | `nanosecond`      | integer    | Nanosecond offset within the second, in range [0, 999999999]            |
  | `rfc3339`         | string     | RFC 3339 style string represents this time object                       |
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
</details>

<details>
<summary><code>timezone</code></summary>

  | Field name      | Type    | Description                  |
  | --------------- | ------- | ---------------------------- |
  | `offsetSeconds` | integer | Offset in seconds            |
  | `short`         | string  | Abbreviated name of the zone |
</details>

<details>
<summary><code>weekday</code></summary>

  | Field name      | Type    | Description             |
  | --------------- | ------- | ----------------------- |
  | `name`          | string  | English name of the day |
</details>

### Functions

- Format conversion

  - Let `dq` guess format of the input
    <details>
    <summary><code>guess</code> (<code>g</code>) </summary>
    Generate $time$ object from input.

    $in: integer \vert float \vert string \rightarrow t:time$

    - $in$: Unix time or known string format of time.
      - $in$ can be provided from input stream or the first item of the arguments. i.e. both of the following are supported:
        - `echo '1666533582' | dq guess`
        - `dq 'guess(1666533582)'`

    - $t$: $time$ object representing the specified time.
    </details>

  - Unix time
    <details>
    <summary><code>fromunix</code> (<code>from_unix</code>)</summary>

      Generate $time$ object from Unix time.

      $in:integer \vert float \vert string \rightarrow t:time $

      - $in$: Unix time represented in integer, floating point number or string. e.g. `1666533582`, `166533582.694357016` or `"1666533582"`

        - $in$ can be provided from input stream or the first item of the arguments. i.e. both of the following are supported:

          - `echo '1666533582' | dq fromunix`
          - `dq 'fromunix(1666533582)'`

      - $t$: $time$ object representing local time.

    </details>

    <details>
    <summary><code>fromunixmilli</code> (<code>from_unixmilli</code>)</summary>

      Generate $time$ object from Unix time (in milliseconds).

      $in:integer \vert float \vert string \rightarrow t:time $

      - $in$: Unix time in milliseconds represented in integer or string. e.g. `1666533582694`, `166533582.694357016` or `"1666533582694"`

        $in$ can be provided from input stream or the first item of the arguments. i.e. both of the following are supported:

        - `echo '1666533582694' | dq fromunixmilli`
        - `dq 'fromunixmilli(1666533582694)'`

      - $t$: $time$ object representing local time.

    </details>

    <details>
    <summary><code>fromunixmicro</code> (<code>from_unixmicro</code>)</summary>

      Generate $time$ object from Unix time (in microseconds).

      $in:integer \vert float \vert string \rightarrow t:time $

      - $in$: Unix time in microseconds represented in integer or string. e.g. `1666533582694357`, `166533582.694357016` or `"1666533582694357"`

        $in$ can be provided from input stream or the first item of the arguments. i.e. both of the following are supported:

        - `echo '1666533582694357' | dq fromunixmicro`
        - `dq 'fromunixmicro(1666533582694357)'`

      - $t$: $time$ object representing local time.

    </details>


    <details>
    <summary><code>fromunixnano</code> (<code>from_unixnano</code>)</summary>

      Generate $time$ object from Unix time (in nanoseconds).

      $in:integer \vert float \vert string \rightarrow t:time $

      - $in$: Unix time in nanoseconds represented in integer or string. e.g. `1666533582694357016`, `166533582.694357016` or `"1666533582694357016"`

        $in$ can be provided from input stream or the first item of the arguments. i.e. both of the following are supported:

        - `echo '1666533582694357016' | dq fromunixnano`
        - `dq 'fromunixnano(1666533582694357016)'`

      - $t$: $time$ object representing local time.

    </details>

  - ANSI C Style
    <details>
    <summary><code>fromansic</code> (<code>from_ansic</code>)</summary>

      Generate $time$ object from an ANSI C style string.

      $in: string \rightarrow t:time$

      - $in$: ANSI C style string. e.g. "Fri Oct 28 05:59:07 2022"

        $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

        - `echo '"Fri Oct 28 05:59:07 2022"' | dq fromansic`
        - `echo 'Fri Oct 28 05:59:07 2022' | dq -R fromansic`
        - `dq 'fromansic("Fri Oct 28 05:59:07 2022")'`

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

  - Unix `date` command style

    <details>
    <summary><code>fromunixdate</code> (<code>from_unixdate</code>)</summary>

      Generate $time$ object from a Unix date style string. "Unix date style" means the output format of `date` command with `LC_TIME=C`.

      $in: string \rightarrow t:time$

      - $in$: Unix date style string. e.g. "Fri Oct 28 05:59:07 JST 2022"

        $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

        - `echo '"Fri Oct 28 05:59:07 JST 2022"' | dq fromunixdate`
        - `echo 'Fri Oct 28 05:59:07 JST 2022' | dq -R fromunixdate`
        - `dq 'fromunixdate("Fri Oct 28 05:59:07 JST 2022")'`

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

  - Ruby style
    <details>
    <summary><code>fromrubydate</code> (<code>from_rubydate</code>)</summary>

      Generate $time$ object from a Ruby Date style string.

      $in: string \rightarrow t:time$

      - $in$: Ruby Date style string. e.g. "Fri Oct 28 05:59:07 +0900 2022"

        $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

        - `echo '"Fri Oct 28 05:59:07 +0900 2022"' | dq fromrubydate`
        - `echo 'Fri Oct 28 05:59:07 +0900 2022' | dq -R fromrubydate`
        - `dq 'fromrubydate("Fri Oct 28 05:59:07 +0900 2022")'`

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

  - RFC 822
    <details>
    <summary><code>fromrfc822</code> (<code>from_rfc822</code>)</summary>

      Generate $time$ object from a RFC822 string.

      $in: string \rightarrow t:time$

      - $in$: RFC822 string. e.g. "28 Oct 22 05:59:07 JST"

        $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

        - `echo '"28 Oct 22 05:59 JST"' | dq fromrfc822`
        - `echo '28 Oct 22 05:59 JST' | dq -R fromrfc822`
        - `dq 'fromrfc822("28 Oct 22 05:59 JST")'`

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
        #=> "28 Oct 22 05:40 JST"
        ```

    </details>

    <details>
    <summary><code>fromrfc822z</code> (<code>from_rfc822z</code>)</summary>

      Generate $time$ object from a RFC822 with numeric zone string.

      $in: string \rightarrow t:time$

      - $in$: RFC822 with numeric zone string. e.g. "28 Oct 22 05:59:07 +0900"

        $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

        - `echo '"28 Oct 22 05:59 +0900"' | dq fromrfc822z`
        - `echo '28 Oct 22 05:59 +0900' | dq -R fromrfc822z`
        - `dq 'fromrfc822z("28 Oct 22 05:59 +0900")'`

      - $t$: $time$ object representing the specified time.

    </details>

    <details>
    <summary><code>torfc822z</code> (<code>to_rfc822z</code>)</summary>

      Generate RFC822 with numeric zone string represents $time$ object.

      $t: time \rightarrow out: string$

      - $t$: $time$ object

      - $out$: RFC822 with numeric zone string represents the time specified by the $time$ object.

        e.g.)
        ```
        echo '1666903217' | TZ=Asia/Tokyo dq 'fromunix | torfc822z'
        #=> "28 Oct 22 05:40 +0900"
        ```

    </details>

  - RFC 850
    <details>
    <summary><code>fromrfc850</code> (<code>from_rfc850</code>)</summary>

      Generate $time$ object from a RFC850 style string.

      $in: string \rightarrow t:time$

      - $in$: RFC850 style string. e.g. "Friday, 28-Oct-22 05:59:07 JST"

        $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

        - `echo '"Friday, 28-Oct-22 05:59:17 JST"' | dq fromrfc850`
        - `echo 'Friday, 28-Oct-22 05:59:17 JST' | dq -R fromrfc850`
        - `dq 'fromrfc850("Friday, 28-Oct-22 05:59:17 JST")'`

      - $t$: $time$ object representing the specified time.

    </details>

    <details>
    <summary><code>torfc850</code> (<code>to_rfc850</code>)</summary>

      Generate RFC850 style string represents $time$ object.

      $t: time \rightarrow out: string$

      - $t$: $time$ object

      - $out$: RFC850 style string represents the time specified by the $time$ object.

        e.g.)
        ```
        echo '1666903217' | TZ=Asia/Tokyo dq 'fromunix | torfc850'
        #=> "Friday, 28-Oct-22 05:40:17 JST"
        ```

    </details>

  - RFC 1123
    <details>
    <summary><code>fromrfc1123</code> (<code>from_rfc1123</code>)</summary>

      Generate $time$ object from a RFC1123 style string.

      $in: string \rightarrow t:time$

      - $in$: RFC1123 style string. e.g. "Fri, 28 Oct 2022 05:40:17 JST"

        $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

        - `echo '"Fri, 28 Oct 2022 05:40:17 JST"' | dq fromrfc1123`
        - `echo 'Fri, 28 Oct 2022 05:40:17 JST' | dq -R fromrfc1123`
        - `dq 'fromrfc1123("Fri, 28 Oct 2022 05:40:17 JST")'`

      - $t$: $time$ object representing the specified time.

    </details>

    <details>
    <summary><code>torfc1123</code> (<code>to_rfc1123</code>)</summary>

      Generate RFC1123 style string represents $time$ object.

      $t: time \rightarrow out: string$

      - $t$: $time$ object

      - $out$: RFC1123 style string represents the time specified by the $time$ object.

        e.g.)
        ```
        echo '1666903217' | TZ=Asia/Tokyo dq 'fromunix | torfc1123'
        #=> "Fri, 28 Oct 2022 05:40:17 JST"
        ```

    </details>

    <details>
    <summary><code>fromrfc1123z</code> (<code>from_rfc1123z</code>)</summary>

      Generate $time$ object from a RFC1123 with numeric zone string.

      $in: string \rightarrow t:time$

      - $in$: RFC1123z style string. e.g. "Fri, 28 Oct 2022 05:40:17 +0900"

        $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

        - `echo '"Fri, 28 Oct 2022 05:40:17 +0900"' | dq fromrfc1123z`
        - `echo 'Fri, 28 Oct 2022 05:40:17 +0900' | dq -R fromrfc1123z`
        - `dq 'fromrfc1123z("Fri, 28 Oct 2022 05:40:17 +0900")'`

      - $t$: $time$ object representing the specified time.

    </details>

    <details>
    <summary><code>torfc1123z</code> (<code>to_rfc1123z</code>)</summary>

      Generate RFC1123 with numeric zone string represents $time$ object.

      $t: time \rightarrow out: string$

      - $t$: $time$ object

      - $out$: RFC1123 with numeric zone string represents the time specified by the $time$ object.

        e.g.)
        ```
        echo '1666903217' | TZ=Asia/Tokyo dq 'fromunix | torfc1123z'
        #=> "Fri, 28 Oct 2022 05:40:17 +0900"
        ```

    </details>

  - RFC 3339
    <details>
    <summary><code>fromrfc3339</code> (<code>from_rfc3339</code>)</summary>

      Generate $time$ object from a RFC3339 style string.

      $in: string \rightarrow t:time$

      - $in$: RFC3339 style string. e.g. "2022-10-23T23:03:01+09:00"

        $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

        - `echo '"2022-10-23T23:03:01+09:00"' | dq fromrfc3339`
        - `echo '2022-10-23T23:03:01+09:00' | dq -R fromrfc3339`
        - `dq 'fromrfc3339("2022-10-23T23:03:01+09:00")'`

      - $t$: $time$ object representing the specified time.

    </details>

    <details>
    <summary><code>torfc3339</code> (<code>to_rfc3339</code>)</summary>

      Generate RFC3339 style string represents $time$ object.

      $t: time \rightarrow out: string$

      - $t$: $time$ object

      - $out$: RFC3339 style string represents the time specified by the $time$ object.

        e.g.)
        ```
        echo '1666533781' | TZ=Asia/Tokyo dq 'fromunix | torfc3339'
        #=> "2022-10-23T23:03:01+09:00"
        ```

    </details>

    <details>
    <summary><code>fromrfc3339nano</code> (<code>from_rfc3339nano</code>)</summary>

      Generate $time$ object from a RFC3339 with nanoseconds string.

      $in: string \rightarrow t:time$

      - $in$: RFC3339 with nanoseconds string. e.g. "2022-10-23T23:03:01.123456789+09:00"

        $in$ can be provided from input stream or the first item of the arguments. i.e. all of the following are supported:

        - `echo '"2022-10-23T23:03:01.123456789+09:00"' | dq fromrfc3339nano`
        - `echo '2022-10-23T23:03:01.123456789+09:00' | dq -R fromrfc3339nano`
        - `dq 'fromrfc3339nano("2022-10-23T23:03:01.123456789+09:00")'`

      - $t$: $time$ object representing the specified time.

    </details>

    <details>
    <summary><code>torfc3339nano</code> (<code>to_rfc3339nano</code>)</summary>

      Generate RFC3339 with nanoseconds string represents $time$ object.

      $t: time \rightarrow out: string$

      - $t$: $time$ object

      - $out$: RFC3339 with nanoseconds string represents the time specified by the $time$ object.

      e.g.)
      ```
      echo '1666533781123456789' | TZ=Asia/Tokyo dq 'fromunixnano | torfc3339nano'
      #=> "2022-10-23T23:03:01.123456789+09:00"
      ```

    </details>

  - `strptime` / `strftime` / `mktime`
    <details>
    <summary><code>strptime</code></summary>

    __This function is derived from__ `jq`

    Interface to the C-library function `strptime()`.

    $in: string, fmt: string \rightarrow out: array of integers$

    - $in$: input string
      - $in$ must be provided via the input stream.
    - $fmt$: format string (refer to the host operating system's documentation for details)
      - $fmt$ must be specified as an argument.
    - $out$: array of integers which can be used an input for `mktime` function

    e.g.)
    ```
    $ echo '"2022-11-03T12:58:47Z"' | dq -c 'strptime("%Y-%m-%dT%H:%M:%SZ")'
    # => [2022,10,3,12,58,47,4,306]
    $ echo '"2022-11-03T12:58:47Z"' | dq 'strptime("%Y-%m-%dT%H:%M:%SZ") | mktime'
    # => 1667480327
    ```

    </details>

    <details>
    <summary><code>strftime</code></summary>

    __This function is derived from__ `jq`

    Interface to the C-library function `strftime()`.

    $in: integer, fmt: string \rightarrow out: string$

    - $in$: unix time
      - $in$ must be provided via the input stream.
    - $fmt$: format string (refer to the host operating system's documentation for details)
      - $fmt$ must be specified as an argument.
    - $out$: string representation of the specified unix time in the specified format

    e.g.)
    ```
    $ dq -r '.unix | strftime("%Y/%m/%d %H:%M:%S")'
    # => 2022/11/21 13:25:06
    ```
    </details>

    <details>
    <summary><code>mktime</code></summary>

    __This function is derived from__ `jq`

    Interface to the C-library function `mktime()`.

    $in: array of integers \rightarrow out: unix time$

    - $in$: array of integers. output of `strptime`.
      - $in$ must be provided via the input stream.
    - $out$: unix time

    e.g.)
    ```
    $ echo '"2022-11-03T12:58:47Z"' | dq 'strptime("%Y-%m-%dT%H:%M:%SZ") | mktime'
    1667480327
    ```
    </details>

  - From Year / Month / Day
    <details>
    <summary><code>fromymd</code> (<code>from_ymd</code>) </summary>

    Generate $time$ object by specifying year, month and day.

    $y, m, d: integer \rightarrow t: time$

    - $y$: year
    - $m$: month [1-12]
    - $d$: day [1-31]
    - $t$: $time$ object representing the specified date

    e.g.)
    ```
    $ dq 'fromymd(2022;11;26) | .rfc3339'
    "2022-11-26T00:00:00+09:00"
    ```
    </details>

  - From Year / Month / Day / Hour / Minute / Second
    <details>
    <summary><code>fromymdhms</code> (<code>from_ymdhms</code>) </summary>

    Generate $time$ object by specifying year, month, day, hour, minute and second.

    $y, mon, d, h, min, s: integer \rightarrow t: time$

    - $y$: year
    - $mon$: month [1-12]
    - $d$: day [1-31]
    - $h$: hour [0-23]
    - $min$: minute [0-59]
    - $sec$: second [0-59]
    - $t$: $time$ object representing the specified date / time

    e.g.)
    ```
    $ dq 'fromymdhms(2022;11;26;9;23;18) | .rfc3339'
    "2022-11-26T09:23:18+09:00"
    ```
    </details>

  - "Kitchen"
    <details>
    <summary><code>fromkitchen</code> (<code>from_kitchen</code>)</summary>

    Generate $time$ object from "kitchen clock" style string.

    $k: string \rightarrow t: time$

    - $k$: kitchen style string e.g.) `"3:04PM"`
    - $t$: $time$ object representing the specified time.

    e.g.)
    ```
    $ dq 'fromkitchen("1:56PM") | .hour'
    13
    ```
    </details>

    <details>
    <summary><code>tokitchen</code> (<code>to_kitchen</code>)</summary>

    Generate "kitchen clock" style string representing the specified $time$ object

    $t: time \rightarrow k: string$

    - $t$: $time$ object
    - $k$: kitchen style string e.g.) `"3:04PM"`

    e.g.)
    ```
    $ dq tokitchen
    "2:03PM"
    ```
    </details>

  - "Stamp"
    <details>
    <summary><code>fromstamp</code> (<code>from_stamp</code>) </summary>

    Generate $time$ object from a "handy timestamp" style string.

    $s: string \rightarrow t: time$

    - $s$: handy timestamp style string. e.g.) `"Nov 26 15:40:53"`
    - $t$: $time$ object

    e.g.)
    ```
    $ dq 'fromstamp("Nov 26 15:40:53") | .daysInMonth'
    30
    ```
    </details>

    <details>
    <summary><code>tostamp</code> (<code>to_stamp</code>) </summary>

    Generate a "handy timestamp" style string representing the specified $time$ object.

    $t: time \rightarrow s: string$

    - $t$: $time$ object
    - $s$: handy timestamp style string. e.g.) `"Nov 26 15:40:53"`

    e.g.)
    ```
    $ dq tostamp
    "Nov 26 15:40:53"
    ```

    </details>

- Calculation

  <details>
  <summary><code>add_date</code></summary>

  $t: time, y: integer, m: integer, d: integer \rightarrow out: time$

  - $t$: $time$ object
    - $t$ must be specified via the input stream
  - $y$: number of years to add. can be zero or negative value.
  - $m$: number of months to add. can be zero or negative value.
  - $d$: number of days to add. can be zero or negative value.
  - $out$ $time$ object

  e.g.)
  ```
  $ dq -r 'fromrfc3339("2022-10-23T23:03:01+09:00") | add_date(0; 0; 1) | .rfc3339'
  # => 2022-10-24T23:03:01+09:00
  $ dq -r 'fromrfc3339("2022-10-23T23:03:01+09:00") | add_date(0; 1; 0) | .rfc3339'
  # => 2022-11-23T23:03:01+09:00
  $ dq -r 'fromrfc3339("2022-10-23T23:03:01+09:00") | add_date(1; 0; 0) | .rfc3339'
  # => 2023-10-23T23:03:01+09:00
  ```
  </details>
