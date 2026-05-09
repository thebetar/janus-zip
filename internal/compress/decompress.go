package compress

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// decodeNode is a trie node used to decode Huffman-encoded bits.
type decodeNode struct {
	char   byte
	isLeaf bool
	zero   *decodeNode // bit 0
	one    *decodeNode // bit 1
}

func DecompressFile(inputPath, outputPath string) error {
	raw, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if len(raw) == 0 {
		return os.WriteFile(outputPath, []byte{}, 0644)
	}

	// The header and body are separated by the first "\n\n"
	parts := strings.SplitN(string(raw), "\n\n", 2)
	if len(parts) < 2 {
		return fmt.Errorf("invalid compressed file format")
	}

	codeMap, err := parseHeader([]byte(parts[0]))
	if err != nil {
		return err
	}

	body := []byte(parts[1])
	if len(body) < 1 {
		return fmt.Errorf("invalid compressed file: missing data")
	}

	padding := int(body[0])
	packedData := body[1:]

	root := buildDecodeTree(codeMap)
	decompressed, err := decodeBits(packedData, padding, root)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, decompressed, 0644)
}

// parseHeader reads the "byte_val:binary_code,..." header into a code map.
func parseHeader(headerData []byte) (map[byte]string, error) {
	codeMap := make(map[byte]string)
	pairs := strings.Split(string(headerData), ",")

	for _, pair := range pairs {
		if pair == "" {
			continue
		}

		kv := strings.SplitN(pair, ":", 2)

		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid header format")
		}

		charVal, err := strconv.Atoi(kv[0])

		if err != nil || charVal < 0 || charVal > 255 {
			return nil, fmt.Errorf("invalid header: bad byte value %q", kv[0])
		}

		codeMap[byte(charVal)] = kv[1]
	}

	return codeMap, nil
}

// buildDecodeTree constructs a binary trie from the Huffman code map.
func buildDecodeTree(codeMap map[byte]string) *decodeNode {
	root := &decodeNode{}
	for char, code := range codeMap {
		node := root

		for _, bit := range code {
			if bit == '0' {
				if node.zero == nil {
					node.zero = &decodeNode{}
				}

				node = node.zero
			} else {
				if node.one == nil {
					node.one = &decodeNode{}
				}

				node = node.one
			}
		}

		node.isLeaf = true
		node.char = char
	}

	return root
}

// decodeBits walks the decode trie bit-by-bit (MSB first) and returns the
// original bytes. The last `padding` bits of the final byte are ignored.
func decodeBits(data []byte, padding int, root *decodeNode) ([]byte, error) {
	var result []byte
	node := root
	totalBits := len(data)*8 - padding

	for byteIdx, b := range data {
		for bitPos := 7; bitPos >= 0; bitPos-- {
			if byteIdx * 8 + (7 - bitPos) >= totalBits {
				break
			}

			if (b >> uint(bitPos)) & 1 == 0 {
				node = node.zero
			} else {
				node = node.one
			}

			if node == nil {
				return nil, fmt.Errorf("invalid compressed data: unexpected bit sequence")
			}

			if node.isLeaf {
				result = append(result, node.char)
				node = root
			}
		}
	}

	return result, nil
}

