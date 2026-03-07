// Copyright (c) 2019-2026, Sylabs Inc. All rights reserved.
// Copyright (c) 2026, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathExists(t *testing.T) {
	dir := t.TempDir()
	if !PathExists(dir) {
		t.Error("PathExists(tempDir) = false, want true")
	}
	if PathExists(filepath.Join(dir, "nonexistent")) {
		t.Error("PathExists(nonexistent) = true, want false")
	}
	if PathExists("") {
		t.Error("PathExists(empty string) = true, want false")
	}
}

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file")
	if err := os.WriteFile(f, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if !FileExists(f) {
		t.Error("FileExists(regularFile) = false, want true")
	}
	if FileExists(dir) {
		t.Error("FileExists(directory) = true, want false")
	}
	if FileExists(filepath.Join(dir, "nonexistent")) {
		t.Error("FileExists(nonexistent) = true, want false")
	}
}

func TestIsDir(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file")
	if err := os.WriteFile(f, nil, 0644); err != nil {
		t.Fatal(err)
	}
	if !IsDir(dir) {
		t.Error("IsDir(dir) = false, want true")
	}
	if IsDir(f) {
		t.Error("IsDir(file) = true, want false")
	}
	if IsDir(filepath.Join(dir, "nonexistent")) {
		t.Error("IsDir(nonexistent) = true, want false")
	}
}

func TestGetTarArgs(t *testing.T) {
	tests := []struct {
		format string
		want   string
	}{
		{FormatBZ2, "-xjf"},
		{FormatGZ, "-xzf"},
		{FormatTGZ, "-xzf"},
		{FormatTAR, "-xf"},
		{"unknown", ""},
		{"", ""},
	}
	for _, tt := range tests {
		if got := GetTarArgs(tt.format); got != tt.want {
			t.Errorf("GetTarArgs(%q) = %q, want %q", tt.format, got, tt.want)
		}
	}
}

func TestDetectTarballFormat(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"archive.tar.bz2", FormatBZ2},
		{"archive.bz2", FormatBZ2},
		{"archive.tar.gz", FormatGZ},
		{"archive.gz", FormatGZ},
		{"archive.tgz", FormatTGZ},
		{"archive.tar", FormatTAR},
		{"archive.zip", UnknownFormat},
		{"archive", UnknownFormat},
		{"", UnknownFormat},
	}
	for _, tt := range tests {
		if got := DetectTarballFormat(tt.path); got != tt.want {
			t.Errorf("DetectTarballFormat(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dst := filepath.Join(dir, "dst")
	content := []byte("hello world")
	if err := os.WriteFile(src, content, 0644); err != nil {
		t.Fatal(err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() error = %v", err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("CopyFile: read %q, want %q", got, content)
	}

	// Mode is preserved (e.g. executable bit)
	srcMode := os.FileMode(0750)
	if err := os.Chmod(src, srcMode); err != nil {
		t.Fatal(err)
	}
	dst2 := filepath.Join(dir, "dst2")
	if err := CopyFile(src, dst2); err != nil {
		t.Fatalf("CopyFile() error = %v", err)
	}
	info, err := os.Stat(dst2)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode() != srcMode {
		t.Errorf("CopyFile: dst mode %v, want %v", info.Mode(), srcMode)
	}

	// Same src and dst is a no-op
	if err := CopyFile(src, src); err != nil {
		t.Errorf("CopyFile(src, src) error = %v", err)
	}
}

func TestCopyFileErrors(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "existing")
	if err := os.WriteFile(existing, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("empty source path", func(t *testing.T) {
		err := CopyFile("", filepath.Join(dir, "dst"))
		if err == nil {
			t.Fatal("CopyFile(empty src) expected error")
		}
	})
	t.Run("empty destination path", func(t *testing.T) {
		err := CopyFile(existing, "")
		if err == nil {
			t.Fatal("CopyFile(empty dst) expected error")
		}
	})
	t.Run("source does not exist", func(t *testing.T) {
		err := CopyFile(filepath.Join(dir, "nosuch"), filepath.Join(dir, "dst"))
		if err == nil {
			t.Fatal("CopyFile(nonexistent) expected error")
		}
	})

	t.Run("source is directory", func(t *testing.T) {
		subdir := filepath.Join(dir, "subdir")
		if err := os.Mkdir(subdir, 0755); err != nil {
			t.Fatal(err)
		}
		err := CopyFile(subdir, filepath.Join(dir, "dst"))
		if err == nil {
			t.Fatal("CopyFile(directory) expected error")
		}
	})

	t.Run("destination is directory", func(t *testing.T) {
		dstdir := filepath.Join(dir, "dstdir")
		if err := os.Mkdir(dstdir, 0755); err != nil {
			t.Fatal(err)
		}
		err := CopyFile(existing, dstdir)
		if err == nil {
			t.Fatal("CopyFile(dst=directory) expected error")
		}
	})
}

func TestDirInit(t *testing.T) {
	if err := DirInit(""); err == nil {
		t.Fatal("DirInit(\"\") expected error")
	}

	dir := t.TempDir()
	target := filepath.Join(dir, "initme")
	if err := os.MkdirAll(target, 0755); err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(target, "file")
	if err := os.WriteFile(f, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := DirInit(target); err != nil {
		t.Fatalf("DirInit() error = %v", err)
	}
	if !IsDir(target) {
		t.Error("DirInit: target is not a directory after init")
	}
	if PathExists(f) {
		t.Error("DirInit: old file still exists, expected empty dir")
	}

	// Idempotent: init again
	if err := DirInit(target); err != nil {
		t.Fatalf("DirInit() second call error = %v", err)
	}
}

func TestOpenResultsFile(t *testing.T) {
	dir := t.TempDir()
	fpath := filepath.Join(dir, "results")
	f, err := OpenResultsFile(fpath)
	if err != nil {
		t.Fatalf("OpenResultsFile() error: %v", err)
	}
	f.Close()
	// New file is created with 0644 (not executable)
	info, err := os.Stat(fpath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0644 {
		t.Errorf("new results file perm %v, want 0644", info.Mode().Perm())
	}
	// Open again (append mode)
	f2, err := OpenResultsFile(fpath)
	if err != nil {
		t.Fatalf("OpenResultsFile() second open error: %v", err)
	}
	f2.Close()
}

func TestOpenResultsFileError(t *testing.T) {
	if _, err := OpenResultsFile(""); err == nil {
		t.Fatal("OpenResultsFile(\"\") expected error")
	}
	// Parent directory does not exist -> open must fail
	_, err := OpenResultsFile("/nonexistent/parent/results")
	if err == nil {
		t.Fatal("OpenResultsFile(nonexistent parent) expected error")
	}
}

func TestOpenLogFile(t *testing.T) {
	dir := t.TempDir()
	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if _, err := OpenLogFile("", "name"); err == nil {
		t.Error("OpenLogFile(\"\", \"name\") should return error")
	}
	if _, err := OpenLogFile("prefix", ""); err == nil {
		t.Error("OpenLogFile(\"prefix\", \"\") should return error")
	}
	f, err := OpenLogFile("prefix", "mylog")
	if err != nil {
		t.Fatalf("OpenLogFile(prefix, mylog) error: %v", err)
	}
	f.Close()
	if !FileExists("prefix-mylog.log") {
		t.Error("OpenLogFile did not create expected file")
	}
	info, err := os.Stat("prefix-mylog.log")
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0644 {
		t.Errorf("new log file perm %v, want 0644", info.Mode().Perm())
	}
}
