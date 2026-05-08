package main

import (
	"fmt"
	"janus-zip/internal/compress"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: janus-zip <file>")
		return
	}

	var decompress_flag bool = false
	var inputFile string

	if len(os.Args) == 3 && (os.Args[1] == "-d" || os.Args[1] == "--decompress") {
		decompress_flag = true
		inputFile = os.Args[2]
	} else {
		inputFile = os.Args[1]
	}


	_, err := os.Stat(inputFile)

	if os.IsNotExist(err) {
		fmt.Printf("File %s does not exist\n", inputFile)
		return
	}

	var outputFile string

	if decompress_flag {
		outputFile = inputFile + ".unzip"
	} else {
		outputFile = inputFile + ".zip"
	}

	if decompress_flag {
		err = compress.DecompressFile(inputFile, outputFile)
	} else {
		err = compress.CompressFile(inputFile, outputFile)
	}

	if err != nil {
		if decompress_flag {
			fmt.Printf("Error decompressing file: %v\n", err)
		} else {
			fmt.Printf("Error compressing file: %v\n", err)
		}
		return
	}

	if decompress_flag {
		fmt.Printf("File decompressed to %s\n", outputFile)
	} else {
		fmt.Printf("File compressed to %s\n", outputFile)
	}
}