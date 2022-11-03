#!/usr/bin/env bash
set -Eeuo pipefail

d="$( cd "$( dirname "$0" )"; cd ..; pwd -P )"

bin=${DQ_BIN:-"$d/cmd/dq/dq"}

RESET=$'\e[0m'
BOLD=$'\e[1m'
GREEN=$'\e[0;32m'
RED=$'\e[0;31m'

test_result=1
trap 'tear_down' 0
tear_down() {
  : "Report result" && {
    if [ "$test_result" -eq 0 ]; then
      echo
      echo -e "${GREEN}${BOLD}TEST OK${RESET}"
      echo
    else
      echo
      echo -e "${RED}${BOLD}TEST FAILED${RESET}"
      echo
    fi
    exit $test_result
  }
}

progress() {
  local message=$1
  echo -n "${BOLD}${message}${RESET} " 1>&2
}

print_fail() {
    echo "${RED}FAIL${RESET}" 1>&2
}

print_ok() {
    echo "${GREEN}OK${RESET}" 1>&2
}

fail_with_message() {
    local message=$1

    print_fail
    echo -e "$message" 1>&2
    exit 1
}

compact() {
  echo "$1" | jq -c
}

assert_eq() {
  if [ "$1" != "$2" ]; then
    fail_with_message "expected: $2\nactual:   $1"
  fi
}

assert_json() {
  local x=$1
  if ! echo "$x" | jq >/dev/null 2>&1; then
    fail_with_message "$( compact "$x" ) is not json"
  fi
}

assert_json_has_field() {
  local x=$1
  local field=$2
  local val
  val="$( echo "$x" | jq ".$field" )"
  if [ "$val" == "null" ]; then
    fail_with_message "$( compact "$x" ) does not have field '$field'"
  fi
}

assert_json_does_not_have_field() {
  local x=$1
  local field=$2
  local val
  val="$( echo "$x" | jq ".$field" )"
  if [ "$val" != "null" ]; then
    fail_with_message "$( compact "$x" ) has field '$field'"
  fi
}

assert_json_field_has_value() {
  local x=$1
  local field=$2
  local expected=$3
  local actual
  actual="$( echo "$x" | jq ".$field" )"
  if [ "$actual" != "$expected" ]; then
    fail_with_message "expected: $expected\nactual:   $actual"
  fi
}

assert_match() {
  local x=$1
  local regexp=$2
  [[ "$x" =~ $regexp ]]
}

dq_without_arguments() {
  progress "dq without arguments (expects json representing current date / time to be printed)"
  result="$( $bin )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "year" "$( date +%Y )"
  assert_json_does_not_have_field "$result" "_source"  # should not expose internal field '_source'
  print_ok
}

dq_with_a_simple_filter() {
  progress "dq with a simple filter"
  result="$( $bin .year )"
  assert_eq "$result" "$( date +%Y )"
  print_ok
}

dq_supports_fromunix_filter() {
  progress "dq supports fromunix() filter"
  result="$( $bin 'fromunix(1666533582)' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666533582"
  print_ok
}

dq_fromunix_filter_supports_stdin() {
  progress "dq fromunix() filter supports stdin"
  result="$( echo '1666533582' | $bin fromunix )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666533582"
  print_ok
}

dq_supports_fromunixmilli_filter() {
  progress "dq supports fromunixmilli() filter"
  result="$( $bin 'fromunixmilli(1666533582694)' )"
  assert_json "$result"
  assert_json_has_field "$result" "unixMilli"
  assert_json_field_has_value "$result" "unixMilli" "1666533582694"
  print_ok
}

dq_fromunixmilli_filter_supports_stdin() {
  progress "dq fromunixmilli() filter supports stdin"
  result="$( echo '1666533582694' | $bin fromunixmilli )"
  assert_json "$result"
  assert_json_has_field "$result" "unixMilli"
  assert_json_field_has_value "$result" "unixMilli" "1666533582694"
  print_ok
}

dq_supports_fromunixmicro_filter() {
  progress "dq supports fromunixmicro() filter"
  result="$( $bin 'fromunixmicro(1666533582694357)' )"
  assert_json "$result"
  assert_json_has_field "$result" "unixMicro"
  assert_json_field_has_value "$result" "unixMicro" "1666533582694357"
  print_ok
}

dq_fromunixmicro_filter_supports_stdin() {
  progress "dq fromunixmicro() filter supports stdin"
  result="$( echo '1666533582694357' | $bin fromunixmicro )"
  assert_json "$result"
  assert_json_has_field "$result" "unixMicro"
  assert_json_field_has_value "$result" "unixMicro" "1666533582694357"
  print_ok
}

dq_supports_fromunixnano_filter() {
  progress "dq supports fromunixnano() filter"
  result="$( $bin 'fromunixnano(1666533582694357016)' )"
  assert_json "$result"
  assert_json_has_field "$result" "unixNanoString"
  assert_json_field_has_value "$result" "unixNanoString" '"1666533582694357016"'
  print_ok
}

dq_fromunixnano_filter_supports_stdin() {
  progress "dq fromunixnano() filter supports stdin"
  result="$( echo '1666533582694357016' | $bin fromunixnano )"
  assert_json "$result"
  assert_json_has_field "$result" "unixNanoString"
  assert_json_field_has_value "$result" "unixNanoString" '"1666533582694357016"'
  print_ok
}

dq_supports_fromansic_filter() {
  progress "dq supports fromansic() filter"
  result="$( $bin 'fromansic("Fri Oct 28 05:59:07 2022")' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666936747"
  print_ok
}

dq_fromansic_filter_supports_stdin() {
  progress "dq fromansic() filter supports stdin"
  result="$( echo 'Fri Oct 28 05:59:07 2022' | $bin -R fromansic )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666936747"
  print_ok
}

dq_supports_toansic_filter() {
  progress "dq supports toansic() filter"
  result="$( echo '1666936747' | $bin 'fromunix | utc | toansic' )"
  assert_eq "$result" '"Fri Oct 28 05:59:07 2022"'
  print_ok
}

dq_supports_fromunixdate_filter() {
  progress "dq supports fromunixdate() filter"
  result="$( $bin 'fromunixdate("Fri Oct 28 05:40:17 JST 2022")' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903217"
  print_ok
}

dq_fromunixdate_filter_supports_stdin() {
  progress "dq fromunixdate() filter supports stdin"
  result="$( echo 'Fri Oct 28 05:40:17 JST 2022' | $bin -R fromunixdate )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903217"
  print_ok
}

dq_supports_tounixdate_filter() {
  progress "dq supports tounixdate() filter"
  result="$( echo '1666903217' | $bin 'fromunix | tounixdate ' )"
  assert_eq "$result" '"Fri Oct 28 05:40:17 JST 2022"'
  print_ok
}

dq_supports_fromrubydate_filter() {
  progress "dq supports fromrubydate() filter"
  result="$( $bin 'fromrubydate("Fri Oct 28 05:40:17 +0900 2022")' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903217"
  print_ok
}

dq_fromrubydate_filter_supports_stdin() {
  progress "dq fromrubydate() filter supports stdin"
  result="$( echo 'Fri Oct 28 05:40:17 +0900 2022' | $bin -R fromrubydate )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903217"
  print_ok
}

dq_supports_torubydate_filter() {
  progress "dq supports torubydate() filter"
  result="$( echo '1666903217' | $bin 'fromunix | torubydate ' )"
  assert_eq "$result" '"Fri Oct 28 05:40:17 +0900 2022"'
  print_ok
}

dq_supports_fromrfc822_filter() {
  progress "dq supports fromrfc822() filter"
  result="$( $bin 'fromrfc822("28 Oct 22 05:40 JST")' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903200"
  print_ok
}

dq_fromrfc822_filter_supports_stdin() {
  progress "dq fromrfc822() filter supports stdin"
  result="$( echo '28 Oct 22 05:40 JST' | $bin -R fromrfc822 )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903200"
  print_ok
}

dq_supports_torfc822_filter() {
  progress "dq supports torfc822() filter"
  result="$( echo '1666903200' | $bin 'fromunix | torfc822' )"
  assert_eq "$result" '"28 Oct 22 05:40 JST"'
  print_ok
}

dq_supports_fromrfc822z_filter() {
  progress "dq supports fromrfc822z() filter"
  result="$( $bin 'fromrfc822z("28 Oct 22 05:40 +0900")' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903200"
  print_ok
}

dq_fromrfc822z_filter_supports_stdin() {
  progress "dq fromrfc822z() filter supports stdin"
  result="$( echo '28 Oct 22 05:40 +0900' | $bin -R fromrfc822z )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903200"
  print_ok
}

dq_supports_torfc822z_filter() {
  progress "dq supports torfc822z() filter"
  result="$( echo '1666903200' | $bin 'fromunix | torfc822z' )"
  assert_eq "$result" '"28 Oct 22 05:40 +0900"'
  print_ok
}

dq_supports_fromrfc850_filter() {
  progress "dq supports fromrfc850() filter"
  result="$( $bin 'fromrfc850("Friday, 28-Oct-22 05:40:17 JST")' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903217"
  print_ok
}

dq_fromrfc850_filter_supports_stdin() {
  progress "dq fromrfc850() filter supports stdin"
  result="$( echo 'Friday, 28-Oct-22 05:40:17 JST' | $bin -R fromrfc850 )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903217"
  print_ok
}

dq_supports_torfc850_filter() {
  progress "dq supports torfc850() filter"
  result="$( echo '1666903217' | $bin 'fromunix | torfc850' )"
  assert_eq "$result" '"Friday, 28-Oct-22 05:40:17 JST"'
  print_ok
}

dq_supports_fromrfc1123_filter() {
  progress "dq supports fromrfc1123() filter"
  result="$( $bin 'fromrfc1123("Fri, 28 Oct 2022 05:40:17 JST")' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903217"
  print_ok
}

dq_fromrfc1123_filter_supports_stdin() {
  progress "dq fromrfc1123() filter supports stdin"
  result="$( echo 'Fri, 28 Oct 2022 05:40:17 JST' | $bin -R fromrfc1123 )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903217"
  print_ok
}

dq_supports_torfc1123_filter() {
  progress "dq supports torfc1123() filter"
  result="$( echo '1666903217' | $bin 'fromunix | torfc1123' )"
  assert_eq "$result" '"Fri, 28 Oct 2022 05:40:17 JST"'
  print_ok
}

dq_supports_fromrfc1123z_filter() {
  progress "dq supports fromrfc1123z() filter"
  result="$( $bin 'fromrfc1123z("Fri, 28 Oct 2022 05:40:17 +0900")' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903217"
  print_ok
}

dq_fromrfc1123z_filter_supports_stdin() {
  progress "dq fromrfc1123z() filter supports stdin"
  result="$( echo 'Fri, 28 Oct 2022 05:40:17 +0900' | $bin -R fromrfc1123z )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666903217"
  print_ok
}

dq_supports_torfc1123z_filter() {
  progress "dq supports torfc1123z() filter"
  result="$( echo '1666903217' | $bin 'fromunix | torfc1123z' )"
  assert_eq "$result" '"Fri, 28 Oct 2022 05:40:17 +0900"'
  print_ok
}

dq_supports_fromrfc3339_filter() {
  progress "dq supports fromrfc3339() filter"
  result="$( $bin 'fromrfc3339("2022-10-23T23:03:01+09:00")' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666533781"
  print_ok
}

dq_fromrfc3339_filter_supports_stdin() {
  progress "dq fromrfc3339() filter supports stdin"
  result="$( echo '2022-10-23T23:03:01+09:00' | $bin -R fromrfc3339 )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666533781"
  print_ok
}

dq_supports_torfc3339_filter() {
  progress "dq supports torfc3339() filter"
  result="$( echo '1666533781' | $bin 'fromunix | torfc3339' )"
  assert_eq "$result" '"2022-10-23T23:03:01+09:00"'
  print_ok
}

dq_supports_fromrfc3339nano_filter() {
  progress "dq supports fromrfc3339nano() filter"
  result="$( $bin 'fromrfc3339nano("2022-10-23T23:03:01.123456789+09:00")' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666533781"
  assert_json_has_field "$result" "unixNano"
  assert_json_field_has_value "$result" "unixNano" "1666533781123456800"  # somehow it's rounded
  print_ok
}

dq_fromrfc3339nano_filter_supports_stdin() {
  progress "dq fromrfc3339nano() filter supports stdin"
  result="$( echo '2022-10-23T23:03:01.123456788+09:00' | $bin -R fromrfc3339nano )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666533781"
  assert_json_has_field "$result" "unixNano"
  assert_json_field_has_value "$result" "unixNano" "1666533781123456800"  # somehow it's rounded
  print_ok
}

dq_supports_torfc3339nano_filter() {
  progress "dq supports torfc3339nano() filter"
  result="$( echo '1666533781123456789' | $bin 'fromunixnano | torfc3339nano' )"
  assert_eq "$result" '"2022-10-23T23:03:01.123456789+09:00"'
  print_ok
}

dq_fromunix_can_parse_floating_point_unix_time() {
  progress "dq fromunix can parse floating point unix time"
  result="$( $bin 'now | fromunix' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  print_ok
}

dq_fromunixmilli_can_parse_floating_point_unix_time() {
  progress "dq fromunixmilli can parse floating point unix time"
  result="$( $bin 'now | fromunixmilli' )"
  assert_json "$result"
  assert_json_has_field "$result" "unixMilli"
  print_ok
}

dq_fromunixmicro_can_parse_floating_point_unix_time() {
  progress "dq fromunixmicro can parse floating point unix time"
  result="$( $bin 'now | fromunixmicro' )"
  assert_json "$result"
  assert_json_has_field "$result" "unixMicro"
  print_ok
}

dq_fromunixnano_can_parse_floating_point_unix_time() {
  progress "dq fromunixnano can parse floating point unix time"
  result="$( $bin 'now | fromunixnano' )"
  assert_json "$result"
  assert_json_has_field "$result" "unixNano"
  print_ok
}

dq_can_use_strptime_output() {
  progress "dq can use strptime output"
  result="$( echo '"2022-11-03T12:58:47Z"' | $bin 'strptime("%Y-%m-%dT%H:%M:%SZ") | mktime | fromunix | utc' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "year" "2022"
  assert_json_field_has_value "$result" "month" "11"
  assert_json_field_has_value "$result" "day" "3"
  assert_json_field_has_value "$result" "hour" "12"
  assert_json_field_has_value "$result" "minute" "58"
  assert_json_field_has_value "$result" "second" "47"
  print_ok
}

dq_can_use_strftime() {
  progress "dq can use strftime"
  result="$( $bin -r '.unix | strftime("%Y-%m-%dT%H:%M:%SZ")' )"
  assert_match "$result" "^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$"
  print_ok
}


dq_supports_add_day_filter() {
  progress "dq supports add_day() filter"
  result="$( $bin 'fromrfc3339("2022-10-23T23:03:01+09:00") | add_day(1) | .weekday.name' )"
  assert_eq "$result" '"Monday"'
}

dq_supports_raw_output() {
  progress "dq supports raw output"
  result="$( $bin -r 'fromrfc3339("2022-10-23T23:03:01+09:00") | add_day(1) | .weekday.name' )"
  assert_eq "$result" 'Monday'
}

# basics
dq_without_arguments
dq_with_a_simple_filter

# fromunix()
dq_supports_fromunix_filter
dq_fromunix_filter_supports_stdin

# fromunixmilli()
dq_supports_fromunixmilli_filter
dq_fromunixmilli_filter_supports_stdin

# fromunixmicro()
dq_supports_fromunixmicro_filter
dq_fromunixmicro_filter_supports_stdin

# fromunixnano()
dq_supports_fromunixnano_filter
dq_fromunixnano_filter_supports_stdin

# fromansic()
dq_supports_fromansic_filter
dq_fromansic_filter_supports_stdin

# toansic()
dq_supports_toansic_filter

# fromunixdate()
dq_supports_fromunixdate_filter
dq_fromunixdate_filter_supports_stdin

# tounixdate()
dq_supports_tounixdate_filter

# fromrubydate()
dq_supports_fromrubydate_filter
dq_fromrubydate_filter_supports_stdin

# torubydate()
dq_supports_torubydate_filter

# fromrfc822()
dq_supports_fromrfc822_filter
dq_fromrfc822_filter_supports_stdin

# torfc822()
dq_supports_torfc822_filter

# fromrfc822z()
dq_supports_fromrfc822z_filter
dq_fromrfc822z_filter_supports_stdin

# torfc822z()
dq_supports_torfc822z_filter

# fromrfc850()
dq_supports_fromrfc850_filter
dq_fromrfc850_filter_supports_stdin

# torfc850()
dq_supports_torfc850_filter

# fromrfc1123()
dq_supports_fromrfc1123_filter
dq_fromrfc1123_filter_supports_stdin

# torfc1123()
dq_supports_torfc1123_filter

# fromrfc1123z()
dq_supports_fromrfc1123z_filter
dq_fromrfc1123z_filter_supports_stdin

# torfc1123z()
dq_supports_torfc1123z_filter

# fromrfc3339()
dq_supports_fromrfc3339_filter
dq_fromrfc3339_filter_supports_stdin

# torfc3339()
dq_supports_torfc3339_filter

# fromrfc3339nano()
dq_supports_fromrfc3339nano_filter
dq_fromrfc3339nano_filter_supports_stdin

# torfc3339nano()
dq_supports_torfc3339nano_filter

# interop with jq's date/time functions
dq_fromunix_can_parse_floating_point_unix_time
dq_fromunixmilli_can_parse_floating_point_unix_time
dq_fromunixmicro_can_parse_floating_point_unix_time
dq_fromunixnano_can_parse_floating_point_unix_time
dq_can_use_strptime_output
dq_can_use_strftime

# add_day()
dq_supports_add_day_filter
dq_supports_raw_output

test_result=0
