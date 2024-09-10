package upload

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/okanay/file-upload-go/types"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	defer file.Close()

	// Check if file bigger than max upload size
	if header.Size > MAX_UPLOAD_SIZE {
		c.JSON(404, gin.H{"error": "Max upload size exceeded."})
		return
	}

	// Check if file extension is allowed
	err = h.service.CheckFileType(header)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type: " + err.Error()})
		return
	}

	uniqueFileName := h.service.CreateUniqueFileName(header)
	// Create record for database
	assetReq := types.CreateAssetReq{
		Creator:     "admin",
		Name:        uniqueFileName.ID,
		Type:        uniqueFileName.Type,
		Filename:    uniqueFileName.IdWithExt,
		Description: c.PostForm("description"),
		Size:        header.Size,
	}

	// Save file
	err = h.service.SaveAssetImage(file, uniqueFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file: " + err.Error()})
		return
	}

	// Save record to database
	asset, err := h.service.uploadRepo.CreateAssetRecord(assetReq)
	if err != nil {
		err = h.service.DeleteImage(uniqueFileName.IdWithExt)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating asset: " + err.Error()})
		return
	}

	fmt.Println("[UPLOAD ASSET] Asset created: ", asset)
	// Return response
	c.JSON(http.StatusOK, gin.H{
		"asset": asset,
	})
}
