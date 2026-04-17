# ✂️ cutter

Cutter is a simple cli program for cutting large text/csv files into chunks that are digestible by Excel.

## How it Works

Excel can only open files that have approximately a million rows or less. Files bigger than that fail to open.

Cutter solves this problem by:

- Splitting a large text file into numbered chunks, each a million rows or less
- Preserving the column headers on each chunk

## Why not just use Split?

What makes `cutter` different from the POSIX `split` command?

- It automatically defaults to Excel friendly defaults
- No need to futz around with command line arguments
- File chunks are automatically renamed to the original file name, followed by `_part1`, `_part2`, instead of `xaa`, `xbb` and etc..
- Drag and drop works in windows

## How to use it?

1. Save `cutter.exe` into a folder
2. Drag and drop the file you want to split onto `cutter.exe`

The split chunks will be safely deposited in the folder from which you dragged over the file.

## CLI Usage

If you want to use the terminal:
  
    Usage: cutter.exe [options] <filename>
        Options:
          -v, --version    Print version information and exit
          -h, --help       Print this message and exit

## Installation

### Multi-Platform:

You can install gndn using `go install`:

```bash
go install github.com/maciakl/cutter@latest
```

### macOS and Linux:

Use [grab](https://github.com/maciakl/grab):

```bash
grab maciakl/cutter
```

### Windows:

Use [scoop](https://scoop.sh):

```bash
scoop add maciak https://github.com/maciakl/bucket
scoop update
scoop install cutter
```
