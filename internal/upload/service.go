package upload

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/okanay/file-upload-go/types"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type Service struct {
	uploadRepo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{uploadRepo: r}
}

func (s *Service) CreateUniqueFileName(header *multipart.FileHeader) types.UniqueFileName {
	// (my-file-name)
	fileBase := filepath.Base(strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename)))

	// (12345678)
	id := uuid.New().String()[0:8]

	// (.jpg)
	fileExt := filepath.Ext(header.Filename)

	// (12345678.jpg)
	idWithExt := fmt.Sprintf("%s%s", id, fileExt)

	// (my-file-name-12345678.jpg)
	filename := fmt.Sprintf("%s-%s%s", fileBase, id, fileExt)

	return types.UniqueFileName{
		Filename:  filename,
		ID:        id,
		IdWithExt: idWithExt,
		Type:      fileExt,
		Base:      fileBase,
	}
}

func (s *Service) SaveAssetImage(file multipart.File, name types.UniqueFileName) error {
	// ./public/my-file-name-12345678.jpg
	publicDir := "./public"
	if err := os.MkdirAll(publicDir, os.ModePerm); err != nil {
		return err
	}

	// Create empty file
	dst, err := os.Create(filepath.Join(publicDir, name.IdWithExt))
	if err != nil {
		return err
	}
	defer dst.Close()

	// Write the content from POST to the empty file
	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteImage(filename string) error {
	// ./public/my-file-name-12345678.jpg
	publicDir := "./public"
	if err := os.MkdirAll(publicDir, os.ModePerm); err != nil {
		return err
	}

	// Delete the file
	err := os.Remove(filepath.Join(publicDir, filename))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CheckFileType(header *multipart.FileHeader) error {
	allowed := false
	for _, ext := range ALLOWED_EXTENSIONS {
		if ext == filepath.Ext(header.Filename) {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("File type not allowed. Allowed types: %v", ALLOWED_EXTENSIONS)
	}

	return nil
}
