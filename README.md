# tl

`tl` is a minimal command-line task list application written in Go.

## Installation

`go get -v github.com/foundling/tl`

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

## `tl` In Action

[![demo](https://asciinema.org/a/229292.svg)](https://asciinema.org/a/229292?autoplay=1)
