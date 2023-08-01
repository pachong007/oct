#!/bin/bash
go build
mv comics tx
/home/comics/tx > tx.log 2>&1 &