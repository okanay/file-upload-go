package asset

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/okanay/file-upload-go/types"
	"github.com/okanay/file-upload-go/utils"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type AssetHandler struct {
	PublicDir     string
	BlurDir       string
	OptimizedDir  string
	AutoClean     bool
	CleanInterval time.Duration
	mutex         sync.RWMutex
	service       *Service
}

func NewAssetHandler(s *Service, publicDir, blurDir, optimizedDir string, autoClean bool, cleanInterval time.Duration) *AssetHandler {
	handler := &AssetHandler{
		service:       s,
		PublicDir:     publicDir,
		BlurDir:       blurDir,
		OptimizedDir:  optimizedDir,
		AutoClean:     autoClean,
		CleanInterval: cleanInterval,
	}

	if autoClean {
		fmt.Println("[ASSET CLEANER] Asset cleaner is enabled.")
		go handler.autoCleanRoutine()
	}
	return handler
}

func (h *AssetHandler) autoCleanRoutine() {
	h.ClearDirectories()

	ticker := time.NewTicker(h.CleanInterval)
	defer ticker.Stop()

	for range ticker.C {
		h.ClearDirectories()
		fmt.Println("[ASSET CLEANER] Asset directories have been cleaned.")
	}
}

func (h *AssetHandler) ClearDirectories() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Clear Blur Directory
	if err := h.clearDirectory(h.BlurDir); err != nil {
		// Log error
	}

	// Clear Optimized Directory
	if err := h.clearDirectory(h.OptimizedDir); err != nil {
		// Log error
	}
}

func (h *AssetHandler) clearDirectory(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		fmt.Println("[ASSET CLEANER] Deleting", dir, name)
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	fmt.Println("[ASSET CLEANER] Directory", dir, "has been cleaned.")
	return nil
}

func (h *AssetHandler) GetAsset(c *gin.Context) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	filename := c.Param("filename")

	if !utils.FileIsExist(filepath.Join(h.PublicDir, filename)) {
		c.JSON(404, gin.H{"message": "The requested " + filename + " was not found."})
		return
	}

	if path := h.handleQualityOptimization(c, filename); path != "" {
		c.File(path)
		return
	}

	if path := h.handleBlur(c, filename); path != "" {
		c.File(path)
		return
	}

	if path := h.getOriginalFile(filename); path != "" {
		c.File(path)
		return
	}
}

func (h *AssetHandler) GetAllAssets(c *gin.Context) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	assets, err := h.service.repository.GetAllAssets()
	if err != nil {
		c.JSON(500, gin.H{"message": "Error fetching assets: " + err.Error()})
		return
	}

	if assets == nil || len(assets) == 0 {
		c.JSON(200, gin.H{"assets": []types.Assets{}})
		return
	}

	c.JSON(200, gin.H{"assets": assets})
}

func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// get filename from post body
	filename := c.PostForm("filename")

	if err := h.service.repository.DeleteAsset(filename); err != nil {
		c.JSON(500, gin.H{"message": "Error deleting asset: " + err.Error()})
		return
	}

	normalPath := filepath.Join(h.PublicDir, filename)
	if utils.FileIsExist(normalPath) {
		if err := os.Remove(normalPath); err != nil {
			c.JSON(500, gin.H{"message": "Error deleting asset: " + err.Error()})
			fmt.Println("[ERROR] Error deleting asset:", err)
			return
		}
	}

	fmt.Println("[ASSET DELETED] Asset has been deleted:", filename)
	c.JSON(200, gin.H{"message": "Asset has been deleted."})
}

func (h *AssetHandler) handleQualityOptimization(c *gin.Context, filename string) string {
	quality := c.Query("quality")
	if quality == "" {
		return ""
	}

	percentage, err := strconv.Atoi(quality)
	if err != nil || percentage < 1 || percentage > 100 {
		return ""
	}

	optimizedFilename := CreateOptimizedFileName(filename, percentage)
	optimizePath := filepath.Join(h.OptimizedDir, optimizedFilename)

	if utils.FileIsExist(optimizePath) {
		return optimizePath
	}

	err = OptimizeImage(h.PublicDir, h.OptimizedDir, filename, percentage)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error optimizing image: " + err.Error()})
		return ""
	}

	if utils.FileIsExist(optimizePath) {
		return optimizePath
	}

	return ""
}

func (h *AssetHandler) handleBlur(c *gin.Context, filename string) string {
	blur := c.Query("blur")
	if blur != "yes" {
		return ""
	}

	err := BlurImage(h.PublicDir, h.BlurDir, filename)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error processing the image: " + err.Error()})
		return ""
	}

	blurredPath := filepath.Join(h.BlurDir, filename)
	if _, err := os.Stat(blurredPath); err == nil {
		return blurredPath
	}

	return ""
}

func (h *AssetHandler) getOriginalFile(filename string) string {
	normalPath := filepath.Join(h.PublicDir, filename)
	if utils.FileIsExist(normalPath) {
		return normalPath
	}
	return ""
}
