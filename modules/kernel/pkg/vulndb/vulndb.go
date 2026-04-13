package vulndb

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/judozi/judozi/modules/kernel/pkg/kernel"
)

//go:embed exploits.json
var exploitsData []byte

type Exploit struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	MinKernel    string   `json:"min_kernel"`
	MaxKernel    string   `json:"max_kernel"`
	Arch         []string `json:"arch"`
	Source       string   `json:"source"`
	Binary       string   `json:"binary,omitempty"`
	Compile      string   `json:"compile"`
	Execute      string   `json:"execute"`
	Requirements []string `json:"requirements"`
	Tags         []string `json:"tags"`
	References   []string `json:"references"`
}

type Database struct {
	Exploits []Exploit `json:"exploits"`
}

func Load() (*Database, error) {
	var db Database
	if err := json.Unmarshal(exploitsData, &db); err != nil {
		return nil, fmt.Errorf("failed to parse exploit database: %w", err)
	}
	return &db, nil
}

func (db *Database) Search(kv kernel.Version, arch string) []Exploit {
	var matches []Exploit
	for _, e := range db.Exploits {
		minV, err := kernel.ParseVersion(e.MinKernel)
		if err != nil {
			continue
		}
		maxV, err := kernel.ParseVersion(e.MaxKernel)
		if err != nil {
			continue
		}
		if !kv.GTE(minV) || !kv.LTE(maxV) {
			continue
		}
		archOK := false
		for _, a := range e.Arch {
			if strings.EqualFold(a, arch) {
				archOK = true
				break
			}
		}
		if archOK {
			matches = append(matches, e)
		}
	}
	return matches
}

func (db *Database) GetByID(id string) *Exploit {
	for i := range db.Exploits {
		if strings.EqualFold(db.Exploits[i].ID, id) {
			return &db.Exploits[i]
		}
	}
	return nil
}
