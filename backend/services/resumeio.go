package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

const metadataResourceURL = "https://ssr.resume.tools/meta/%s?cache=%s"
const imageResourceURL = "https://ssr.resume.tools/to-image/%s-%d.%s?cache=%s&size=%d"

func GeneratePDF(renderingToken string) ([]byte, error) {
	cacheDate := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	metadata, err := getMetadata(renderingToken, cacheDate)
	if err != nil {
		return nil, err
	}

	if len(metadata.Pages) == 0 {
		return nil, fmt.Errorf("no pages found in metadata")
	}

	imagesBytes := make([][]byte, len(metadata.Pages))

	for i := range metadata.Pages {
		imgBytes, err := downloadImageFromURL(renderingToken, i+1, cacheDate)
		if err != nil {
			return nil, err
		}
		imagesBytes[i] = imgBytes
	}

	pageWidth, pageHeight := metadata.Pages[0].Viewport.Width, metadata.Pages[0].Viewport.Height

	pdfConfig := gofpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: pageWidth, Ht: pageHeight},
	}
	pdf := gofpdf.NewCustom(&pdfConfig)
	defer pdf.Close()

	for i, imageBytes := range imagesBytes {
		pdf.AddPage()

		img := bytes.NewReader(imageBytes)
		imageName := "image_" + strconv.Itoa(i)
		pdf.RegisterImageReader(imageName, "jpeg", img)
		pdf.Image(imageName, 0, 0, pageWidth, pageHeight, false, "JPEG", 0, "")
	}

	output := new(bytes.Buffer)
	if err := pdf.Output(output); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

type resumeMetadata struct {
	Pages []struct {
		Viewport struct {
			Width  float64 `json:"width,omitempty"`
			Height float64 `json:"height,omitempty"`
		} `json:"viewport,omitempty"`
		Links []struct {
			Left   float64 `json:"left,omitempty"`
			Top    float64 `json:"top,omitempty"`
			Width  float64 `json:"width,omitempty"`
			Height float64 `json:"height,omitempty"`
			URL    string  `json:"url,omitempty"`
		}
	} `json:"pages"`
}

func getMetadata(renderingToken, cacheDate string) (*resumeMetadata, error) {
	resp, err := http.Get(fmt.Sprintf(metadataResourceURL, renderingToken, cacheDate))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var metadata resumeMetadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func downloadImageFromURL(renderingToken string, pageID int, cacheDate string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf(imageResourceURL, renderingToken, pageID, "jpeg", cacheDate, 3000))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}
