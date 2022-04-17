package celeritas

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
)

func (c *Celeritas) TakeScreenShot(pageURL, testName string, w, h float64) {
	page := rod.New().MustConnect().MustIgnoreCertErrors(true).MustPage(pageURL).MustWaitLoad()

	img, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Width:  w,
			Height: h,
			Scale:  1,
		},
		FromSurface: true,
	})
	if err != nil {
		fmt.Println(err)
	}

	fileName := time.Now().Format("2006-01-02-15-04-05.000000")
	err = utils.OutputFile(fmt.Sprintf("%s/screenshots/%s-%s.png", c.RootPath, testName, fileName), img)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Celeritas) FetchPage(pageURL string) *rod.Page {
	return rod.New().MustConnect().MustIgnoreCertErrors(true).MustPage(pageURL).MustWaitLoad()
}

func (c *Celeritas) SelectElementByID(page *rod.Page, id string) *rod.Element {
	// Commented out - I do not think this is how this function should be used
	// It returns an error
	//return page.MustElementByJS(fmt.Sprintf("document.getElementById('%s')", id))
	return page.MustElement("#" + id)
}

func (c *Celeritas) TakeScreenShotOfPage(page *rod.Page, testName string, w, h float64) {

	img, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Width:  w,
			Height: h,
			Scale:  1,
		},
		FromSurface: true,
	})
	if err != nil {
		fmt.Println(err)
	}

	fileName := time.Now().Format("2006-01-02-15-04-05.000000")
	err = utils.OutputFile(fmt.Sprintf("%s/screenshots/%s-%s.png", c.RootPath, testName, fileName), img)
	if err != nil {
		fmt.Println(err)
	}
}
