#!/usr/bin/env bash

dir=$HOME/report/00/final

mkdir -p $dir

git checkout 'develop'

develop=$dir/develop
rm -f $develop
for no in {1..10}; do
  go test -run - -benchmem -bench 00_C ./sql 2>/dev/null >> $develop
done

git checkout 'final'
llrb=$dir/llrb
rm -f $llrb
for no in {1..10}; do
  go test -run - -benchmem -bench 00_C ./sql 2>/dev/null >> $llrb
done

btree=$dir/btree
rm -f $btree
for no in {1..10}; do
  go test -tags 'btree' -run - -benchmem -bench 00_C ./sql 2>/dev/null >> $btree
done
