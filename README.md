# tl

`tl` is a minimal command-line task list application written in Go.

## Installation

`go get -v github.com/foundling/tl`

## Usage

```
Usage of tl:
  -a string
    	task text to append
  -c	toggle task complete status
  -d string
    	task number to delete
  -f string
    	alternate task data filepath to ~/tl.csv (default "/Users/alex/tl.csv")
  -h	usage
  -p string
    	task to prepend
  -t string
    	task update text
  -u int
    	task number to update (default -1)
  -v	print verbose information
```

## `tl` In Action

[![demo](https://asciinema.org/a/229292.svg)](https://asciinema.org/a/229292?autoplay=1)
