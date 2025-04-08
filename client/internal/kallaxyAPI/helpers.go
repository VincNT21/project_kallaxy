package kallaxyapi

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"

	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/VincNT21/kallaxy/client/models"
	"github.com/disintegration/imaging"
)

func (c *HelpersClient) GetMediaTypes(mediaRecords models.MediaWithRecords) map[string]bool {

	mediaTypes := make(map[string]bool)

	// Iterate over MediaRecords to check for different media types
	for mediaType := range mediaRecords.MediaRecords {
		mediaTypes[mediaType] = true
	}

	// Return data
	log.Println("--DEBUG-- GetMediaTypes() OK")
	return mediaTypes
}

func (c *HelpersClient) GetImage(imageUrl string) (*bytes.Buffer, error) {
	// Check if image exists in cache
	buf, exists := c.apiClient.Helpers.GetFromCache(imageUrl)
	if !exists {
		// If image doesn't exist, fetch it online
		buf, err := c.apiClient.Helpers.FetchImage(imageUrl)
		if err != nil {
			return nil, err
		}
		// Return the newly fetched image
		log.Println("--DEBUG-- GetImage() OK, image fetched")
		return buf, nil
	}
	// Return the cached image
	log.Println("--DEBUG-- GetImage() OK, image was on cache")
	return buf, nil
}

func (c *HelpersClient) GetFromCache(key string) (*bytes.Buffer, bool) {
	data, exists := c.apiClient.Cache.Get(key)                 // In stored cache
	dataTemp, existsTemp := c.apiClient.Cache.GetFromTemp(key) // In temp cache
	if !exists && !existsTemp {
		return nil, false
	} else if exists {
		return bytes.NewBuffer(data), true
	}
	return bytes.NewBuffer(dataTemp), true
}

func (c *HelpersClient) FetchImage(imageUrl string) (*bytes.Buffer, error) {

	// Make request
	r, err := http.Get(imageUrl)
	if err != nil {
		log.Printf("--ERROR-- with FetchImage(): %v\n", err)
		return nil, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with FetchImage(). Response status code: %v\n", r.StatusCode)
		return nil, fmt.Errorf("problem with FetchImage() request, status code: %v", r.StatusCode)
	}

	// Read the response body
	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("--ERROR-- with FetchImage(), couldn't io.ReadAll the response body: %v\n", err)
		return nil, err
	}

	// Resize image before caching and returning
	// (to max dimensions of 200x200)
	resizedImageData, err := resizeImage(imageData, 200, 200)
	if err != nil {
		log.Printf("--ERROR-- with FetchImage(), couldn't resize image: %v\n", err)
		return nil, err
	}

	// Store data in temp cache using the URL as the key
	c.apiClient.Cache.AddToTemp(imageUrl, resizedImageData)

	// Return data
	return bytes.NewBuffer(resizedImageData), nil
}

func resizeImage(originalImage []byte, maxWidth, maxHeight int) ([]byte, error) {

	// Decode the image bytes into an image.image and detect format
	img, format, err := image.Decode(bytes.NewReader(originalImage))
	if err != nil {
		log.Printf("--ERROR-- with resizeImage(), couldn't decode image and format: %v\n", err)
		return nil, err
	}

	// Resize the image while maintaining aspect ratio
	resizedImg := imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)

	// Encode back to bytes using the same format
	buf := new(bytes.Buffer)

	switch format {
	case "jpeg":
		err = jpeg.Encode(buf, resizedImg, nil)
	case "png":
		err = png.Encode(buf, resizedImg)
	case "gif":
		err = gif.Encode(buf, resizedImg, nil)
	default:
		// Default to JPEG if unknown
		err = jpeg.Encode(buf, resizedImg, nil)
	}

	if err != nil {
		log.Printf("--ERROR-- with resizeImage(), couldn't encode resized image according to format: %v\n", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
