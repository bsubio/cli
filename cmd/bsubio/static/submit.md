# bsubio submit

Submit a job for processing

## Usage

```
bsubio submit [options] <type> <input_file> [<input_file2> ...]
```

## Options

- `-o <file>` - Output file path (requires -w, cannot be used with multiple files)
- `-w` - Wait for job to complete

## Arguments

- `type` - Job type
- `input_file` - Path to the input file(s)

## Examples

Submit a job:

```
bsubio submit json_format data.json
```

Submit multiple files:

```
bsubio submit json_format file1.json file2.json file3.json
```

Submit and wait for completion:

```
bsubio submit -w passthru input.txt
```

Submit multiple files, wait, and print all outputs to stdout:

```
bsubio submit -w passthru *.pdf
```

Submit, wait, and save output to file (single file only):

```
bsubio submit -w -o result.txt passthru input.txt
```
