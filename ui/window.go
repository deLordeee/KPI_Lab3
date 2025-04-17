package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	
	window      screen.Window
	textureChan chan screen.Texture
	quitChan    chan struct{}


	currentSize size.Event
	drawRect    image.Rectangle
	clickPos    image.Point
}

func (pw *Visualizer) Main() {
	pw.textureChan = make(chan screen.Texture)
	pw.quitChan = make(chan struct{})
	pw.drawRect.Max.X = 200
	pw.drawRect.Max.Y = 200
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.textureChan <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	win, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  pw.Title,
		Width:  800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize window:", err)
	}
	defer func() {
		win.Release()
		close(pw.quitChan)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.window = win

	events := make(chan any)
	go func() {
		for {
			e := win.NextEvent()
			if pw.Debug {
				log.Printf("Event received: %v", e)
			}
			if pw.detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var currentTex screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, currentTex)

		case currentTex = <-pw.textureChan:
			win.Send(paint.Event{})
		}
	}
}

func (pw *Visualizer) detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, tex screen.Texture) {
	switch e := e.(type) {
	case size.Event:
		pw.currentSize = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if e.Button == mouse.ButtonRight && e.Direction == mouse.DirPress {
			pw.clickPos = image.Point{X: int(e.X), Y: int(e.Y)}
			pw.window.Send(paint.Event{})
		}

	case paint.Event:
		if tex == nil {
			pw.drawDefaultUI()
		} else {
			pw.window.Scale(pw.currentSize.Bounds(), tex, tex.Bounds(), draw.Src, nil)
		}
		pw.window.Publish()
	}
}


func (pw *Visualizer) drawDefaultUI() {
	pw.window.Fill(pw.currentSize.Bounds(), color.RGBA{G: 255, A: 255}, draw.Src)

	x, y := 400, 400
	if pw.clickPos != (image.Point{}) {
		x, y = pw.clickPos.X, pw.clickPos.Y
	}

	barWidth := 125
	verticalBar := image.Rect(x-barWidth/2, y, x+barWidth/2, y+200)

	horizontalWidth := 400
	horizontalBar := image.Rect(x-horizontalWidth/2, y-150, x+horizontalWidth/2, y)

	pw.window.Fill(verticalBar, color.RGBA{R: 255, G: 255, A: 255}, draw.Src)
	pw.window.Fill(horizontalBar, color.RGBA{R: 255, G: 255, A: 255}, draw.Src)
}
