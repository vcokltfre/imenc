package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
)

func readFile(filename string) []byte {
	if filename == ":0" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("Error reading stdin: %s\n", err)
			os.Exit(1)
		}

		return data
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	return data
}

func writeFile(filename string, data []byte) {
	if filename == ":1" || filename == ":2" {
		var w io.Writer
		if filename == ":1" {
			w = os.Stdout
		} else {
			w = os.Stderr
		}

		_, err := w.Write(data)
		if err != nil {
			fmt.Printf("Error writing stdout: %s\n", err)
			os.Exit(1)
		}

		return
	}

	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file %s: %s\n", filename, err)
		os.Exit(1)
	}
}

func encode(infile, outfile string) {
	data := readFile(infile)
	size := int(math.Ceil(math.Sqrt(float64(len(data))) / 4))

	img := image.NewNRGBA(image.Rect(0, 0, size, size))

	copy(img.Pix, data)

	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, img)
	if err != nil {
		fmt.Printf("Error encoding image: %s\n", err)
		os.Exit(1)
	}

	writeFile(outfile, buf.Bytes())
}

func decode(infile, outfile string) {
	data := readFile(infile)

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		fmt.Printf("Error decoding image: %s\n", err)
		os.Exit(1)
	}

	if img.ColorModel() != color.NRGBAModel {
		fmt.Printf("Image is not RGBA\n")
		os.Exit(1)
	}

	buf := make([]byte, img.Bounds().Dx()*img.Bounds().Dy()*4)
	copy(buf, img.(*image.NRGBA).Pix)

	writeFile(outfile, buf)
}

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <encode|decode> <infile> <outfile>\n", os.Args[0])
		os.Exit(1)
	}

	cmd := os.Args[1]
	infile := os.Args[2]
	outfile := os.Args[3]

	switch cmd {
	case "encode":
		encode(infile, outfile)
	case "decode":
		decode(infile, outfile)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}
