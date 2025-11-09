# bsubio submit

Submit a job for processing

## Usage

```
bsubio submit [options] <type> <input_file> [<input_file2> ...]
```

## Options

- `-o <file>` - Output file path (requires -w)
- `-w` - Wait for job to complete

## Arguments

- `type` - Job type
- `input_file` - Path to the input file

## Examples

Submit a job:

```
bsubio submit json_format data.json
```

Submit and wait for completion:

```
bsubio submit -w passthrought input.txt
```

Submit, wait, and save output to file:

```
bsubio submit -w -o result.txt passthrough input.txt
```
