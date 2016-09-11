#!/usr/bin/env bash

set -euo pipefail

dir=$(pwd)
report_dir=$dir/report
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
    go test -tags "'${tags}'" -run - -benchmem -bench 00_C github.com/cockroachdb/cockroach/sql 2>/dev/null >> $report
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
  rm -f $report
  for no in `seq $count`; do
    go test -tags "'${tags}'" -benchmem -bench . ./bench >> $report
  done
}

if [[ $# -ge 1 ]]; then
  count=$1
fi
echo "iteration count ${count}"
mkdir -p $report_dir

# micro benchmarks
branch='develop'
(
  cd $ck
  git checkout $branch
)
sql_00_C $branch 'llrb'

branch='final'
(
  cd $ck
  git checkout $branch
)
micro $branch 'llrb'
micro $branch 'btree'
sql_00_C $branch 'llrb'
sql_00_C $branch 'btree'

branch='no-degree'
(
  cd $ck
  git checkout $branch
)
micro $branch 'btree'
sql_00_C $branch 'btree'
