package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

const (
	whiteIndex = iota
	redIndex
	greenIndex
	blueIndex
	blackIndex
)

var palette = []color.Color{
	color.White,
	redIndex:   color.RGBA{0xFF, 0x00, 0x00, 0xFF},
	greenIndex: color.RGBA{0x00, 0xFF, 0x00, 0xFF},
	blueIndex:  color.RGBA{0x00, 0x00, 0xFF, 0xFF},
	blackIndex: color.Black,
}

var (
	mu    sync.Mutex
	count int
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/lissajous", liassajousHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func lissajousRenderer(out io.Writer, cycles float64) {
	if cycles == 0 {
		cycles = 5
	}
	const (
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 640   // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				uint8(rand.Uint64()%blackIndex+1))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

func liassajousHandler(w http.ResponseWriter, r *http.Request) {
	cycles, err := strconv.Atoi(r.URL.Query().Get("cycles"))
	if err != nil {
		cycles = 0
	}
	lissajousRenderer(w, float64(cycles))
}

// handler echoes the Path component of the requested URL.
func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// counter echoes the number of calls so far.
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}
