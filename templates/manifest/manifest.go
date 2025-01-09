package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

func NewManager(path string, env string, devHost string) (*Manager, error) {
	mm := &Manager{
		env:           env,
		viteServerURL: fmt.Sprintf("http://%s:5173/", devHost),
	}

	if env != "development" {
		err := mm.Load(path)
		return mm, err
	}

	return mm, nil
}

func (mm *Manager) Load(path string) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var manifest ViteManifest
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return err
	}

	mm.manifest = manifest
	return nil
}

func (mm *Manager) GetAsset(entry string) string {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if manifestEntry, ok := mm.manifest[entry]; ok {
		return manifestEntry.File
	}
	return mm.viteServerURL + entry
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
