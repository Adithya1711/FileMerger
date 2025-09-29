# File Merger Tool

A lightweight command-line tool written in Go that helps you collect file contents from a project directory into a single `data.txt` file. You can interactively choose which files to include, apply ignore rules, and even use flexible selection syntax.

---

## Features

- **Recursive file listing** – Lists all files under a project directory.
- **Ignore support** – Skips files and directories defined in a `.ignore` file (similar to `.gitignore`).
- **Negation rules** – Use `!pattern` in `.ignore` to explicitly re-include certain files.
- **Selection syntax** –

  - `0,3,5` → choose specific files by index
  - `*` → choose all files
  - `* !1,2` → choose all files except indices 1 and 2
  - `* !1-3` → choose all files except indices 1 through 3 (range support)

- **Output format** – Writes selected files into `data.txt` with clear headers:

  ```
  // path/to/file.ext
  file content here...

  // another/file.ext
  more content...
  ```

---

## Installation

1. Make sure you have [Go installed](https://go.dev/dl/).
2. Clone this repository or copy the source file.
3. Build the binary:

   ```bash
   go build -o filemerger main.go
   ```

---

## Usage

Run the tool from the command line:

```bash
./filemerger
```

You will be prompted to enter a project directory path. The tool will then:

1. List all files (excluding those ignored).
2. Ask you which files to include using the selection syntax.
3. Write the chosen files into `data.txt` in the current working directory.

---

## The `.ignore` File

Place a `.ignore` file in the root of your project directory to control which files and directories should be skipped.

Examples:

```gitignore
# Ignore all log files
*.log

# Ignore a directory
build/

# Ignore all .exe files
*.exe

# Re-include one specific file
!important.log
!build/config.yaml
```

---

## Examples

### Select all files

```
*
```

### Select specific files

```
0,2,5
```

### Select all except some

```
* !1,3
```

### Select all except a range

```
* !2-5
```

---

## Output Example

If you selected two files (`main.go` and `README.md`), your `data.txt` will look like this:

```text
// main.go
package main

func main() {
    println("Hello, world")
}

// README.md
# Project Title
Some description here.
```

---

## License

This project is released under the MIT License. You are free to use, modify, and distribute it.
