#!/usr/bin/env bash
goose -env=test -pgschema=price $1
echo 'goose completed...'
