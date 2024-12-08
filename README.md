# 1brc without concurrency

Just another solution to the 1 billion row challenge but without concurrency.
Reasonably all solutions in Go that I know use concurrency.
But why not trying out what's possible without directly using concurrency?
True to the motto: concurrency should only be used when all other options have been exhausted.

Of course does Go parallelize things in the background and thus this solution is not 100% single threaded.
To achieve this we would need to do something like compiling to WASM since WASM enforces single threaded execution.
Here's an interesting post about determinism in Go: https://www.polarsignals.com/blog/posts/2024/05/28/mostly-dst-in-go

The rules of 1brc can be found here [https://github.com/gunnarmorling/1brc](https://github.com/gunnarmorling/1brc).

For comparison please note that I ran all my tests on an Apple M1 Pro.

## Project

Every `run_X` package contains one solution. Each ascending folder contains another change.
Tests and Benchmarks can be found in `run_test.go` and in the respective `run_X` packages.

I still did one implementation with concurrency in "concurrent_1". 
Though this version doesn't contain all improvements I made along the way.

## Changelog

### run1 - 119s

Folder: https://github.com/Erik7354/1brc-go/tree/main/run_1

Run 1 is just the first and naive implementation without any optimizations in mind.

### run2 - 108s

Folder: https://github.com/Erik7354/1brc-go/tree/main/run_2

Run 2 does a little optimization. Each line is now not split anymore using `strings.Split` but `strings.SplitN`. 
The latter only runs till it split the string into maximum `n := 2` substrings. 
This way we save us scanning half of the line. The performance gains aren't phenomenal but noticeable.

### run3 - 82s

Folder: https://github.com/Erik7354/1brc-go/tree/main/run_3

Run 3 switched from using `scanner.Text()` to `scanner.Bytes()`. 
Since strings are immutable in Go the first allocates a new string for every call made. 
The latter in contrast reuses the same internal buffer halving the allocations from >2.000.000.000 to >1.000.000.000.
`bytes.Cut(...)` is now used to split the lines. It only splits once, comparable to `strings.SplitN` with `n := 2`.

### run4 - 61s

Folder: https://github.com/Erik7354/1brc-go/tree/main/run_4

Run 4 replaces floats with ints. 
The new func `lineToInt(bs []byte)` takes exactly the bytes of the temperature and parses them to an int.
One of the constraints of the challenge is that the temperature value is a "non-null double between -99.9 (inclusive) and 99.9 (inclusive), always with one fractional digit".
That said the integer part of the value is just multiplied by 10 and the decimal part becomes the ones place.
Only at the end when printing the results the ints are cast back to floats.

### run5 - 56s

Folder: https://github.com/Erik7354/1brc-go/tree/main/run_5

Run 5 again is just a small improvement saving a few seconds.
Here I now use a custom func `splitLine` instead of `bytes.Cut`.

```Go
// splitLine
// "München;23.5" => "München", []byte("23.5")
// "München;-10.5" => "München", []byte("-10.5")
func splitLine(line []byte) (string, []byte) {
	l := len(line)

	switch {
	case line[l-4] == ';': // 1.2
		return string(line[:l-4]), line[l-3:]
	case line[l-5] == ';': // 12.3 or -1.2
		return string(line[:l-5]), line[l-4:]
	case line[l-6] == ';': // -12.3
		return string(line[:l-6]), line[l-5:]
	default:
		panic("unknown format")
	}
}
```

### run6 - 54s

Folder: https://github.com/Erik7354/1brc-go/tree/main/run_6

Run 6 is even smaller.
I migrated `lineToInt` from using a for-loop to using if-else.
Benchmarks can be found in the directory.

### run7 - 48s

Folder: https://github.com/Erik7354/1brc-go/tree/main/run_7

Run 7 is at least a bigger code change and again a small performance gain.
I benchmarked `bufio.Scanner` against `file.Read` with a small and a bigger buffer and noticed how much faster read was.
This makes sense since the scanner needs to look at every byte and find the right place to split. Read just reads.
Now I can't just use file read and read every byte by byte - this is also very slow.
It's clear I need some kind of buffering but with `bufio.Scanner` a lot of bytes are processed twice leading to worse performance.
So I created my own `StationScanner` that buffers but also tries to process every byte only once.

The most interesting func of `StationScanner` is probably `updateChunk` which checks if there are still enough unprocessed bytes in `chunk` that it could hold the longest line possible.
If not then it copies the remaining unprocessed one to the start of `chunk` and reads in new bytes.
That way only one buffer is needed for the whole work.

```Go
type StationScanner struct {
	f *os.File

	chunk [chunkSize]byte
	start int
	end   int

	eof bool
}

func (s *StationScanner) updateChunk() {
    if s.end-s.start >= maxLineLength && !s.eof {
        return // still at least one whole line left in s.chunk
    }
    
    copy(s.chunk[:], s.chunk[s.start:s.end])
    s.end = s.end - s.start
    s.start = 0
    
    n, err := s.f.Read(s.chunk[s.end:])
    if err == io.EOF {
        s.eof = true
    }
    if err != nil && err != io.EOF {
        panic(err)
    }
    s.end += n
}

// ...other functions
```

## run8 - 39s

Folder: https://github.com/Erik7354/1brc-go/tree/main/run_8

Run8 gets serious and uses `unsafe.String(unsafe.SliceData(bName), len(bName))` to "convert" the name-bytes to a string.

This unfortunately is not directly possible since the underlying buffer is `StationScanner.chunk`.
As written in run7 this buffer is volatile and so can't serve as storage for static strings.

First and naive we could make a copy of the name-slice to a new static buffer, like below.
But in doing so we essentially mimic what `string(...)` is probably doing under the hood - creating a new immutable buffer for the string.
Mimicking this doesn't make much sense and neither would it improve performance.
```Go
bName := make([]byte, len(lines[:l]))
copy(bName, lines[:l])
name = unsafe.String(unsafe.SliceData(bName), len(bName))
```

The other option is to make regular copies (backups if you will) of the whole `StationScanner.chunk`.
That way `bName` still holds a valid memory reference, and we win some performance because we just do one memcopy per chunk and not per line.
Speaking facts, this change reduces our allocations from >1 billion to >2900 and our ns/op by about 20%.
```Go
// func updateChunk()
backup := s.chunk
s.chunk = make([]byte, chunkSize)
copy(s.chunk, backup)


// func Line()
bName := lines[:l]
name = unsafe.String(unsafe.SliceData(bName), len(bName))
```

### run9 - 29s

Folder: https://github.com/Erik7354/1brc-go/tree/main/run_9

Run9 using a custom `indexByte` func instead of `bytes.IndexByte` decreased the ns/op by another 25%.

```Go
func indexByte(b []byte, c byte) int {
	for i, bb := range b {
		if bb == c {
			return i
		}
	}

	panic("not a line or end reached")
}
```

But 25% looked too good to be true for me (for such a simple functionality) so I benchmark both version against each other. 
Surprisingly, the custom function takes about 2x the time of the stdlib func to complete...on short strings. 
On long strings the time difference is 10x.

Tbh I'm missing a good explanation for this. 
I first benchmarked the long slice and thought maybe the custom func is just faster on short ones.
This way the final benchmark would make sense since the final dataset we use mostly has short station names.
Of course this thesis proved false.
My next thought was _maybe_ it's because I import an additional library `bytes` - though this already seemed unlikely to me.
Importing `_ "bytes"` didn't change anything of the final benchmark so some kind of `init` side effects weren't the reason either.

## Further Ideas

For further improvements I think the biggest leverage is using more unsafe Go or a custom map.

When reading the CPU profiler for [run9](Folder: https://github.com/Erik7354/1brc-go/tree/main/run_9) about 35% of the total CPU time is `runtime.mapaccess2_faststr`.
Map-access-2 only happens in line 47 in run9 `c, ok := stations[name]`.
It could be tried to improve the percentage using a specialised custom map.
Since we only have at max 10.000 unique keys a simple hash algorithm like FNV would be sufficient and maybe more performant.
Maybe even better would be to try [xxHash](https://xxhash.com/).
Secondly (not directly related to maps) it could be tried to not use `string` but `[]byte` as map key - though that's not possible with the standard map.
That way it wouldn't be necessary to cast the bytes to strings in the first place.

I also often read about memory mapping `mmap` to increase performance, but I'm not sure about that - though I could easily be wrong.
Let me give my two cents.
Mmap is used to map files bigger than physical memory to virtual memory to provide random access to programs.
As far as I know this is widely used in for example databases where you need to read from many different positions from big files.
For this, I assume, mmap has to use `file.Read` (or some variation of that) under the hood. 
Now with 1brc we read the measurements only sequentially and not randomly, at least from what I have seen.
So wouldn't using mmap just add a layer on top of a plain read?

