#!/bin/bash

generate_signature() {
  local signingKey="$1"
  local method="$2"
  local url="$3"
  local timestamp="$4"

  local dataToSign="$method|$url|$timestamp"

  local signature=$(echo -n "$dataToSign" | openssl dgst -sha256 -hmac "$signingKey" | sed 's/^.* //')

  echo "$signature"
}

if [ "$#" -ne 2 ]; then
  echo "usage: $0 <method> <url>"
  exit 1
fi

secretKey="$TRINQUET_REQUEST_SIGNING_KEY"
method="$1"
url="$2"
timestamp=$(($(date +%s) * 1000))

signature=$(generate_signature "$secretKey" "$method" "$url" "$timestamp")

echo "signature: $signature timestamp: $timestamp"
