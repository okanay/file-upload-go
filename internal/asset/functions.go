package asset

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/okanay/file-upload-go/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func OptimizeImage(publicDir, optimizedDir, filename string, quality int) error {
	// Check if directories exist
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		return fmt.Errorf("public directory does not exist: %s", publicDir)
	}
	if _, err := os.Stat(optimizedDir); os.IsNotExist(err) {
		return fmt.Errorf("optimized images directory does not exist: %s", optimizedDir)
	}

	inputPath := filepath.Join(publicDir, filename)
	if !utils.FileIsExist(inputPath) {
		return fmt.Errorf("input file does not exist: %s", inputPath)
	}

	optimizedFilename := CreateOptimizedFileName(filename, quality)
	outputPath := filepath.Join(optimizedDir, optimizedFilename)

	err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	if utils.FileIsExist(outputPath) {
		return nil // Dosya zaten var, işlem yapmaya gerek yok
	}

	ext := filepath.Ext(filename)
	var cmd *exec.Cmd

	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		ffmpegQuality := 31 - int(float64(quality)/100*29)
		cmd = exec.Command("ffmpeg",
			"-i", inputPath,
			"-q:v", strconv.Itoa(ffmpegQuality),
			"-y",
			outputPath,
		)
	case ".png":
		scale := float64(quality) / 100
		cmd = exec.Command("ffmpeg",
			"-i", inputPath,
			"-vf", fmt.Sprintf("scale=iw*%f:-1", scale),
			"-y",
			outputPath,
		)
	case ".webp":
		cmd = exec.Command("ffmpeg",
			"-i", inputPath,
			"-c:v", "libwebp",
			"-quality", strconv.Itoa(quality),
			"-y",
			outputPath,
		)
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %v, output: %s", err, string(output))
	}

	fmt.Println("[OPTIMIZE IMAGE] ", outputPath)
	return nil
}

func BlurImage(publicDir, blurDir, filename string) error {
	// Check if the directories exist
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		return errors.New("public directory does not exist")
	}
	if _, err := os.Stat(blurDir); os.IsNotExist(err) {
		return errors.New("blur directory does not exist")
	}

	normalPath := filepath.Join(publicDir, filename)
	blurredPath := filepath.Join(blurDir, filename)

	// Check if the original file exists
	if _, err := os.Stat(normalPath); os.IsNotExist(err) {
		return errors.New("original file does not exist")
	}

	// Check if the blurred image already exists
	if _, err := os.Stat(blurredPath); err == nil {
		return nil // Blurred image already exists, no need to recreate
	}

	srcImage, err := imaging.Open(normalPath)
	if err != nil {
		return err
	}

	blurredImage := imaging.Blur(srcImage, 20.0)
	resizeImage := imaging.Resize(blurredImage, 50, 0, imaging.Lanczos)

	err = imaging.Save(resizeImage, blurredPath)
	if err != nil {
		return err
	}

	fmt.Println("[BLUR IMAGE] ", blurredPath)
	return nil
}
