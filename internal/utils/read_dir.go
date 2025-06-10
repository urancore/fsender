package utils

import (
	"os"
)

func ReadDir(path string) ([]os.FileInfo, error) {
	entries, err := os.ReadDir(path) // Получаем []fs.DirEntry
	if err != nil {
		return nil, err
	}

	infos := make([]os.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info() // Получаем FileInfo из DirEntry
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}

	return infos, nil
}
