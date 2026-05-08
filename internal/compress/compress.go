package compress

import (
	"os"
	"sort"
	"strconv"
	"strings"
)

func CompressFile(inputPath, outputPath string) error {
	// Read input file
	inputData, err := os.ReadFile(inputPath)

	if err != nil {
		return err
	}

	charFrequencyMap := GetCharFrequencyMap(inputData)
	encodingTree := GetEncodingTree(charFrequencyMap)

	// Byte stream to write compressed data to
	var compressedData []byte

	// Compress data using encoding tree
	for _, char := range inputData {
		encodedValue := encodingTree[string(char)]
		compressedData = append(compressedData, byte(encodedValue))
	}

	// First write encoding tree to output file for decompression
	encodingTreeString := GetEncodingTreeString(encodingTree)
	err = WriteCompressedDataToFile(outputPath, encodingTreeString, compressedData)
	if err != nil {
		return err
	}


	return nil
}

func GetCharFrequencyMap(inputData []byte) map[string]uint64 {
	var characterFrequencyMap map[string]uint64 = make(map[string]uint64)

	for _, character := range inputData {
		if characterFrequencyMap[string(character)] == 0 {
			characterFrequencyMap[string(character)] = 1
			continue
		}

		characterFrequencyMap[string(character)] += 1
	}

	// Sort map by frequency
	type kv struct {
		Key   string
		Value uint64
	}

	var sorted []kv

	for k, v := range characterFrequencyMap {
		sorted = append(sorted, kv{k, v})
	}

	// Sort the slice by frequency
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	return characterFrequencyMap
}

func GetEncodingTree(characterFrequencyMap map[string]uint64) map[string]uint64 {
	encodingTree := make(map[string]uint64)
	var idx uint = 0

	for char, _ := range characterFrequencyMap {
		encodingTree[char] = uint64(idx)
		idx++
	}

	return encodingTree
}

func GetEncodingTreeString(encodingTree map[string]uint64) string {
	var stringBuilder strings.Builder

	for character, idx := range encodingTree {
		idxStr := strconv.FormatUint(idx, 10)

		stringBuilder.WriteString(character);
		stringBuilder.WriteString(":");
		stringBuilder.WriteString(idxStr);
		stringBuilder.WriteString(",");
	}

	stringBuilder.WriteString("\n\n")

	return stringBuilder.String()
}

func WriteCompressedDataToFile(outputPath string, encodingTreeString string, compressedData []byte) error {
	outputFile, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.Write([]byte(encodingTreeString + "\n"))
	if err != nil {
		return err
	}

	_, err = outputFile.Write(compressedData)
	if err != nil {
		return err
	}

	return nil
}