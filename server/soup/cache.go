package soup

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"time"

	"github.com/Mopsgamer/space-soup/server/environment"
)

var ErrCacheExpired = errors.New("cache expired")

type SoupCache struct {
	ExpiresAt  time.Time
	IsImmortal bool

	TestList []MovementTest

	ExcelBytes     []byte
	PlotImageBytes []byte
}

func NewCache(tests []MovementTest, rang, description string) (cache *SoupCache, err error) {
	movementTestList := tests
	if rang != "" {
		movementTestList, err = Range(tests, rang)
		if err != nil {
			return
		}
	}

	imageBytes, err := NewImageBytes(movementTestList, description)
	if err != nil {
		return
	}

	excelBytes, err := NewFileExcelBytes(movementTestList)
	if err != nil {
		return
	}

	cache = &SoupCache{
		TestList:       movementTestList,
		PlotImageBytes: imageBytes,
		ExcelBytes:     excelBytes,
	}
	return
}

func (cache SoupCache) IsExpired() bool {
	return !cache.IsImmortal && cache.ExpiresAt.Before(time.Now())
}

// Increases ExpiresAt time.
func (cache *SoupCache) Live() error {
	if cache == nil {
		return ErrCacheExpired
	}
	cache.ExpiresAt = time.Now().Add(environment.AlgCacheDuration)
	return nil
}

type FileHashCacheMap map[string]*SoupCache

func (m FileHashCacheMap) Free() {
	for hash, cache := range m {
		if cache.ExpiresAt.After(time.Now()) {
			continue
		}

		delete(m, hash)
	}
}

func (m FileHashCacheMap) Live(key string) error {
	err := m[key].Live()
	if err != nil {
		return err
	}
	m.Free()
	return nil
}

func (m FileHashCacheMap) Add(key string, cache *SoupCache) error {
	err := cache.Live()
	if err != nil {
		return err
	}
	m.Free()
	m[key] = cache
	return nil
}

func NewImageBytes(tests []MovementTest, description string) (imageBytes []byte, err error) {
	var imageWriterTo io.WriterTo
	imageWriterTo, err = VisualizeAlphaDelta(tests, description)
	if err != nil {
		return
	}

	imageBuffer := bytes.NewBuffer([]byte{})
	_, err = imageWriterTo.WriteTo(imageBuffer)
	if err != nil {
		return
	}
	imageBytes = imageBuffer.Bytes()
	return
}

func NewFormFileBytes(formFile *multipart.FileHeader) (fileBytes []byte, err error) {
	var file multipart.File
	file, err = formFile.Open()
	if err != nil {
		return
	}
	if _, err = file.Read(fileBytes); err != nil {
		file.Close()
		return
	}
	file.Close()
	return
}
