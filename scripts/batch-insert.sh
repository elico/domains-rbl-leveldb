#!/usr/bin/env bash

curl -i -X POST -H "Content-Type: multipart/form-data" -F "listfile=@BL/porn/domains" "http://localhost:8080/batch/insert/"
