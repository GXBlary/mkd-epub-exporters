package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"mkd-epub-exporters/pkg/converter"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ProgressData holds the status of active conversion jobs
type ProgressData struct {
	Current int    `json:"current"`
	Total   int    `json:"total"`
	Name    string `json:"name"`
}

// ConversionResult holds the outcome of a file conversion
type ConversionResult struct {
	SourcePath string `json:"SourcePath"`
	DestPath   string `json:"DestPath"`
	Success    bool   `json:"Success"`
	Error      string `json:"Error"`
}

// SelectFiles opens a file selection dialog and returns the selected paths
func (a *App) SelectFiles() ([]string, error) {
	files, err := wailsRuntime.OpenMultipleFilesDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select Documents to Convert",
		Filters: []wailsRuntime.FileFilter{
			{
				DisplayName: "All Supported Documents",
				Pattern:     "*.docx;*.xlsx;*.pdf;*.html;*.htm;*.txt;*.md;*.pptx;*.rtf;*.epub;*.odt;*.tex;*.wiki",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// SelectFolder opens a directory selection dialog and returns the path
func (a *App) SelectFolder() (string, error) {
	folder, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select Output Folder",
	})
	if err != nil {
		return "", err
	}
	return folder, nil
}

// OpenFile opens a file using the host OS default application handler
func (a *App) OpenFile(path string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "start", "", path)
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", path)
	} else {
		cmd = exec.Command("xdg-open", path)
	}
	_ = cmd.Start()
}

// Convert converts all queued files into the specified format
func (a *App) Convert(files []string, outDir string, format string, embedImages bool) ([]ConversionResult, error) {
	allFiles, err := converter.CollectFiles(files)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	results := make([]ConversionResult, 0, len(allFiles))
	exportMD := format == "md" || format == "both"
	exportEpub := format == "epub" || format == "both"

	for i, fPath := range allFiles {
		// Emit progress to Frontend
		wailsRuntime.EventsEmit(a.ctx, "conversion-progress", ProgressData{
			Current: i + 1,
			Total:   len(allFiles),
			Name:    filepath.Base(fPath),
		})

		mdContent, err := converter.ConvertFile(fPath, outDir, embedImages)
		if err != nil {
			results = append(results, ConversionResult{
				SourcePath: fPath,
				Success:    false,
				Error:      err.Error(),
			})
			continue
		}

		fileSuccess := true
		var lastOutPath string

		if exportMD {
			stem := strings.TrimSuffix(filepath.Base(fPath), filepath.Ext(fPath))
			outPath := filepath.Join(outDir, stem+".md")

			counter := 1
			for {
				if _, err := os.Stat(outPath); os.IsNotExist(err) {
					break
				}
				outPath = filepath.Join(outDir, fmt.Sprintf("%s_%d.md", stem, counter))
				counter++
			}

			err = os.WriteFile(outPath, []byte(mdContent), 0644)
			if err != nil {
				fileSuccess = false
			} else {
				lastOutPath = outPath
				results = append(results, ConversionResult{
					SourcePath: fPath,
					DestPath:   outPath,
					Success:    true,
				})
			}
		}

		if exportEpub {
			stem := strings.TrimSuffix(filepath.Base(fPath), filepath.Ext(fPath))
			outPath := filepath.Join(outDir, stem+".epub")

			counter := 1
			for {
				if _, err := os.Stat(outPath); os.IsNotExist(err) {
					break
				}
				outPath = filepath.Join(outDir, fmt.Sprintf("%s_%d.epub", stem, counter))
				counter++
			}

			err = converter.ConvertToEpub(mdContent, stem, outPath)
			if err != nil {
				fileSuccess = false
			} else {
				lastOutPath = outPath
				// If we exported both, avoid duplicate entries in the successful list
				if !exportMD {
					results = append(results, ConversionResult{
						SourcePath: fPath,
						DestPath:   outPath,
						Success:    true,
					})
				}
			}
		}

		if !fileSuccess {
			results = append(results, ConversionResult{
				SourcePath: fPath,
				Success:    false,
				Error:      "Failed to write exported files",
			})
		} else if exportMD && exportEpub {
			// If we successfully wrote both files, return the path of the last generated file
			// (usually the EPUB) as target, which is fine for double-click opening.
			// Let's update the existing entry to show the last output path.
			for idx := range results {
				if results[idx].SourcePath == fPath && results[idx].Success {
					results[idx].DestPath = lastOutPath
				}
			}
		}
	}

	return results, nil
}
