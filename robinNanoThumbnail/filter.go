package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/nfnt/resize"
	"github.com/prometheus/common/log"
)

func readThumbNail(filename string) ([]byte, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	started := false
	ended := false
	buf := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasPrefix(s, "; thumbnail begin") {
			started = true
			continue
		}

		if !started {
			continue
		}

		if started {
			if strings.HasPrefix(s, "; thumbnail end") {
				ended = true
				break
			}
		}
		//		fmt.Println(s)
		if strings.HasPrefix(s, "; ") {
			buf += s[2:]
		}
	}

	if !started {
		return nil, errors.New("No thumbnail in file")
	}

	if !ended {
		return nil, errors.New("No thumbnail end in file")
	}

	data, err := base64.StdEncoding.DecodeString(buf)
	if nil != err {
		return nil, err
	}

	return data, nil
}

func printImage(img image.Image, w io.Writer) {

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {

		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			color := img.At(x, y)
			r, g, b, a := color.RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			a = a >> 8

			if a < 100 {
				w.Write([]byte("0000"))
				continue
			}

			rb := uint16(r >> 3)
			gb := uint16(g >> 2)
			bb := uint16(b >> 3)
			rgb := 0xffff & ((rb << 11) | (gb << 5) | bb)

			w.Write([]byte(fmt.Sprintf("%04x", uint16(rgb))))
		}

		w.Write([]byte("\rM10086 ;"))
	}
}

func main() {

	if len(os.Args) != 2 {
		log.Fatalln("Please use with gcode filename")
	}

	raw, err := readThumbNail(os.Args[1])
	if nil != err {
		log.Fatalln("Error reading thumbnail from gcode:", err)
	}

	img, err := png.Decode(bytes.NewReader(raw))
	if nil != err {
		log.Fatalln("Error decoding thumbnail image from gcode:", err)
	}

	small := resize.Resize(100, 100, img, resize.Lanczos3)
	if img.Bounds().Dx() != 200 || img.Bounds().Dy() != 200 {
		img = resize.Resize(200, 200, img, resize.Lanczos3)
	}

	file, err := os.Create(os.Args[1] + ".tmp")
	if nil != err {
		log.Fatalln("Error creating temporary file"+os.Args[1]+".tmp:", err)
	}

	file.Write([]byte(";simage:"))
	printImage(small, file)
	file.Write([]byte("\r;;gimage:"))
	printImage(img, file)
	file.Write([]byte("\n"))

	input, err := os.Open(os.Args[1])
	if nil != err {
		log.Fatalln("Error re-opening input file"+os.Args[1]+":", err)
	}

	io.Copy(file, input)

	input.Close()
	file.Close()

	os.Rename(os.Args[1]+".tmp", os.Args[1])
}
