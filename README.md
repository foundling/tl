# tl

tl is a simple command-line tasklist application written in Go.

## Installation

`go get github.com/foundling/tl`

## Usage

```
tl [-aduv]
  -v 
    print tasks in verbose format
  -a <text>
  -u <task #> [-t updated task text] [-c]
    update task by task #
    -t <task text>
    -c
      mark complete
  -d <task #>
    delete task by task #
```

## Examples

```
# add a task 
$ tl -a "do laundry"
$ tl 
total 1
1: do laundry ✖

...

# add a task and mark it complete 
$ tl -a "triage github issues" -c
$ tl -a "vote"

...

# print tasks
$ tl
total 2
1: do laundry ✖
2: triage github issues ✓
3: vote ✖

...

# print tasks in verbose format
$ tl -v
total 3

to do
~~~~~~
1: do laundry ✖
3: vote ✖


complete
~~~~~~~~~
2: triage github issues ✓

...

# update a task
$ tl -u 1 -t "do laundry (wash new jeans separately)"
$ tl
total 3
1: do laundry (wash new jeans separately) ✖
2: triage github issues ✓
3: vote ✖
```
