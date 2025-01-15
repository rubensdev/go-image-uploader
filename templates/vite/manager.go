package vite

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
	manifest     ViteManifest
	env          string
	devServerURL string
	mu           sync.RWMutex
}

func NewManager(jsonStr string, env string) (*Manager, error) {
	m := &Manager{
		env: env,
	}

	if env != "development" {
		var manifest ViteManifest
		err := json.Unmarshal([]byte(jsonStr), &manifest)
		if err != nil {
			return nil, err
		}
		m.manifest = manifest
	}

	return m, nil
}

func (m *Manager) GetAsset(path string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.InDevMode() {
		return m.devServerURL + path
	}

	if manifestEntry, ok := m.manifest[path]; ok {
		return manifestEntry.File
	}
	// TODO: Should panic here?
	panic(fmt.Errorf("asset %s not found", path))
}

func (m *Manager) GetEntry(entry string) (EntryProps, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if manifestEntry, ok := m.manifest[entry]; ok {
		return manifestEntry, nil
	}
	return EntryProps{}, errors.New("entry not found")
}

func (m *Manager) SetDevServerURL(host string) {
	m.devServerURL = fmt.Sprintf("http://%s:5173/", host)
}

func (m *Manager) InDevMode() bool {
	return m.env == "development"
}

// func checkResourceAvailable(url string) (bool, error) {
// 	res, err := http.Head(url)
// 	if err != nil {
// 		return false, err
// 	}
// 	defer res.Body.Close()

// 	return res.StatusCode != http.StatusNotFound, nil
// }
