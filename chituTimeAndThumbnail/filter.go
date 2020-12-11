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
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func readThumbNail(scanner *bufio.Scanner) ([]byte, error) {

	started := false
	ended := false
	buf := ""

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

	w.Write([]byte(fmt.Sprintf("M4010 X%d Y%d\n", img.Bounds().Max.X-img.Bounds().Min.X, img.Bounds().Max.Y-img.Bounds().Min.Y)))

	index_pixel := 0
	pixel_num := 0
	pixel_data := ""
	last_color := uint16(0)
	mask := uint16(32)
	unmask := ^mask
	color := uint16(0)
	first := true
	same_pixel := 1

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			pixel := img.At(x, y)
			r, g, b, a := pixel.RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			a = a >> 8

			if a < 100 {
				r = 255
				g = 255
				b = 255
			}

			rb := uint16(r >> 3)
			gb := uint16(g >> 2)
			bb := uint16(b >> 3)
			color = ((rb << 11) | (gb << 5) | bb) & unmask

			if first {
				last_color = color
				first = false
			} else {
				if last_color == color && same_pixel < 4095 {
					same_pixel++
				} else {
					if same_pixel >= 2 {

						pixel_data += fmt.Sprintf("%04x", uint16(last_color|mask))
						pixel_data += fmt.Sprintf("%04x", uint16(12288|same_pixel))
						pixel_num += same_pixel
						last_color = color
						same_pixel = 1

					} else {
						pixel_data += fmt.Sprintf("%04x", uint16(last_color))
						last_color = color
						pixel_num++
					}
				}
			}

			if len(pixel_data) >= 180 {
				w.Write([]byte(fmt.Sprintf("M4010 I%d T%d '%s'\n", index_pixel, pixel_num, pixel_data)))
				pixel_data = ""
				index_pixel += pixel_num
				pixel_num = 0
			}

		}
	}

	if pixel_num > 0 {

		if same_pixel >= 2 {
			pixel_data += fmt.Sprintf("%04x", uint16(last_color|mask))
			pixel_data += fmt.Sprintf("%04x", uint16(12288|same_pixel))
			pixel_num += same_pixel

		} else {
			pixel_data += fmt.Sprintf("%04x", uint16(last_color))
			pixel_num++
		}
		w.Write([]byte(fmt.Sprintf("M4010 I%d T%d '%s'\n", index_pixel, pixel_num, pixel_data)))

	}

}

func main() {

	if len(os.Args) != 2 {
		log.Fatalln("Please use with gcode filename")
	}

	input, err := os.Open(os.Args[1])
	if nil != err {
		log.Fatalln("Error opening input file"+os.Args[1]+":", err)
	}
	scanner := bufio.NewScanner(input)

	rawImage, err := readThumbNail(scanner)
	if nil != err {
		log.Fatalln("Error reading image from gcode:", err)
	}

	img, err := png.Decode(bytes.NewReader(rawImage))
	if nil != err {
		log.Fatalln("Error decoding thumbnail image from gcode:", err)
	}

	output, err := os.Create(os.Args[1] + ".tmp")
	if nil != err {
		log.Fatalln("Error creating temporary file"+os.Args[1]+".tmp:", err)
	}

	printImage(img, output)

	var M73 = regexp.MustCompile(`^M73 P([0-9]+) R([0-9]+)`)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasPrefix(s, "M73") {
			x := M73.FindStringSubmatch(s)
			if len(x) == 3 {
				p, _ := strconv.Atoi(x[1])
				r, _ := strconv.Atoi(x[2])
				output.WriteString(fmt.Sprintf("M2100 T%d\nM2101 T%d", r*60, p*60))
			}
		} else {
			output.WriteString(s)
		}
		output.Write([]byte("\n"))
	}

	input.Close()
	output.Close()

	os.Rename(os.Args[1]+".tmp", os.Args[1])
}
