package manifest

import (
	"encoding/json"
	"errors"
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
	viteServerURL string
	mu            sync.RWMutex
}

func NewManager(path string, viteServerUrl string) (*Manager, error) {
	mm := &Manager{viteServerURL: viteServerUrl}
	err := mm.Load(path)
	return mm, err
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
