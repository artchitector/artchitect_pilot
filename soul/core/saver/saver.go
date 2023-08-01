package saver

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"mime/multipart"
	"net/http"
)

// Saver send binary image to saver-server, which lives in memory-server (near mother-database)
type Saver struct {
	memorySaverURL  string
	storageSaverURL string
}

func NewSaver(saverURL string, storageSaverURL string) *Saver {
	return &Saver{saverURL, storageSaverURL}
}

func (s *Saver) SaveArt(ctx context.Context, artID uint, imageData []byte) error {
	// Buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Initialize the file field
	fileWriter, err := multiPartWriter.CreateFormFile("file", fmt.Sprintf("art-%d.jpg", artID))
	if err != nil {
		return errors.Wrapf(err, "[saver] failed art id=%d image saving", artID)
	}

	// Copy the actual file content to the field field's writer
	r := bytes.NewReader(imageData)
	_, err = io.Copy(fileWriter, r)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed art id=%d image saving", artID)
	}

	// Populate other fields
	fieldWriter, err := multiPartWriter.CreateFormField("art_id")
	if err != nil {
		return errors.Wrapf(err, "[saver] failed art id=%d image saving", artID)
	}

	_, err = fieldWriter.Write([]byte(fmt.Sprintf("%d", artID)))
	if err != nil {
		return errors.Wrapf(err, "[saver] failed art id=%d image saving", artID)
	}

	// We completed adding the file and the fields, let's close the multipart writer
	// So it writes the ending boundary
	multiPartWriter.Close()

	// By now our original request body should have been populated, so let's just use it with our custom request
	pth := fmt.Sprintf("%s/upload_art", s.memorySaverURL)
	req, err := http.NewRequest("POST", pth, &requestBody)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed art id=%d image saving", artID)
	}
	// We need to set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	// Do the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed art id=%d image saving", artID)
	}

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("[saver] failed to upload art. URL: %s. Status: %d", pth, res.StatusCode)
	}

	log.Info().Msgf("[saver] uploaded art %d to saver. URL: %s, Status: %d", artID, pth, res.StatusCode)

	return nil
}

func (s *Saver) SaveUnity(ctx context.Context, filename string, imgFile []byte) error {
	// Buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Initialize the file field
	fileWriter, err := multiPartWriter.CreateFormFile("file", fmt.Sprintf("%s.jpg", filename))
	if err != nil {
		return errors.Wrapf(err, "[saver] failed createFromFile %s", filename)
	}

	// Copy the actual file content to the field field's writer
	r := bytes.NewReader(imgFile)
	_, err = io.Copy(fileWriter, r)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed image saving %s", filename)
	}

	// Populate other fields
	fieldWriter, err := multiPartWriter.CreateFormField("filename")
	if err != nil {
		return errors.Wrapf(err, "[saver] failed create field filename, saving %s", filename)
	}

	_, err = fieldWriter.Write([]byte(filename))
	if err != nil {
		return errors.Wrapf(err, "[saver] failed write filename to field, saving %s", filename)
	}

	// We completed adding the file and the fields, let's close the multipart writer
	// So it writes the ending boundary
	multiPartWriter.Close()

	// By now our original request body should have been populated, so let's just use it with our custom request
	pth := fmt.Sprintf("%s/upload_unity", s.memorySaverURL)
	req, err := http.NewRequest("POST", pth, &requestBody)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed %s image saving", filename)
	}
	// We need to set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	// Do the request
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed make request %s", filename)
	}

	log.Info().Msgf("[saver] upload unity to saver %s", filename)

	return nil
}

func (s *Saver) SaveFullsize(ctx context.Context, artID uint, imageData []byte) error {
	// Buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Initialize the file field
	fileWriter, err := multiPartWriter.CreateFormFile("file", fmt.Sprintf("art-%d.jpg", artID))
	if err != nil {
		return errors.Wrapf(err, "[saver] failed art id=%d image saving", artID)
	}

	// Copy the actual file content to the field field's writer
	r := bytes.NewReader(imageData)
	_, err = io.Copy(fileWriter, r)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed art id=%d image saving", artID)
	}

	// Populate other fields
	fieldWriter, err := multiPartWriter.CreateFormField("art_id")
	if err != nil {
		return errors.Wrapf(err, "[saver] failed art id=%d image saving", artID)
	}

	_, err = fieldWriter.Write([]byte(fmt.Sprintf("%d", artID)))
	if err != nil {
		return errors.Wrapf(err, "[saver] failed art id=%d image saving", artID)
	}

	// We completed adding the file and the fields, let's close the multipart writer
	// So it writes the ending boundary
	multiPartWriter.Close()

	// By now our original request body should have been populated, so let's just use it with our custom request
	pth := fmt.Sprintf("%s/upload_fullsize", s.storageSaverURL)
	req, err := http.NewRequest("POST", pth, &requestBody)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed art id=%d image saving", artID)
	}
	// We need to set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	// Do the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed art id=%d image saving", artID)
	}

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("[saver] failed to upload art. URL: %s. Status: %d", pth, res.StatusCode)
	}

	log.Info().Msgf("[saver] uploaded art %d to saver. URL: %s, Status: %d", artID, pth, res.StatusCode)

	return nil
}
