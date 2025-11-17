# bsubio submit

Submit one or more jobs for processing

## Usage

```
bsubio submit [options] <type> <input_file> [input_file2 ...]
```

## Options

- `-o <file>` - Output file path (requires -w, only works with single file)
- `-w` - Wait for all jobs to complete

## Arguments

- `type` - Job type
- `input_file` - Path to input file(s) (can specify multiple)

## Examples

Submit a single job:

```
bsubio submit json_format data.json
```

Submit multiple jobs:

```
bsubio submit json_format file1.json file2.json file3.json
```

Submit multiple jobs using glob patterns:

```
bsubio submit passthru *.pdf
```

Submit and wait for all jobs to complete:

```
bsubio submit -w passthru input1.txt input2.txt
```

Submit single job, wait, and save output to file:

```
bsubio submit -w -o result.txt passthru input.txt
```

## Notes

- When multiple files are specified, one job is submitted per file
- With `-w` flag, waits for all jobs to complete and prints outputs in the order files were provided
- The `-o` flag cannot be used with multiple input files
