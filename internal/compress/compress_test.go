package compress

import (
	"os"
	"testing"
)

func TestCompressSmallFile(t *testing.T) {
	testFilePath := "../../mock_data/test.txt"
	outputFilePath := "../../mock_data/test.txt.zip"

	err := CompressFile(testFilePath, outputFilePath)
	if err != nil {
		t.Fatalf("CompressFile failed: %v", err)
	}

	zipFile, err := os.Stat(outputFilePath)
	if os.IsNotExist(err) {
		t.Fatalf("Compressed file was not created")
	}

	// Check file size is smaller than original
	origFile, err := os.Stat(testFilePath)
	if os.IsNotExist(err) {
		t.Fatalf("Original file not found: %v", err)
	}

	if zipFile.Size() >= origFile.Size() {
		t.Errorf("Compressed file is not smaller than original: %d >= %d", zipFile.Size(), origFile.Size())
	}

	// Clean up
	os.Remove(outputFilePath)
}

func TestCompressBigFile(t *testing.T) {
	testFilePath := "../../mock_data/big_test.txt"
	outputFilePath := "../../mock_data/big_test.txt.zip"

	err := CompressFile(testFilePath, outputFilePath)
	if err != nil {
		t.Fatalf("CompressFile failed: %v", err)
	}

	zipFile, err := os.Stat(outputFilePath)
	if os.IsNotExist(err) {
		t.Fatalf("Compressed file was not created")
	}

	// Check file size is smaller than original
	origFile, err := os.Stat(testFilePath)
	if os.IsNotExist(err) {
		t.Fatalf("Original file not found: %v", err)
	}

	if zipFile.Size() >= origFile.Size() {
		t.Errorf("Compressed file is not smaller than original: %d >= %d", zipFile.Size(), origFile.Size())
	}

	// Clean up
	os.Remove(outputFilePath)
}

