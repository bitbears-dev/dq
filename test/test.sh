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

dq_supports_fromrfc3339_filter() {
  progress "dq supports fromrfc3339() filter"
  result="$( $bin 'fromrfc3339("2022-10-23T23:03:01+09:00")' )"
  assert_json "$result"
  assert_json_has_field "$result" "unix"
  assert_json_field_has_value "$result" "unix" "1666533781"
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

dq_without_arguments
dq_with_a_simple_filter
dq_supports_fromunix_filter
dq_fromunix_filter_supports_stdin
dq_supports_fromrfc3339_filter
dq_supports_add_day_filter
dq_supports_raw_output

test_result=0
