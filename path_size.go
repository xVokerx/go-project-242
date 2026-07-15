package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GetPathSize(path string, isHuman bool, isAll bool, isRecursive bool) (string, error) {
	var size int64
	file, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if !file.IsDir() {
		size = file.Size()
	} else {
		dirEntry, err := os.ReadDir(path)
		if err != nil {
			return "", err
		}
		for _, file := range dirEntry {
			if !isAll && strings.HasPrefix(file.Name(), ".") {
				continue
			}
			if file.IsDir() {
				if isRecursive {
					dirSize, err := recursiveDirSize(filepath.Join(path, file.Name()), isAll)
					if err != nil {
						return "", err
					}
					size += dirSize
				}
				continue
			}
			fileinfo, err := file.Info()
			if err != nil {
				return "", err
			}
			size += fileinfo.Size()
		}
	}
	if isHuman {
		return fmt.Sprintf("%s\t%s", humanize(size), path), nil
	}
	return fmt.Sprintf("%dB\t%s", size, path), nil
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
	units := []string{"KB", "MB", "GB", "TB"}
	i := 0
	sizeInFloat := float64(size)
	for sizeInFloat > 1024 && i < len(units)-1 {
		sizeInFloat /= 1024
		i++
	}
	return strconv.FormatFloat(sizeInFloat, 'f', 1, 64) + units[i]
}
