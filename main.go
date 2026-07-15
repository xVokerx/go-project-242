package code

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"
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
		return fmt.Sprintf("%s \t %s", humanize(size), path), nil
	}
	return fmt.Sprintf("%dB \t %s", size, path), nil
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
	units := []string{"B", "KB", "MB", "GB", "TB"}
	i := 0
	sizeInFloat := float64(size)
	for sizeInFloat >= 1024 && i < len(units)-1 {
		sizeInFloat /= 1024
		i++
	}
	return strconv.FormatFloat(sizeInFloat, 'f', 1, 64) + units[i]
}

func main() {
	cmd := &cli.Command{
		Usage:     "print size of a file or directory; supports -r (recursive), -H (human-readable), -a (include hidden)",
		Name:      "hexlet-path-size",
		UsageText: "hexlet-path-size [global options] <path>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "human",
				Aliases: []string{"H"},
				Usage:   "human-readable sizes (auto-select unit) (default: false)",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "include hidden files and directories (default: false)",
			},
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "recursive size of directories (default: false)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.NArg() < 1 {
				return fmt.Errorf("missing <path> argument")
			}
			path := cmd.Args().Get(0)
			isHuman := cmd.Bool("human")
			isAll := cmd.Bool("all")
			isRecursive := cmd.Bool("recursive")
			file, err := GetPathSize(path, isHuman, isAll, isRecursive)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(file)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
	}

}
