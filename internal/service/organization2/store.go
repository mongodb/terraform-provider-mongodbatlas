package organization2

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	storeMu     sync.RWMutex
	store       map[string]*orgState
	storeLoaded sync.Once
)

type persistedOrgState struct {
	NextRenewal      time.Time `json:"next_renewal"`
	ExpiresAt        time.Time `json:"expires_at"`
	SecretCreatedAt  time.Time `json:"secret_created_at"`
	Name             string    `json:"name"`
	OrgID            string    `json:"org_id"`
	ClientID         string    `json:"client_id"`
	ClientSecret     string    `json:"client_secret"`
	Interval         string    `json:"interval,omitempty"`
	CurrentSecretID  string    `json:"current_secret_id"`
	OldSecretID      string    `json:"old_secret_id,omitempty"`
	SecretVersion    int64     `json:"secret_version"`
	HasRotationBlock bool      `json:"has_rotation_block"`
}

func storePath() string {
	if path := os.Getenv("MONGODB_ATLAS_ORGANIZATION2_POC_STORE"); path != "" {
		return path
	}
	return filepath.Join(os.TempDir(), "mongodbatlas-organization2-poc-store.json")
}

func ensureStoreLoaded() {
	storeLoaded.Do(func() {
		store = map[string]*orgState{}
		if err := loadStoreFromDisk(); err != nil {
			store = map[string]*orgState{}
		}
	})
}

func loadStoreFromDisk() error {
	data, err := os.ReadFile(storePath())
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	var persisted map[string]persistedOrgState
	if err := json.Unmarshal(data, &persisted); err != nil {
		return err
	}

	loaded := make(map[string]*orgState, len(persisted))
	for name := range persisted {
		entry := persisted[name]
		loaded[name] = persistedToOrgState(&entry)
	}
	store = loaded
	return nil
}

func saveStoreToDisk() error {
	persisted := make(map[string]persistedOrgState, len(store))
	for name, entry := range store {
		persisted[name] = orgStateToPersisted(entry)
	}

	data, err := json.MarshalIndent(persisted, "", "  ")
	if err != nil {
		return err
	}

	path := storePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	tempPath := path + ".tmp"
	if err := os.WriteFile(tempPath, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tempPath, path)
}

func getStoreEntry(name string) (*orgState, bool) {
	ensureStoreLoaded()
	storeMu.RLock()
	defer storeMu.RUnlock()
	entry, ok := store[name]
	return entry, ok
}

func putStoreEntry(name string, state *orgState) error {
	ensureStoreLoaded()
	storeMu.Lock()
	store[name] = state
	err := saveStoreToDisk()
	storeMu.Unlock()
	return err
}

func deleteStoreEntry(name string) error {
	ensureStoreLoaded()
	storeMu.Lock()
	delete(store, name)
	err := saveStoreToDisk()
	storeMu.Unlock()
	return err
}

func resetStoreLocked() {
	storeMu.Lock()
	store = map[string]*orgState{}
	storeLoaded = sync.Once{}
	_ = os.Remove(storePath())
	storeMu.Unlock()
}

func orgStateToPersisted(state *orgState) persistedOrgState {
	return persistedOrgState{
		Name:             state.name,
		OrgID:            state.orgID,
		ClientID:         state.clientID,
		ClientSecret:     state.clientSecret,
		Interval:         state.interval,
		SecretVersion:    state.secretVersion,
		NextRenewal:      state.nextRenewal,
		ExpiresAt:        state.expiresAt,
		CurrentSecretID:  state.currentSecretID,
		OldSecretID:      state.oldSecretID,
		SecretCreatedAt:  state.secretCreatedAt,
		HasRotationBlock: state.hasRotationBlock,
	}
}

func persistedToOrgState(state *persistedOrgState) *orgState {
	return &orgState{
		name:             state.Name,
		orgID:            state.OrgID,
		clientID:         state.ClientID,
		clientSecret:     state.ClientSecret,
		interval:         state.Interval,
		secretVersion:    state.SecretVersion,
		nextRenewal:      state.NextRenewal,
		expiresAt:        state.ExpiresAt,
		currentSecretID:  state.CurrentSecretID,
		oldSecretID:      state.OldSecretID,
		secretCreatedAt:  state.SecretCreatedAt,
		hasRotationBlock: state.HasRotationBlock,
	}
}

func updateStoreEntry(name string, update func(*orgState) error) error {
	ensureStoreLoaded()
	storeMu.Lock()
	defer storeMu.Unlock()

	entry, ok := store[name]
	if !ok {
		return fmt.Errorf("organization2 %q not found", name)
	}
	if err := update(entry); err != nil {
		return err
	}
	return saveStoreToDisk()
}
