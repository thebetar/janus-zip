package compress

import (
	"os"
	"sort"
	"strconv"
	"strings"
)

// huffmanNode is a node in the Huffman tree used during encoding.
type huffmanNode struct {
	char        byte
	isLeaf      bool
	freq        uint64
	left, right *huffmanNode
}

func CompressFile(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return os.WriteFile(outputPath, []byte{}, 0644)
	}

	byteFrequencyMap := buildByteFrequencyMap(data)
	codeMap := buildEncodingMap(byteFrequencyMap)

	// Encode every byte into its Huffman bit string.
	var bitBuilder strings.Builder
	for _, b := range data {
		bitBuilder.WriteString(codeMap[b])
	}
	bits := bitBuilder.String()

	// Pad the bit string to a multiple of 8 and record how many bits were added.
	padding := (8 - len(bits)%8) % 8
	for i := 0; i < padding; i++ {
		bits += "0"
	}

	// Pack bits into bytes, MSB first. Prepend a padding-count byte.
	packed := make([]byte, 1+len(bits)/8)
	packed[0] = byte(padding)
	for i := 0; i < len(bits); i += 8 {
		var b byte
		for j := 0; j < 8; j++ {
			b <<= 1
			if bits[i+j] == '1' {
				b |= 1
			}
		}
		packed[1+i/8] = b
	}

	header := serializeCodeMap(codeMap)
	return writeCompressedFile(outputPath, header, packed)
}

func buildByteFrequencyMap(data []byte) map[byte]uint64 {
	byteFrequencyMap := make(map[byte]uint64)

	for _, dataByte := range data {
		byteFrequencyMap[dataByte]++
	}

	return byteFrequencyMap
}

func buildEncodingMap(byteFrequencyMap map[byte]uint64) map[byte]string {
	nodes := make([]*huffmanNode, 0, len(byteFrequencyMap))

	// Build initial list of leaf nodes, sorted for determinism.
	chars := make([]byte, 0, len(byteFrequencyMap))
	for c := range byteFrequencyMap {
		chars = append(chars, c)
	}

	sort.Slice(
		chars, 
		func(i, j int) bool { 
			return chars[i] < chars[j] 
		},
	)

	// Create all leaf nods
	for _, c := range chars {
		nodes = append(nodes, &huffmanNode{char: c, isLeaf: true, freq: byteFrequencyMap[c]})
	}

	// Edge case: only one distinct byte value.
	if len(nodes) == 1 {
		return map[byte]string{chars[0]: "0"}
	}

	sortNodes := func() {
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].freq < nodes[j].freq
		})
	}

	for len(nodes) > 1 {
		sortNodes()
		// Get first 2 nodes
		left, right := nodes[0], nodes[1]
		// Remove first two nodes from list
		nodes = nodes[2:]

		// Combine them into a new parent node and add back to list
		nodes = append(nodes, &huffmanNode{
			freq:  left.freq + right.freq,
			left:  left,
			right: right,
		})
	}

	codeMap := make(map[byte]string)
	generateCodes(nodes[0], "", codeMap)

	return codeMap
}

func generateCodes(node *huffmanNode, prefix string, codeMap map[byte]string) {
	if node == nil {
		return
	}
	if node.left == nil && node.right == nil {
		codeMap[node.char] = prefix
		return
	}
	generateCodes(node.left, prefix+"0", codeMap)
	generateCodes(node.right, prefix+"1", codeMap)
}

func serializeCodeMap(codeMap map[byte]string) string {
	chars := make([]byte, 0, len(codeMap))

	for c := range codeMap {
		chars = append(chars, c)
	}
	
	sort.Slice(chars, func(i, j int) bool { return chars[i] < chars[j] })

	var stringBuilder strings.Builder

	for _, c := range chars {
		stringBuilder.WriteString(strconv.Itoa(int(c)))
		stringBuilder.WriteByte(':')
		stringBuilder.WriteString(codeMap[c])
		stringBuilder.WriteByte(',')
	}
	
	stringBuilder.WriteString("\n\n")

	return stringBuilder.String()
}

func writeCompressedFile(outputPath, header string, data []byte) error {
	f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.WriteString(header); err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}