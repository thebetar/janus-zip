package compress

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func DecompressFile(inputPath, outputPath string) error {
	// Read compressed file
	compressedData, err := os.ReadFile(inputPath)

	if err != nil {
		return err
	}

	// Split header and compressed data
	parts := strings.SplitN(string(compressedData), "\n\n", 2)

	if len(parts) < 2 {
		return fmt.Errorf("Invalid compressed file format")
	}

	headerData := []byte(parts[0])
	compressedData = []byte(parts[1])

	// First read encoding tree from file
	var decodingTree map[uint64]string
	decodingTree, err = ReadDecodingTreeFromFile(headerData)

	if err != nil {
		return err
	}

	// Decompress data using encoding tree
	var decompressedData []byte

	for _, byteValue := range compressedData {
		character, exists := decodingTree[uint64(byteValue)]

		if !exists {
			return fmt.Errorf("Invalid compressed data: byte value %d not found in decoding tree", byteValue)
		}

		decompressedData = append(decompressedData, []byte(character)...)
	}

	// Write decompressed data to output file
	err = os.WriteFile(outputPath, decompressedData, 0644)

	if err != nil {
		return err
	}

	return nil
}

func ReadDecodingTreeFromFile(headerData []byte) (map[uint64]string, error) {
	// Read encoding tree from file
	decodingTree := make(map[uint64]string)

	// Parse encoding tree string
	pairs := strings.Split(string(headerData), ",")

	for _, pair := range pairs {
		if pair == "" {
			continue
		}

		keyValue := strings.SplitN(pair, ":", 2)

		if len(keyValue) != 2 {
			return nil, fmt.Errorf("Invalid encoding tree format")
		}

		character := keyValue[0]
		idx, err := strconv.ParseUint(keyValue[1], 10, 64)

		if err != nil {
			return nil, fmt.Errorf("Invalid encoding tree format: %v", err)
		}

		decodingTree[idx] = character
	}

	return decodingTree, nil
}

