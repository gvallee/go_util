// Copyright (c) 2019-2026, Sylabs Inc. All rights reserved.
// Copyright (c) 2026, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

// Constants defining the format of the MPI package
const (
	// FormatBZ2 represents a bz2 file
	FormatBZ2 = "bz2"

	// FormatGZ represents a GZ file
	FormatGZ = "gz"

	// FormatTGZ represents a GZ file
	FormatTGZ = "tgz"

	// FormatTAR represents a simple TAR file
	FormatTAR = "tar"

	// UnknownFormat is the value returned when the format cannot be detected
	UnknownFormat = ""
)

// PathExists checks whether a given path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if err == nil {
		return true
	}

	return false
}

// FileExists is a utility function that checks whether a path exists and is a regular file.
// Returns false if the path does not exist, is a directory, or on any stat error.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// IsDir checks whether a path is pointing at a directory or not
func IsDir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.Mode().IsDir()
}

// GetTarArgs returns the arguments to use with tar based on the format of the tarball
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

// DetectTarballFormat detects the format of a tarball so we can know how
// to untar it.
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

// CopyFile is a utility functions that copies a file
func CopyFile(src string, dst string) error {
	if src == dst {
		return nil
	}

	log.Printf("* Copying %s to %s", src, dst)
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

	return nil
}

// OpenResultsFile opens the result files associated to the experiments
func OpenResultsFile(filepath string) *os.File {
	mode := os.O_APPEND | os.O_RDWR
	if !FileExists(filepath) {
		mode = os.O_RDWR | os.O_CREATE
	}
	f, err := os.OpenFile(filepath, mode, 0755)
	if err != nil {
		log.Printf("failed to open file %s: %s", filepath, err)
		return nil
	}

	return f
}

// OpenLogFile opens the log file for the execution of a command
func OpenLogFile(prefix string, name string) *os.File {
	if prefix == "" || name == "" {
		return nil
	}

	filename := prefix + "-" + name + ".log"
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("failed to create log file: %s", err)
		return nil
	}

	return logFile
}

// DirInit ensures that the target directory exists and is empty.
// If the path already exists, it is removed and recreated.
func DirInit(path string) error {
	if _, err := os.Stat(path); err == nil {
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove existing path %s: %w", path, err)
		}
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}
