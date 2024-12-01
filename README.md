# YAM - Yet Another Merger

**YAM**YAM is a command-line tool for merging multiple YAML files with support for nested structures, arrays, and customizable conflict resolution.

## Installation

Install the latest version with:

```bash
go install github.com/darmiel/yamyam@latest
```

## Usage

Merge YAML files with flexible options:

```bash
yamyam [flags] <file1.yaml> <file2.yaml> ...
```

### Flags

- `-overwrite`: Overwrite duplicate keys during merging.
- `-skip`: Skip duplicate keys instead of throwing errors.
- `-skipKey <key>`: Skip specific keys. Can be used multiple times.
- `-output <file>`: Write the result to a file. Default: stdout.
- `-verbose`: Enable more verbose output.

### Examples

1. Merge files and overwrite duplicates:

   ```bash
   yamyam -overwrite file1.yaml file2.yaml
   ```

2. Skip duplicates and write output to a file:

   ```bash
   yamyam -skip -output merged.yaml file1.yaml file2.yaml
   ```

3. Skip specific keys during merging:

   ```bash
   yamyam -skipKey "apiKey" -skipKey "debug" file1.yaml file2.yaml
   ```

4. Enable verbose logging:

   ```bash
   yamyam -verbose file1.yaml file2.yaml
   ```

## Requirements

- Go 1.20 or later
