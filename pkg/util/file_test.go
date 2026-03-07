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
}

func TestDirInit(t *testing.T) {
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
	f := OpenResultsFile(fpath)
	if f == nil {
		t.Fatal("OpenResultsFile() returned nil")
	}
	f.Close()
	// Open again (append mode)
	f2 := OpenResultsFile(fpath)
	if f2 == nil {
		t.Fatal("OpenResultsFile() second open returned nil")
	}
	f2.Close()
}

func TestOpenLogFile(t *testing.T) {
	dir := t.TempDir()
	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if got := OpenLogFile("", "name"); got != nil {
		t.Error("OpenLogFile(\"\", \"name\") should return nil")
	}
	if got := OpenLogFile("prefix", ""); got != nil {
		t.Error("OpenLogFile(\"prefix\", \"\") should return nil")
	}
	f := OpenLogFile("prefix", "mylog")
	if f == nil {
		t.Fatal("OpenLogFile(prefix, mylog) returned nil")
	}
	f.Close()
	if !FileExists("prefix-mylog.log") {
		t.Error("OpenLogFile did not create expected file")
	}
}
