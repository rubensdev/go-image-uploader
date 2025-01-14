package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type EntryProps struct {
	File     string   `json:"file"`
	CssFiles []string `json:"css,omitempty"`
	Src      string   `json:"src"`
}

type ViteManifest map[string]EntryProps

type Manager struct {
	manifest      ViteManifest
	env           string
	viteServerURL string
	mu            sync.RWMutex
}

func NewManager(jsonStr string, env string, devHost string) (*Manager, error) {
	mm := &Manager{
		env:           env,
		viteServerURL: fmt.Sprintf("http://%s:5173/", devHost),
	}

	if env != "development" {
		var manifest ViteManifest
		err := json.Unmarshal([]byte(jsonStr), &manifest)
		if err != nil {
			return nil, err
		}
		mm.manifest = manifest
	}

	return mm, nil
}

func (mm *Manager) GetAsset(path string) string {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if mm.env == "development" {
		return mm.viteServerURL + path
	}

	if manifestEntry, ok := mm.manifest[path]; ok {
		return manifestEntry.File
	}
	// TODO: Should panic here?
	panic(fmt.Errorf("asset %s not found", path))
}

func (mm *Manager) GetEntry(entry string) (EntryProps, error) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if manifestEntry, ok := mm.manifest[entry]; ok {
		return manifestEntry, nil
	}
	return EntryProps{}, errors.New("entry not found")
}

func (mm *Manager) ViteServerURL() string {
	return mm.viteServerURL
}

func (mm *Manager) IsDevelopment() bool {
	return mm.env == "development"
}
