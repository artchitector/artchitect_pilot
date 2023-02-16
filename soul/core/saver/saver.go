package saver

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"mime/multipart"
	"net/http"
)

// Saver send binary image to saver-server, which lives in memory-server (near mother-database)
type Saver struct {
	saverURL string
}

func NewSaver(saverURL string) *Saver {
	return &Saver{saverURL}
}

func (s *Saver) SaveImage(cardID uint, imageData []byte) error {
	// Buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Initialize the file field
	fileWriter, err := multiPartWriter.CreateFormFile("file", fmt.Sprintf("card-%d.jpg", cardID))
	if err != nil {
		return errors.Wrapf(err, "[saver] failed card id=%d image saving", cardID)
	}

	// Copy the actual file content to the field field's writer
	r := bytes.NewReader(imageData)
	_, err = io.Copy(fileWriter, r)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed card id=%d image saving", cardID)
	}

	// Populate other fields
	fieldWriter, err := multiPartWriter.CreateFormField("card_id")
	if err != nil {
		return errors.Wrapf(err, "[saver] failed card id=%d image saving", cardID)
	}

	_, err = fieldWriter.Write([]byte(fmt.Sprintf("%d", cardID)))
	if err != nil {
		return errors.Wrapf(err, "[saver] failed card id=%d image saving", cardID)
	}

	// We completed adding the file and the fields, let's close the multipart writer
	// So it writes the ending boundary
	multiPartWriter.Close()

	// By now our original request body should have been populated, so let's just use it with our custom request
	pth := fmt.Sprintf("%s/upload", s.saverURL)
	req, err := http.NewRequest("POST", pth, &requestBody)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed card id=%d image saving", cardID)
	}
	// We need to set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	// Do the request
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed card id=%d image saving", cardID)
	}

	log.Info().Msgf("[saver] upload card %d to saver", cardID)

	return nil
}

func (s *Saver) SaveHundred(rank uint, hundred uint, imgFile []byte) error {
	// Buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Initialize the file field
	fileWriter, err := multiPartWriter.CreateFormFile("file", fmt.Sprintf("r-%d-h-%d.jpg", rank, hundred))
	if err != nil {
		return errors.Wrapf(err, "[saver] failed card r:%d h:%d", rank, hundred)
	}

	// Copy the actual file content to the field field's writer
	r := bytes.NewReader(imgFile)
	_, err = io.Copy(fileWriter, r)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed image saving r:%d h:%d", rank, hundred)
	}

	// Populate other fields
	fieldWriter, err := multiPartWriter.CreateFormField("rank")
	if err != nil {
		return errors.Wrapf(err, "[saver] failed create field rank, saving r:%d h:%d", rank, hundred)
	}

	_, err = fieldWriter.Write([]byte(fmt.Sprintf("%d", rank)))
	if err != nil {
		return errors.Wrapf(err, "[saver] failed write rank to field, saving r:%d h:%d", rank, hundred)
	}

	// Populate other fields
	fieldWriter, err = multiPartWriter.CreateFormField("hundred")
	if err != nil {
		return errors.Wrapf(err, "[saver] failed create field hundred, saving r:%d h:%d", rank, hundred)
	}

	_, err = fieldWriter.Write([]byte(fmt.Sprintf("%d", hundred)))
	if err != nil {
		return errors.Wrapf(err, "[saver] failed write field hundred", rank)
	}

	// We completed adding the file and the fields, let's close the multipart writer
	// So it writes the ending boundary
	multiPartWriter.Close()

	// By now our original request body should have been populated, so let's just use it with our custom request
	pth := fmt.Sprintf("%s/upload_hundred", s.saverURL)
	req, err := http.NewRequest("POST", pth, &requestBody)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed r:%d h:%d image saving", rank, hundred)
	}
	// We need to set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	// Do the request
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed make request r:%d, h:%d", rank, hundred)
	}

	log.Info().Msgf("[saver] upload image to saver r:%d h:%d", rank, hundred)

	return nil
}
