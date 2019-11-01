// Copyright 2019 Roosevelt Purification. All rights reserved.

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

func main() {

	runtime.GOMAXPROCS(5)

	fmt.Println("How many variations do you want for each name? ")

	in := bufio.NewReader(os.Stdin)
	inputLine, inputError := in.ReadString('\n')
	if inputError != nil {
		log.Fatal(inputError)
	}

	numberOfVariationsDesired, conversionError := strconv.Atoi(strings.TrimSuffix(inputLine, "\n"))
	if conversionError != nil {
		log.Fatal(conversionError)
	}

	// open the CSV file
	csvFile, _ := os.Open("names.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	var waitGrp sync.WaitGroup

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		// generate logo variations for each name
		for i := 1; i < numberOfVariationsDesired+1; i++ {
			// Generate the logos
			waitGrp.Add(1)

			go func(variation_number_to_use int) {
				defer waitGrp.Done()
				generateLogo(line[0], variation_number_to_use)
			}(i)
		}
	}

	waitGrp.Wait()
}

func generateLogo(companyName string, variation_number int) {

	fmt.Println(fmt.Sprintf("Generating %s Variation %d", companyName, variation_number))

	// regex for a clean file name
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal("Bad Regex: %s", err)
	}

	logoFileSaveName := fmt.Sprintf("%s_Variation_%d_%s", reg.ReplaceAllString(companyName, ""), variation_number, time.Now().Format("20060102150405"))

	// setup the background
	width := 1920
	height := 1080

	upperLeft := image.Point{0, 0}
	lowerRight := image.Point{width, height}

	myLogo := image.NewRGBA(image.Rectangle{upperLeft, lowerRight})

	// background color
	backgroundColor := image.NewUniform(getRandomColor())

	// draw the background
	draw.Draw(myLogo, myLogo.Bounds(), backgroundColor, image.ZP, draw.Src)

	// write the company name
	fontSize := float64(200)
	fontSpacing := float64(.7)

	myFontContext := getRandomFontAndContext()
	myFontContext.SetClip(myLogo.Bounds())
	myFontContext.SetDst(myLogo)
	myFontContext.SetFontSize(fontSize) //font size in points

	// center the logo
	companyNameLength := len(companyName)
	centerPointX := (width / companyNameLength) + (width / companyNameLength) - 70

	pt := freetype.Pt(centerPointX, height/2+70)

	for _, str := range companyName {
		_, err := myFontContext.DrawString(string(str), pt)
		if err != nil {
			fmt.Println(err)
			return
		}

		pt.X += myFontContext.PointToFixed(fontSize * fontSpacing)
	}

	// save the file
	f, _ := os.Create(fmt.Sprintf("output/%s.png", logoFileSaveName))
	png.Encode(f, myLogo)
}

func getRandomFontAndContext() *freetype.Context {
	// set the font

	fontFile := fmt.Sprintf("fonts/%s.ttf", getRandomFont())
	fontDPI := float64(72)
	fontContext := new(freetype.Context)
	utf8Font := new(truetype.Font)
	fontColor := color.RGBA{255, 255, 255, 255}

	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		log.Fatal(err)
	}

	utf8Font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	fontForeGroundColor := image.NewUniform(fontColor)

	fontContext = freetype.NewContext()
	fontContext.SetDPI(fontDPI) //screen resolution in Dots Per Inch
	fontContext.SetFont(utf8Font)
	fontContext.SetSrc(fontForeGroundColor)

	return fontContext
}

func getRandomColor() color.RGBA {

	rand.Seed(time.Now().UnixNano())

	max := 255
	min := 0

	red := uint8(rand.Intn(max-min) + min)
	green := uint8(rand.Intn(max-min) + min)
	blue := uint8(rand.Intn(max-min) + min)

	return color.RGBA{red, green, blue, 0xff}
}

func getRandomFont() string {

	fonts := []string{
		"ZCOOLXiaoWei-Regular",
		"DancingScript-Bold",
		"DancingScript-Regular",
		"Arvo-Bold",
	}

	randomIndex := rand.Int() % len(fonts)

	return fonts[randomIndex]
}
