package reposted

import (
	"encoding/gob"
	"os"
	"sync"
)

var DBFiles = map[string]any{
	"repostedHashes.db":   &ImgHashes,
	"repostedScores.db":   &Scores,
	"repostedLastPost.db": &LastPosts,
}

var SaveMutex = sync.Mutex{}
var LoadMutex = sync.Mutex{}

func SaveDB() error {
	SaveMutex.Lock()
	defer SaveMutex.Unlock()
	tempFiles := make(map[string]*os.File, len(DBFiles))
	err := func() error {
		// Create Temporary files
		for fileName := range DBFiles {
			f, err := os.CreateTemp(".", fileName)
			if err != nil {
				return err
			}
			tempFiles[fileName] = f
			defer tempFiles[fileName].Close()
		}
		// Encode to Temporary files
		for fileName, data := range DBFiles {
			enc := gob.NewEncoder(tempFiles[fileName])
			err := enc.Encode(data)
			if err != nil {
				return err
			}
		}
		return nil
	}()
	if err != nil {
		return err
	}
	// Now Move the temporary files over the old DBs
	for fileName := range DBFiles {
		err := os.Rename(tempFiles[fileName].Name(), fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadDB() error {
	LoadMutex.Lock()
	defer LoadMutex.Unlock()
	for fileName, data := range DBFiles {
		f, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer f.Close()
		dec := gob.NewDecoder(f)
		err = dec.Decode(data)
		if err != nil {
			return err
		}
	}
	return nil
}
