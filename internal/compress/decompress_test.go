package compress

import (
	"os"
	"testing"
)

func TestDecompressSmallFile(t *testing.T) {
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

	// Decompress and verify content matches original
	decompressedFilePath := "../../mock_data/test.txt.unzip"
	err = DecompressFile(outputFilePath, decompressedFilePath)
	if err != nil {
		t.Fatalf("DecompressFile failed: %v", err)
	}

	decompressedData, err := os.ReadFile(decompressedFilePath)
	if err != nil {
		t.Fatalf("Failed to read decompressed file: %v", err)
	}

	originalData, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read original file: %v", err)
	}

	if string(decompressedData) != string(originalData) {
		t.Errorf("Decompressed data does not match original")
	}

	// Clean up
	os.Remove(outputFilePath)
	os.Remove(decompressedFilePath)
}

func TestDecompressBigFile(t *testing.T) {
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

	// Decompress and verify content matches original
	decompressedFilePath := "../../mock_data/big_test.txt.unzip"
	err = DecompressFile(outputFilePath, decompressedFilePath)
	if err != nil {
		t.Fatalf("DecompressFile failed: %v", err)
	}

	decompressedData, err := os.ReadFile(decompressedFilePath)
	if err != nil {
		t.Fatalf("Failed to read decompressed file: %v", err)
	}

	originalData, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read original file: %v", err)
	}

	if string(decompressedData) != string(originalData) {
		t.Errorf("Decompressed data does not match original")
	}

	// Clean up
	os.Remove(outputFilePath)
	os.Remove(decompressedFilePath)
}

