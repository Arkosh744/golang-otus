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

type copyBufferStruct struct {
	buffer           []byte
	totalBytesCopied int64
	limit            int64
}

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

	bufferStruct := copyBufferStruct{
		buffer:           make([]byte, 8),
		totalBytesCopied: 0,
		limit:            limit,
	}

	bufferStruct.totalBytesCopied, err = processCopy(srcFile, dstFile, bufferStruct)
	if err != nil {
		return err
	}

	fmt.Println("\nCopying finished.")
	fmt.Println("Result file size:", bufferStruct.totalBytesCopied, "bytes")
	fmt.Println("Result file path:", toPath)
	return nil
}

func processCopy(
	srcFile *os.File,
	dstFile *os.File,
	bufferStruct copyBufferStruct,
) (int64, error) {
	for {
		bytesRead, readErr := srcFile.Read(bufferStruct.buffer)
		if readErr != nil && readErr != io.EOF {
			return 0, readErr
		}

		if bytesRead <= 0 {
			break
		}

		bytesToWrite := bytesRead
		if bufferStruct.totalBytesCopied+int64(bytesRead) > limit {
			bytesToWrite = int(limit - bufferStruct.totalBytesCopied)
		}

		bytesWritten, writeErr := dstFile.Write(bufferStruct.buffer[:bytesToWrite])
		if writeErr != nil {
			return 0, writeErr
		}

		bufferStruct.totalBytesCopied += int64(bytesWritten)
		printProgress(bufferStruct.totalBytesCopied, limit)

		if bufferStruct.totalBytesCopied >= limit {
			break
		}

	}
	return bufferStruct.totalBytesCopied, nil
}

func isSpecialMode(fileMode fs.FileMode) bool {
	return fileMode&os.ModeDevice != 0 || fileMode&os.ModeNamedPipe != 0 || fileMode&os.ModeSocket != 0
}

func printProgress(current, total int64) {
	progress := float64(current) / float64(total) * 100
	fmt.Printf("\rCopy progress: %6.2f%%", progress)
}
