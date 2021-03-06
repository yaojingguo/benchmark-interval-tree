#!/usr/bin/env bash

set -euo pipefail

function checkout_issue_6465() {
  (
    cd $fork
    git checkout 'issue-6465'
  )
}

function checkout_issue_6465_opt() {
  (
    cd $fork
    git checkout 'issue-6465-opt'
  )
}

function checkout_rightmosts_guided_bst() {
  (
    cd $fork
    git checkout rightmosts-guided-bst
  )
}

function approach() {
  checkout_rightmosts_guided_bst
  bst_btree_report=$report_dir/bst_btree
  rm -fr $bst_btree_report
  for no in {1..5}; do
    go test -bench . "$btree_based" >> $bst_btree_report
  done

  checkout_issue_6465_opt
  scan_btree_report=$report_dir/scan_btree
  for no in {1..5}; do
    go test -bench . "$btree_based" >> $scan_btree_report
  done
}

function degree() {
  checkout_issue_6465_opt

  local degrees=(2 4 8 16 32 64 128 256)
  rm -f $report_dir/degree_*
  local reports=""
  for d in ${degrees[@]}; do
    echo "degree: $d"
    local report=$report_dir/degree_${d}_btree
    for no in {1..5}; do
      go test -bench "$basic" "$btree_based" -degree $d >> $report
    done
    reports="$reports $report"
  done
  benchstat "$report"
}

function opt() {
  checkout_issue_6465_opt

  btree_wo_opt_report="$report_dir/btree_wo_opt"
  rm -f $btree_wo_opt_report
  checkout_issue_6465 
  for no in {1..5}; do
    go test -bench . "$btree_based" >> $btree_wo_opt_report
  done
  btree_w_opt_report="$report_dir/btree_w_opt"
  rm -f $btree_w_opt_report
  checkout_issue_6465_opt
  for no in {1..5}; do
    go test -bench . "$btree_based" >> $btree_w_opt_report
  done
  benchstat $btree_wo_opt_report $btree_w_opt_report
}

function vs_llrb() {
  checkout_issue_6465_opt

  no_random_btree_report=$report_dir/no_random_btree
  rm -fr $no_random_btree_report
  for no in {1..10}; do
    # go test -benchmem -bench "$no_random" "$btree_based" >> $no_random_btree_report
    go test -benchmem -bench . "$btree_based" >> $no_random_btree_report
  done

  no_random_llrb_report=$report_dir/no_random_llrb
  rm -fr $no_random_llrb_report
  for no in {1..10}; do
    # go test -benchmem -bench "$no_random" "$llrb_based"  >> $no_random_llrb_report
    go test -benchmem -bench . "$llrb_based"  >> $no_random_llrb_report
  done
  benchstat "$no_random_llrb_report" "$no_random_btree_report"
}

function random() {
  checkout_issue_6465_opt

  local lens=(16 32 64 128 256 512 1024)
  for len in ${lens[@]}; do
    local random_btree_report="$report_dir/random_btree_$len"
    rm -fr $random_btree_report
    for no in {1..10}; do
      go test -bench "$random" "$btree_based" -length $len >> $random_btree_report
    done

    local random_llrb_report="$report_dir/random_llrb_$len"
    rm -fr $random_llrb_report
    for no in {1..10}; do
      go test -bench "$random" "$llrb_based"  -length $len >> $random_llrb_report
    done

    echo "benchstat for slice with a random length between 1 and $len"
    local random_stat="$report_dir/random_stat_$len"
    benchstat $random_llrb_report $random_btree_report > $random_stat
  done
}

function llrb_btree_no_free_list() {
  (
    cd $ck
    git checkout 'issue-6465-develop'
  )
  new_llrb_report=$report_dir/new_llrb
  rm -f $new_llrb_report
  for no in {1..10}; do
    go test -benchmem -bench . ./bench >> $new_llrb_report
  done

  (
    cd $ck
    git checkout 'issue-6465-develop'
  )
  new_btree_report=$report_dir/new_btree
  rm -f $new_btree_report
  for no in {1..10}; do
    go test -benchmem -bench . ./bench -impl btree >> $new_btree_report
  done

  (
    cd $ck
    git checkout 'no-free-list'
  )
  no_free_list_report=$report_dir/no_free_list
  rm -f $no_free_list_report
  for no in {1..10}; do
    go test -benchmem -bench . ./bench -impl btree >> $no_free_list_report
  done
}

if [[ $# -ne 1 ]]; then
  echo "command must be provided"
  exit 1
fi

cmd=$1

fork=$GOPATH/src/github.com/yaojingguo/cockroach
ck=$GOPATH/src/github.com/cockroachdb/cockroach
bench_url='github.com/yaojingguo/benchmark-interval-tree/bench'
btree_based="$bench_url/btree_based"
llrb_based="$bench_url/llrb_based"
report_dir='report'
basic="^Benchmark((Insert)|(FastInsert)|(Delete)|(Get))\$"
no_random='.*Benchmark[^R].*'
random='Random'

case $cmd in
  approach)
    approach;;
  degree)
    degree;;
  opt)
    opt;;
  vs_llrb)
    vs_llrb;;
  random)
    random;;
  new)
    llrb_btree_no_free_list;;
  init)
    rm -fr $report_dir
    mkdir $report_dir;;
esac
