# waitgroupext [![Build Status](https://travis-ci.org/abligh/waitgroupext.svg?branch=master)](https://travis-ci.org/abligh/waitgroupext) [![GoDoc](http://godoc.org/github.com/abligh/waitgroupext?status.png)](http://godoc.org/github.com/abligh/waitgroupext) [![GitHub release](https://img.shields.io/github/release/abligh/waitgroupext.svg)](https://github.com/abligh/waitgroupext/releases)

`waitgroupext` provides extended WaitGroups that are similar in usage
and API to `sync.WaitGroup` save for the presence of a `WaitChan()`
function. This returns a channel that can be selected upon to wait
for a WaitGroup, allowing an interruptible wait.

Benchmarks are a little slower than conventional waitgroups:

*Conventional Waitgroups*
```
$ go test -bench WaitGroup
PASS
BenchmarkWaitGroupUncontended-4		100000000		15.4 ns/op
BenchmarkWaitGroupAddDone-4     	30000000		50.2 ns/op
BenchmarkWaitGroupAddDoneWork-4 	30000000		60.7 ns/op
BenchmarkWaitGroupWait-4        	1000000000		2.97 ns/op
BenchmarkWaitGroupWaitWork-4    	30000000		46.9 ns/op
BenchmarkWaitGroupActuallyWait-4	 5000000		235 ns/op
	16 B/op        1 allocs/op
ok  	sync	       11.662s
```

*waitgroupext WaitGroups*

```
$ go test -bench .
PASS
BenchmarkWaitGroupUncontended-4		 5000000		328 ns/op
BenchmarkWaitGroupAddDone-4     	 3000000	       	428 ns/op
BenchmarkWaitGroupAddDoneWork-4 	 3000000	       	562 ns/op
BenchmarkWaitGroupWait-4        	 5000000	       	268 ns/op
BenchmarkWaitGroupWaitWork-4    	 5000000		301 ns/op
BenchmarkWaitGroupActuallyWait-4	 3000000		509 ns/op
	128 B/op        2 allocs/op
ok  	github.com/abligh/waitgroupext	  12.559s
```

