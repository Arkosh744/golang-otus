package main

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestCopyFileWithPredefinedFiles(t *testing.T) {
	type args struct {
		fromPath string
		toPath   string
		offset   int64
		limit    int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success copy file 0 0",
			args: args{
				fromPath: "testdata/input.txt",
				offset:   0,
				limit:    0,
			},
			wantErr: nil,
		},
		{
			name: "success copy file 0 10",
			args: args{
				fromPath: "testdata/input.txt",
				offset:   0,
				limit:    10,
			},
			wantErr: nil,
		},
		{
			name: "success copy file 100 1000",
			args: args{
				fromPath: "testdata/input.txt",
				offset:   100,
				limit:    1000,
			},
			wantErr: nil,
		},
		{
			name: "success copy file 6000 1000",
			args: args{
				fromPath: "testdata/input.txt",
				offset:   6000,
				limit:    1000,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.toPath = fmt.Sprintf("testdata/output_offset%d_limit%d.txt", tt.args.offset, tt.args.limit)
			predefinedResultFile, err := os.Open(fmt.Sprintf("testdata/out_offset%d_limit%d.txt", tt.args.offset, tt.args.limit))
			if err != nil {
				t.Errorf("os.Open() error = %v", err)
			}
			defer predefinedResultFile.Close()

			err = CopyFile(tt.args.fromPath, tt.args.toPath, tt.args.offset, tt.args.limit)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer os.Remove(tt.args.toPath)

			resultTestFile, err := os.Open(tt.args.toPath)
			if err != nil {
				t.Errorf("os.Open() error = %v", err)
			}
			defer resultTestFile.Close()

			if !compareFiles(predefinedResultFile, resultTestFile) {
				t.Errorf("compareFiles() not successful")
			}
		})
	}
}

func compareFiles(file1, file2 *os.File) bool {
	buffer1 := make([]byte, 8)
	buffer2 := make([]byte, 8)
	for {
		n1, err1 := file1.Read(buffer1)
		n2, err2 := file2.Read(buffer2)
		if n1 != n2 || !errors.Is(err1, err2) {
			fmt.Println("n1 != n2", n1, n2)
			return false
		}

		if err1 != nil {
			break
		}
	}

	return true
}

func TestCopyFile(t *testing.T) {
	type args struct {
		fromPath string
		toPath   string
		offset   int64
		limit    int64
	}
	tests := []struct {
		name       string
		args       args
		srcContent []byte
		prepare    func(srcContent []byte) (srcFilePath string, dstFile *os.File)
		after      func(dstFile *os.File, srcFilePath string, srcContent []byte)
		wantErr    error
	}{
		{
			name: "fail offset exceeds file size",
			args: args{
				offset: 99999,
				limit:  0,
			},
			srcContent: []byte("This is a test file for CopyFile function."),
			prepare:    prepareTempFiles(t),
			after:      nil,
			wantErr:    ErrOffsetExceedsFileSize,
		},
		{
			name: "fail file not found",
			args: args{
				offset: 0,
				limit:  0,
			},
			srcContent: []byte("This is a test file for CopyFile function."),
			prepare: func(srcContent []byte) (srcFile string, dstFile *os.File) {
				dstFile, err := os.CreateTemp("", "dstFile")
				if err != nil {
					t.Fatal(err)
				}

				return "some thing.good", dstFile
			},
			after:   nil,
			wantErr: ErrFileNotFound,
		},
		{
			name: "fail undefined length",
			args: args{
				offset: 0,
				limit:  0,
			},
			srcContent: []byte("This is a test file for CopyFile function."),
			prepare: func(srcContent []byte) (srcFile string, dstFile *os.File) {
				dstFile, err := os.CreateTemp("", "dstFile")
				if err != nil {
					t.Fatal(err)
				}

				return "/dev/urandom", dstFile
			},
			after:   nil,
			wantErr: ErrUnsupportedFile,
		},
		{
			name: "success with limit",
			args: args{
				offset: 0,
				limit:  99999,
			},
			srcContent: []byte("This is a test file for CopyFile function."),
			prepare:    prepareTempFiles(t),
			after:      checkContentAndDeleteTempFiles(t),
			wantErr:    nil,
		},
		{
			name: "success copy file",
			args: args{
				offset: 0,
				limit:  0,
			},
			srcContent: []byte("This is a test file for CopyFile function."),
			prepare:    prepareTempFiles(t),
			after:      checkContentAndDeleteTempFiles(t),
			wantErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcFilePath, dstFile := tt.prepare(tt.srcContent)
			defer func() {
				if tt.after != nil {
					tt.after(dstFile, srcFilePath, tt.srcContent)
				}
			}()

			if err := CopyFile(srcFilePath, dstFile.Name(), tt.args.offset, tt.args.limit); !errors.Is(err, tt.wantErr) {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func checkContentAndDeleteTempFiles(t *testing.T) func(dstFile *os.File, srcFilePath string, srcContent []byte) {
	t.Helper()

	return func(dstFile *os.File, srcFilePath string, srcContent []byte) {
		defer func() {
			os.Remove(dstFile.Name())
			os.Remove(srcFilePath)
		}()

		dstContent, err := os.ReadFile(dstFile.Name())
		if err != nil {
			t.Fatal(err)
		}

		if string(dstContent) != string(srcContent) {
			t.Fatalf("Expected: %s, Got: %s", string(srcContent), string(dstContent))
		}
	}
}

func prepareTempFiles(t *testing.T) func(srcContent []byte) (srcFilePath string, dstFile *os.File) {
	t.Helper()

	return func(srcContent []byte) (srcFilePath string, dstFile *os.File) {
		srcFile, err := os.CreateTemp("", "srcFile")
		if err != nil {
			t.Fatal(err)
		}

		dstFile, err = os.CreateTemp("", "dstFile")
		if err != nil {
			t.Fatal(err)
		}

		if _, err = srcFile.Write(srcContent); err != nil {
			t.Fatal(err)
		}

		return srcFile.Name(), dstFile
	}
}
