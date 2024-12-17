package templates

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

type ManifestManager struct {
	manifest      ViteManifest
	viteServerURL string
	mu            sync.RWMutex
}

func NewManifestManager(path string, viteServerUrl string) (*ManifestManager, error) {
	mm := &ManifestManager{viteServerURL: viteServerUrl}
	err := mm.Load(path)
	return mm, err
}

func (mm *ManifestManager) Load(path string) error {
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

func (mm *ManifestManager) GetAsset(entry string) string {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if manifestEntry, ok := mm.manifest[entry]; ok {
		return manifestEntry.File
	}
	return mm.viteServerURL + entry
}

func (mm *ManifestManager) GetEntry(entry string) (EntryProps, error) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if manifestEntry, ok := mm.manifest[entry]; ok {
		return manifestEntry, nil
	}
	return EntryProps{}, errors.New("entry not found")
}
