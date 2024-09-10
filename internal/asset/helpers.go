package asset

import (
	"fmt"
	"path/filepath"
	"strings"
)

func CreateOptimizedFileName(filename string, quality int) string {
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)
	return fmt.Sprintf("%s-%d%s", baseName, quality, ext)
}
