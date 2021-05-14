package components

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// PageImages holds the image urls and the page url
type PageImages struct {
	Src     []string
	Page    string
	pageURL *url.URL
}

func (pageImages *PageImages) parsePageURL() (err error) {
	if pageImages.pageURL != nil {
		return
	}

	pageImages.pageURL, err = url.Parse(pageImages.Page)
	return
}

// ToImageSlice for each image src, makes sure it's property formatted and return a slice of
// properly formatted images
func (pageImages *PageImages) ToImageSlice() (ImageSources, error) {
	imgSrcs := make(ImageSources, len(pageImages.Src))

	if err := pageImages.parsePageURL(); err != nil {
		return nil, err
	}

	// Keep track of invalid urls
	skipped := 0

	for i, url := range pageImages.Src {
		if len(url) < 2 {
			// Invalid url
			skipped++
			continue
		} else if url[:2] == "//" {
			// Scheme relative url
			imgSrcs[i] = ImageSrc(pageImages.pageURL.Scheme + ":" + url)
		} else if url[0] == '/' {
			// Page relative url
			baseURL := pageImages.pageURL
			imgSrcs[i] = ImageSrc(baseURL.Scheme + "://" + baseURL.Host + url)
		} else if url[:4] == "http" {
			// Return as is
			imgSrcs[i] = ImageSrc(url)
		} else {
			// This is probably a current location relative url
			imgSrcs[i] = ImageSrc(pageImages.Page + url)
		}
	}

	return imgSrcs[:len(imgSrcs)-skipped], nil
}

// NewPageImages ...
func NewPageImages(src []string, page string) PageImages {
	return PageImages{
		Src:  src,
		Page: page,
	}
}

// Page contains the contents of a page
type Page struct {
	URL  string
	HTML *string
}

// NewPage ...
func NewPage(url, html string) Page {
	return Page{
		URL:  url,
		HTML: &html,
	}
}

// ImageSrc is just an image url string with attached special functions
type ImageSrc string

func (src ImageSrc) String() string {
	return string(src)
}

// Get will retrieve the image pointed to by the URL
func (src ImageSrc) Get() (io.Reader, error) {
	resp, err := http.Get(src.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return resp.Body, nil
	}

	return nil, fmt.Errorf("Received an error code: %d on '%s'", resp.StatusCode, src)
}

// GetData is an enhancement around `Get` that wraps the result in ImageData with
// extra metadata included
func (src ImageSrc) GetData() (ImageData, error) {
	resp, err := http.Get(src.String())
	if err != nil {
		return ImageData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		return ImageData{
			File:    resp.Body,
			Src:     src,
			headers: resp.Header,
		}, nil
	}

	return ImageData{}, fmt.Errorf("Received an error code: %d on '%s'", resp.StatusCode, src)

}

// ImageSources is a slice of ImageSrc entries
type ImageSources []ImageSrc

// ImageData contains the image url and the file data
type ImageData struct {
	Src     ImageSrc
	File    io.Reader
	headers http.Header
}
