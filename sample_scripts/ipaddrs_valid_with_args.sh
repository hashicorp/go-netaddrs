#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0



same_line() {
cat << EndOfMessage
"172.25.41.79"  "172.25.16.77" "172.25.42.80"
EndOfMessage
}

multi_line_ipv6() {
cat << EndOfMessage
"2001:0db8:85a3:0000:0000:8a2e:0370:7334"
    "fe80::1ff:fe23:4567:890a%3"
"2001:db8::"
EndOfMessage
}

if [[ "$1" == "same-line" ]]; then
    same_line
fi

if [[ "$1" == "multi-line" ]]; then
    multi_line_ipv6
fi