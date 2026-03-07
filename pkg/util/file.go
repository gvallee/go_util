// Copyright (c) 2019-2026, Sylabs Inc. All rights reserved.
// Copyright (c) 2026, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

// Package util provides file, path, URL, and string utilities for reuse across projects.
package util

import (
	"fmt"
	"io"
	"os"
	"path"
)

// Tarball format constants for use with DetectTarballFormat and GetTarArgs.
const (
	// FormatBZ2 represents a bz2-compressed archive
	FormatBZ2 = "bz2"

	// FormatGZ represents a gzip-compressed archive
	FormatGZ = "gz"

	// FormatTGZ represents a gzip-compressed tarball (.tgz)
	FormatTGZ = "tgz"

	// FormatTAR represents an uncompressed TAR archive
	FormatTAR = "tar"

	// UnknownFormat is the value returned when the format cannot be detected
	UnknownFormat = ""
)

// PathExists reports whether the named path exists (returns false on any stat error, e.g. permission denied).
func PathExists(pathStr string) bool {
	_, err := os.Stat(pathStr)
	return err == nil
}

// FileExists reports whether pathStr exists and is a regular file.
// It returns false if the path does not exist, is a directory, or on any stat error.
func FileExists(pathStr string) bool {
	info, err := os.Stat(pathStr)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// IsDir reports whether pathStr exists and is a directory.
func IsDir(pathStr string) bool {
	info, err := os.Stat(pathStr)
	if err != nil {
		return false
	}
	return info.Mode().IsDir()
}

// GetTarArgs returns the arguments to use with tar based on the format of the tarball.
// It returns the empty string for unknown or unsupported format; callers should validate before passing to tar.
func GetTarArgs(format string) string {
	switch format {
	case FormatBZ2:
		return "-xjf"
	case FormatGZ:
		return "-xzf"
	case FormatTGZ:
		return "-xzf"
	case FormatTAR:
		return "-xf"
	}
	return ""
}

// DetectTarballFormat returns the archive format based on the file extension
// (e.g. .gz, .tar, .tgz, .bz2). Returns UnknownFormat for unknown or empty path.
func DetectTarballFormat(filepath string) string {
	if path.Ext(filepath) == ".bz2" {
		return FormatBZ2
	}

	if path.Ext(filepath) == ".gz" {
		return FormatGZ
	}

	if path.Ext(filepath) == ".tar" {
		return FormatTAR
	}

	if path.Ext(filepath) == ".tgz" {
		return FormatTGZ
	}

	return UnknownFormat
}

// CopyFile copies the regular file src to dst, preserving content and mode (permissions).
// If src and dst are the same path, it returns nil without copying.
// Returns error if src or dst is empty, if src is not a regular file, or if dst cannot be written (e.g. dst is an existing directory).
func CopyFile(src string, dst string) error {
	if src == dst {
		return nil
	}
	if src == "" {
		return fmt.Errorf("copy file: source path is empty")
	}
	if dst == "" {
		return fmt.Errorf("copy file: destination path is empty")
	}

	// Check that the source file is valid
	srcStat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("cannot access file %s: %w", src, err)
	}

	if !srcStat.Mode().IsRegular() {
		return fmt.Errorf("invalid source file %s: not a regular file", src)
	}

	// Actually do the copy
	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", src, err)
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", dst, err)
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return fmt.Errorf("unable to copy file from %s to %s: %w", src, dst, err)
	}
	if err := d.Sync(); err != nil {
		return fmt.Errorf("failed to sync %s: %w", dst, err)
	}

	// Check whether the copy succeeded by comparing the sizes of the two files
	dstStat, err := d.Stat()
	if err != nil {
		return fmt.Errorf("unable to get stat for %s: %w", d.Name(), err)
	}
	if srcStat.Size() != dstStat.Size() {
		return fmt.Errorf("file copy failed, size is %d instead of %d", srcStat.Size(), dstStat.Size())
	}
	// Preserve source mode so executables and permissions are copied correctly
	if err := os.Chmod(dst, srcStat.Mode()); err != nil {
		return fmt.Errorf("failed to set mode on %s: %w", dst, err)
	}
	return nil
}

// OpenResultsFile opens the result file for appending, or creates it if it does not exist.
// Returns error if path is empty. New files are created with mode 0644.
func OpenResultsFile(filepath string) (*os.File, error) {
	if filepath == "" {
		return nil, fmt.Errorf("results file path is empty")
	}
	mode := os.O_APPEND | os.O_RDWR
	if !FileExists(filepath) {
		mode = os.O_RDWR | os.O_CREATE
	}
	f, err := os.OpenFile(filepath, mode, 0644)
	if err != nil {
		return nil, fmt.Errorf("open results file %s: %w", filepath, err)
	}
	return f, nil
}

// OpenLogFile opens or creates the log file prefix-name.log for appending.
// Returns error if prefix or name is empty. New files are created with mode 0644.
func OpenLogFile(prefix string, name string) (*os.File, error) {
	if prefix == "" || name == "" {
		return nil, fmt.Errorf("open log file: prefix and name required")
	}
	filename := prefix + "-" + name + ".log"
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file %s: %w", filename, err)
	}
	return f, nil
}

// DirInit ensures that the target directory exists and is empty.
// If the path already exists, it is removed and recreated. Returns error if path is empty.
func DirInit(pathStr string) error {
	if pathStr == "" {
		return fmt.Errorf("directory path is empty")
	}
	if _, err := os.Stat(pathStr); err == nil {
		if err := os.RemoveAll(pathStr); err != nil {
			return fmt.Errorf("failed to remove existing path %s: %w", pathStr, err)
		}
	}
	if err := os.MkdirAll(pathStr, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", pathStr, err)
	}
	return nil
}
