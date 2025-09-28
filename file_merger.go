package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func loadIgnorePatterns(projectDir string) ([]string, error) {
    ignoreFile := filepath.Join(projectDir, ".ignore")
    var patterns []string

    data, err := os.ReadFile(ignoreFile)
    if err != nil {
        if os.IsNotExist(err) {
            return patterns, nil // no ignore file, return empty
        }
        return nil, err
    }

    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        patterns = append(patterns, line)
    }
    return patterns, nil
}

func shouldIgnore(relPath string, patterns []string) bool {
    for _, pattern := range patterns {
        // normalize to forward slashes for matching
        relPathUnix := filepath.ToSlash(relPath)

        // match direct file or dir
        if ok, _ := filepath.Match(pattern, relPathUnix); ok {
            return true
        }
        // prefix match for directories (e.g. dir/)
        if strings.HasSuffix(pattern, "/") && strings.HasPrefix(relPathUnix, pattern) {
            return true
        }
        // fallback: match basename
        if ok, _ := filepath.Match(pattern, filepath.Base(relPathUnix)); ok {
            return true
        }
    }
    return false
}

func listFiles(projectDir string, patterns []string) ([]string, error) {
    var files []string
    err := filepath.WalkDir(projectDir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }
        relPath, err := filepath.Rel(projectDir, path)
        if err != nil {
            return err
        }
        if !shouldIgnore(relPath, patterns) {
            files = append(files, relPath)
        }
        return nil
    })
    return files, err
}

func chooseFiles(files []string) ([]string, error) {
    fmt.Println("Available files:")
    for idx, file := range files {
        fmt.Printf("[%d] %s\n", idx, file)
    }

    fmt.Print("Enter the indices of files to include (comma separated): ")
    reader := bufio.NewReader(os.Stdin)
    input, err := reader.ReadString('\n')
    if err != nil {
        return nil, err
    }

    input = strings.TrimSpace(input)
    if input == "*" { // select all files
        return files, nil
    }

    indices := strings.Split(input, ",")
    var chosen []string
    for _, idxStr := range indices {
        idxStr = strings.TrimSpace(idxStr)
        if idx, err := strconv.Atoi(idxStr); err == nil && idx >= 0 && idx < len(files) {
            chosen = append(chosen, files[idx])
        }
    }

    return chosen, nil
}

func writeDataFile(projectDir string, chosenFiles []string, outputFile string) error {
    out, err := os.Create(outputFile)
    if err != nil {
        return err
    }
    defer out.Close()

    writer := bufio.NewWriter(out)
    for _, file := range chosenFiles {
        fmt.Fprintf(writer, "// %s\n", file)
        content, err := os.ReadFile(filepath.Join(projectDir, file))
        if err != nil {
            fmt.Fprintf(writer, "[Error reading %s: %v]\n", file, err)
        } else {
            writer.Write(content)
        }
        writer.WriteString("\n\n")
    }
    return writer.Flush()
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter project directory path: ")
    projectDir, _ := reader.ReadString('\n')
    projectDir = strings.TrimSpace(projectDir)

    fi, err := os.Stat(projectDir)
    if err != nil || !fi.IsDir() {
        fmt.Println("Invalid directory path.")
        return
    }

    patterns, err := loadIgnorePatterns(projectDir)
    if err != nil {
        fmt.Println("Error loading .ignore:", err)
        return
    }

    files, err := listFiles(projectDir, patterns)
    if err != nil {
        fmt.Println("Error listing files:", err)
        return
    }
    if len(files) == 0 {
        fmt.Println("No files found in the given directory.")
        return
    }

    chosen, err := chooseFiles(files)
    if err != nil {
        fmt.Println("Error choosing files:", err)
        return
    }
    if len(chosen) == 0 {
        fmt.Println("No files selected.")
        return
    }

    if err := writeDataFile(projectDir, chosen, "data.txt"); err != nil {
        fmt.Println("Error writing data.txt:", err)
    } else {
        fmt.Println("Data written to data.txt")
    }
}
