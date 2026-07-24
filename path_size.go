package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GetPathSize(path string, isRecursive bool, isHuman bool, isAll bool) (string, error) {
	var size int64
	file, err := os.Lstat(path)
	if err != nil {
		return "", fmt.Errorf("read directory %s: %w", path, err)
	}
	if !file.IsDir() {
		size = file.Size()
	} else {
		dirEntry, err := os.ReadDir(path)
		if err != nil {
			return "", fmt.Errorf("read directory %s: %w", path, err)
		}
		for _, file := range dirEntry {
			if !isAll && strings.HasPrefix(file.Name(), ".") {
				continue
			}
			fullPath := filepath.Join(path, file.Name())
			if file.IsDir() {
				if isRecursive {
					dirSize, err := recursiveDirSize(fullPath, isAll)
					if err != nil {
						fmt.Fprintf(os.Stderr, "warning: %v\n", fmt.Errorf("read dir %s: %w", fullPath, err))
						continue
					}
					size += dirSize
				}
				continue
			}
			fileinfo, err := file.Info()
			if err != nil {
				return "", fmt.Errorf("read directory %s: %w", path, err)
			}
			size += fileinfo.Size()
		}
	}
	if isHuman {
		return humanize(size), nil
	}
	return strconv.FormatInt(size, 10) + "B", nil
}

func recursiveDirSize(path string, isAll bool) (int64, error) {
	var size int64
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return 0, fmt.Errorf("read directory %s: %w", path, err)
	}
	for _, file := range dirEntry {
		if !isAll && strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.IsDir() {
			fullPath := filepath.Join(path, file.Name())
			nestedSize, err := recursiveDirSize(fullPath, isAll)
			if err != nil {
				return 0, fmt.Errorf("read directory %s: %w", fullPath, err)
			}
			size += nestedSize
		} else {
			fileInfo, err := file.Info()
			if err != nil {
				return 0, fmt.Errorf("read file %s: %w", file.Name(), err)
			}
			size += fileInfo.Size()
		}
	}
	return size, nil
}

func humanize(size int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	i := 0
	sizeInFloat := float64(size)
	for sizeInFloat >= 1024 && i < len(units)-1 {
		sizeInFloat /= 1024
		i++
	}
	if i == 0 {
		return strconv.FormatInt(size, 10) + units[i]
	}
	return strconv.FormatFloat(sizeInFloat, 'f', 1, 64) + units[i]
}
