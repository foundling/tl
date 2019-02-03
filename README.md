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
$ tl -a "do laundry"
$ tl 
total 1
1: do laundry ✖

...

$ tl -a "triage github issues" -c
$ tl -a "vote"
$ tl
total 2
1: do laundry ✖
2: triage github issues ✓
3: vote ✖

...

$ tl -v
total 3

to do
~~~~~~
1: do laundry ✖
3: vote ✖


complete
~~~~~~~~~
2: triage github issues ✓

```
