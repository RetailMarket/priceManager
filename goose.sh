#!/usr/bin/env bash
goose -env=test -pgschema=$1 $2
echo 'goose completed...'
