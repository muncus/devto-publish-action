// Syncer allows for updating of existing dev.to articles from a filesystem.

package devto

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

type Syncer struct {
	StateFile string
	IDMap     map[string]int
	StateMap  map[string]SyncerStateRecord
	client    *Client
}

type SyncerStateRecord struct {
	Id int `json:"id"`
}

// Used for testing
func NewEmptySyncer() *Syncer {
	s := &Syncer{
		IDMap:    make(map[string]int),
		StateMap: make(map[string]SyncerStateRecord),
	}
	return s
}

func NewSyncer(statefile string, apikey string) (*Syncer, error) {
	s := &Syncer{
		StateFile: statefile,
		client:    NewClient(apikey),
		IDMap:     make(map[string]int),
		StateMap:  make(map[string]SyncerStateRecord),
	}
	if s.StateFile != "" {
		err := s.LoadStateFile(s.StateFile)
		if err != nil {
			s.IDMap = make(map[string]int)
			return nil, err
		}
	}
	return s, nil
}

// SetDebug enables debugging on the dev.to client used to sync.
func (s *Syncer) SetDebug(d bool) {
	s.client.Debug = d
}

func (s *Syncer) LoadStateFile(file string) error {
	b, err := ioutil.ReadFile(file)
	if err == os.ErrNotExist {
		// File specified does not exist. use empty state.
		s.StateMap = make(map[string]SyncerStateRecord)
		return nil
	}
	if err != nil {
		return err
	}
	return s.LoadState(b)
}
func (s *Syncer) LoadState(content []byte) error {
	// fmt.Printf("%#v\n", s)
	// try to unmarshal into "simple" format, with just IDs
	err := json.Unmarshal(content, &s.IDMap)
	if err == nil {
		// old format seen, convert to SyncerStateRecords.
		for k, v := range s.IDMap {
			s.StateMap[k] = SyncerStateRecord{
				Id: v,
			}
		}
		return nil
	} else {
		// Attempt to unmarshal "new" format.
		err = json.Unmarshal(content, &s.StateMap)
		fmt.Printf("%#v\n", s.StateMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Syncer) DumpState(file string) error {
	b, err := os.Create(file)
	if err != nil {
		return err
	}
	defer b.Close()
	encoder := json.NewEncoder(b)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(&s.StateMap)
	if err != nil {
		return err
	}
	return nil
}

func (s *Syncer) Sync(basedir string, logger zerolog.Logger) ([]*Article, []error) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(basedir)
	files, err := filepath.Glob(filepath.Join(basedir, "*"))
	if err != nil {
		return []*Article{}, []error{err}
	}
	errors := make([]error, 0)
	articles := make([]*Article, 0)
	for _, f := range files {
		time.Sleep(3 * time.Second)
		relpath, err := filepath.Rel(basedir, f)
		if err != nil {
			logger.Error().Str("file", relpath).Msg(err.Error())
			errors = append(errors, err)
			continue
		}
		a, err := s.SyncFile(relpath)
		if err != nil {
			logger.Error().Str("file", relpath).Msg(err.Error())
			errors = append(errors, err)
			continue
		}
		logger.Info().Str("file", relpath).Msg("Success")
		articles = append(articles, a)
	}
	return articles, errors
}

func (s *Syncer) SyncFile(file string) (*Article, error) {
	b, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer b.Close()
	bodycontent, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	a := &Article{
		BodyMarkdown: string(bodycontent),
	}
	if _, ok := s.IDMap[file]; ok {
		if s.IDMap[file] < 0 {
			//This file was previously skipped.
			return nil, fmt.Errorf("Skipping file with previous failures: %s", file)
		}
		// This file has been posted before. update.
		a.ID = s.IDMap[file]
	}
	newArticle, err := s.client.UpsertArticle(a, nil)
	if err != nil {
		s.IDMap[file] = -1
		return nil, fmt.Errorf("Failed to update file '%s': %s", file, err)
	}
	// Update the state map, so we can Update it next time.
	s.IDMap[file] = newArticle.ID
	return newArticle, nil
}
