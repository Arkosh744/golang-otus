package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNotFound          = errors.New("file not found")
	ErrFileAlreadyExists     = errors.New("file already exists")
)

func CopyFile(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotFound
		}

		return err
	}

	defer srcFile.Close()

	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	srcFileSize := srcFileInfo.Size()
	if srcFileInfo.Size() < 0 || isSpecialMode(srcFileInfo.Mode()) {
		return ErrUnsupportedFile
	}

	if offset > srcFileSize {
		return ErrOffsetExceedsFileSize
	}

	if offset > 0 {
		if _, err = srcFile.Seek(offset, io.SeekStart); err != nil {
			return err
		}
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return ErrFileAlreadyExists
	}
	defer dstFile.Close()

	if limit == 0 || limit > srcFileSize-offset {
		limit = srcFileSize - offset
	}

	buffer := make([]byte, 8)
	var totalBytesCopied int64

	totalBytesCopied, err = processCopy(srcFile, buffer, totalBytesCopied, limit, dstFile)
	if err != nil {
		return err
	}

	fmt.Println("\nCopying finished.")
	fmt.Println("Result file size:", totalBytesCopied, "bytes")
	fmt.Println("Result file path:", toPath)
	return nil
}

func processCopy(
	srcFile *os.File,
	buffer []byte,
	totalBytesCopied int64,
	limit int64,
	dstFile *os.File,
) (int64, error) {
	for {
		bytesRead, readErr := srcFile.Read(buffer)
		if readErr != nil && readErr != io.EOF {
			return 0, readErr
		}

		if bytesRead > 0 {
			bytesToWrite := bytesRead
			if totalBytesCopied+int64(bytesRead) > limit {
				bytesToWrite = int(limit - totalBytesCopied)
			}

			bytesWritten, writeErr := dstFile.Write(buffer[:bytesToWrite])
			if writeErr != nil {
				return 0, writeErr
			}

			totalBytesCopied += int64(bytesWritten)
			printProgress(totalBytesCopied, limit)

			if totalBytesCopied >= limit {
				break
			}
		} else {
			break
		}
	}
	return totalBytesCopied, nil
}

func isSpecialMode(fileMode fs.FileMode) bool {
	return fileMode&os.ModeDevice != 0 || fileMode&os.ModeNamedPipe != 0 || fileMode&os.ModeSocket != 0
}

func printProgress(current, total int64) {
	progress := float64(current) / float64(total) * 100
	fmt.Printf("\rCopy progress: %6.2f%%", progress)
}
