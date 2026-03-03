package downloader

import (
	"fmt"
	"io"
	"os"
)

// renameDir renames the source directory to the destination directory.
// This is used to move the downloaded content from a temp dir to the final destination.
// On Windows, os.Rename doesn't work if the destination already exists,
// so we need to copy the content and remove the source manually.
func renameDir(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("couldn't copy to dest from source: %v", err)
	}

	// for Windows, close before trying to remove:
	// https://stackoverflow.com/a/64943554/246801
	inputFile.Close()

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't remove source file: %v", err)
	}
	return nil
}

// removeBackup remove backup once we no longer need it.
func removeBackup(backup string) {
	if backup != "" {
		_ = os.RemoveAll(backup)
	}
}

// restoreBackup put backup back as dst on 304 or download failure.
func restoreBackup(backup, dst string) {
	if backup != "" {
		_ = os.RemoveAll(dst) // remove any partial dst left by a failed download
		_ = renameDir(backup, dst)
	}
}
