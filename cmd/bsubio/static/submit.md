# bsubio submit

Submit a job for processing

## Usage

```
bsubio submit [options] <input_file> <type>
```

## Options

- `-o <file>` - Output file path (requires -w)
- `-w` - Wait for job to complete

## Arguments

- `input_file` - Path to the input file
- `type` - Job type

## Examples

Submit a job:
```
bsubio submit data.json json_format
```

Submit and wait for completion:
```
bsubio submit -w input.txt passthrough
```

Submit, wait, and save output to file:
```
bsubio submit -w -o result.txt input.txt passthrough
```
