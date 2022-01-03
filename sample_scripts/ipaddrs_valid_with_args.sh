#!/bin/bash


same_line() {
cat << EndOfMessage
"172.25.41.79"  "172.25.16.77" "172.25.42.80"
EndOfMessage
}

multi_line() {
cat << EndOfMessage
"172.25.41.79"
    "172.25.16.77"
"172.25.42.80"
EndOfMessage
}

if [[ $1 -eq "same-line" ]]; then
    same_line
fi

if [[ $1 -eq "multi-line" ]]; then
    multi_line
fi