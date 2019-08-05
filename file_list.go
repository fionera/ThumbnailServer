package main

import (
    "os"
    "path/filepath"
)

type FileList []os.FileInfo

func NewFileList(path string) (error, *FileList) {
    files := FileList{}

    err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
        if f.IsDir() {
            return nil
        }

        files = append(files, f)

        return nil
    })

    if err != nil {
        return err, nil
    }

    return nil, &files
}
