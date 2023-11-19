package components

import (
	"log"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/elias-gill/walldo-in-go/globals"
	"github.com/elias-gill/walldo-in-go/utils"
	"github.com/elias-gill/walldo-in-go/wallpaper"
)

type wallpapersGrid struct {
	content *fyne.Container
}

func NewImageGrid() wallpapersGrid {
	res := wallpapersGrid{content: container.NewWithoutLayout()}
	res.defineCardSize()
	return res
}

func (c *wallpapersGrid) GetGridContent() *fyne.Container {
	return c.content
}

func (c *wallpapersGrid) RefreshImgGrid() {
	c.defineCardSize()
	c.FillGrid()
}

// Generates and return a new layout acording to the user configurations
func (c *wallpapersGrid) defineCardSize() {
	// default card size
	size := fyne.NewSize(150, 130)
	// other grid sizes
	switch globals.GridSize {
	case "small":
		size = fyne.NewSize(110, 100)
	case "large":
		size = fyne.NewSize(195, 175)
	}
	c.content.Layout = layout.NewGridWrapLayout(size)
}

// fills the container with the correspondent content
func (c wallpapersGrid) FillGrid() {
	c.content.RemoveAll()
	imagesList := utils.ListImagesRecursivelly() // search original images

	// save all images into a go channel to manage concurrently load/generate thumbnails
	channel := make(chan string, len(imagesList))
	for _, v := range imagesList {
		channel <- v
	}

	// create more "threads" to increase performance
	for i := 0; i < runtime.NumCPU()-2; i++ {
		go c.addNewCard(channel)
	}
	print("\n Usando ", runtime.NumCPU()-2, " Hilos")
}

// Recibes the channel with the list of images and creates a new card from the every entry
// WARN: needs a rework to possibly improve performace
func (c wallpapersGrid) addNewCard(chanel chan string) {
	for image := range chanel {
		button := widget.NewButton("", func() {
			// the button has the index of the original image
			err := wallpaper.SetFromFile(image)
			if err != nil {
				log.Println(err.Error())
			}
		})

		// resize the image and get the thumbnail name
		thumbail := utils.ResizeImage(image)
		aux := canvas.NewImageFromFile(thumbail)
		aux.ScaleMode = canvas.ImageScaleFastest
		aux.FillMode = canvas.ImageFillContain

		// With the max layout we can overlap the button and the thumbnail
		c.content.Add(container.NewMax(aux, button))
		c.content.Refresh()
	}
}