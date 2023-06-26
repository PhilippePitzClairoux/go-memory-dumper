# Golang memory reader
This projects lets you dump the content of a linux memory file and outputs it in the console.
You can pipe the output of the program to store it in a file if you want.

## Try it out yourself
```bash
git clone https://github.com/PhilippePitzClairoux/go-memory-dumper \
    && cd go-memory-dumper \
    && go build cmd/MemoryDumper.go
```

## Usable Arguments
```bash
x@fedora ~/G/memory-reader> ./MemoryDumper --help
Usage of ./MemoryDumper:
  -filter string
        dump all addresses for pids that match filter
  -hide-empty
        hide lines that are completely empty
  -pid string
        pid of process to dump (default "0")
  -read-loop
        read continuously the data stored in memory
```