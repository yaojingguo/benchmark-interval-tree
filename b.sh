#!/usr/bin/env bash

set -euo pipefail

saved_dir=$(pwd)
report_dir=$saved_dir/report
ck=$GOPATH/src/github.com/cockroachdb/cockroach
count=10

function sql_00_C() {
  local branch=$1
  local impl=$2
  case $impl in 
    llrb)
      local tags=''
      ;;
    btree)
      local tags='btree'
      ;;
  esac
  local report=$report_dir/${branch}_${impl}_00_C
  rm -f $report
  for no in `seq $count`; do
    go test -tags "'${tags}'" -run - -benchmem -bench 00_C ./sql 2>/dev/null >> $report
  done
}

function micro() {
  local branch=$1
  local impl=$2
  case $impl in 
    llrb)
      local tags=''
      ;;
    btree)
      local tags='btree'
      ;;
  esac
  local report=$report_dir/${branch}_${impl}_micro
  rm -f $report_dir
  for no in `seq $count`; do
    go test -tags "'${tags}'" -benchmem -bench . ./bench >> $report
  done
}

if [[ $# -ge 1 ]]; then
  count=$1
fi
echo "iteration count ${count}"
mkdir -p $report_dir


branch='develop'
micro $branch 'llrb'
branch='final'
micro $branch 'llrb'
micro $branch 'btree'
branch='remove-degree'
micro $branch 'btree'

cd $ck

branch='develop'
git checkout $branch
sql_00_C $branch 'llrb'

branch='final'
git checkout $branch
sql_00_C $branch 'llrb'
sql_00_C $branch 'btree'

branch='remove-degree'
git checkout $branch
sql_00_C $branch 'btree'
cd $saved_dir
