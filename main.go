package main

import (
	"fmt"
	"image"
	"image/draw"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kbinani/screenshot"
	"github.com/nfnt/resize"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "VNC-HTTP is running")
	})
	http.HandleFunc("/out/video/raw", func(w http.ResponseWriter, r *http.Request) {
		width, height, err := getWidthHeight(r.URL.Query())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		img, err := capture(0, width, height)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		bounds := img.Bounds()
		rgba := image.NewRGBA(bounds)
		draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

		rawData := make([]byte, 0, bounds.Dx()*bounds.Dy()*3)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := rgba.At(x, y).RGBA()
				rawData = append(rawData, byte(r>>8), byte(g>>8), byte(b>>8))
			}
		}

		fmt.Fprint(w, string(rawData))
	})

	http.ListenAndServe(":8080", nil)
}

func getWidthHeight(val url.Values) (widthh, heightt uint, err error) {
	widthStr := val.Get("width")
	heightStr := val.Get("height")

	var width int
	if widthStr != "" {
		width, err = strconv.Atoi(widthStr)
		if err != nil {
			return 0, 0, err
		}
	} else {
		width = 192
	}

	var height int
	if heightStr != "" {
		height, err = strconv.Atoi(heightStr)
		if err != nil {
			return
		}
	} else {
		height = 108
	}

	return uint(width), uint(height), nil
}

func capture(display int, width, height uint) (image.Image, error) {
	img, err := screenshot.CaptureDisplay(display)
	if err != nil {
		return nil, err
	}

	return resize.Resize(width, height, img, resize.NearestNeighbor), nil
}
