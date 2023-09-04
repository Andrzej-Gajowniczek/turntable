package main

import (
	"C"
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

//go:embed "img/turntable-png-17.png"
var armT []byte

//go:embed "img/vinyl.png"
var imageV []byte

//go:embed "73.mp3"
var streaM []byte

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "GoWithAndy",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)

	img, _, err := image.Decode(bytes.NewReader(imageV))
	if err != nil {
		log.Fatal("blad", err)
	}
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())

	img2, _, err := image.Decode(bytes.NewReader(armT))
	if err != nil {
		log.Fatal("blad", err)
	}
	pic2 := pixel.PictureDataFromImage(img2)
	arm := pixel.NewSprite(pic2, pic2.Bounds())

	armPos := pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(180, 40))) //.Add(pixel.V(100, 230)))
	//armPos = armPos.Rotated(win.Bounds().Center(), -6.28/21)

	armPos = armPos.ScaledXY(win.Bounds().Center(), pixel.V(1.5, 1.5))

	angle := 0.0

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		angle -= 3.3 * dt

		win.Clear(color.RGBA{30, 45, 45, 0})

		mat := pixel.IM
		mat = mat.Rotated(pixel.ZV, angle)
		mat = mat.Moved(win.Bounds().Center().Add(pixel.V(-90, 0)))
		sprite.Draw(win, mat)
		arm.Draw(win, armPos)
		win.Update()
	}
}

type ByteReadCloser struct {
	*bytes.Reader
}

func (b *ByteReadCloser) Close() error {
	return nil
}

func main() {
	go func() {
		for { /*
				f, err := os.Open("moon.mp3")
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()
			*/

			mp3Reader := &ByteReadCloser{
				bytes.NewReader(streaM),
			}
			streamer, format, err := mp3.Decode(mp3Reader)
			if err != nil {
				log.Fatal(err)
			}
			defer streamer.Close()

			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
			done := make(chan struct{})

			speaker.Play(beep.Seq(streamer, beep.Callback(func() {
				close(done)
			})))

			<-done

		}
	}()

	pixelgl.Run(run)
}
