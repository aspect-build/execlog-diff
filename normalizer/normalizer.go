package normalizer

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"
)

const nestedPrefix = "nested:"

// OPTIMIZATION: Use a struct for map keys to avoid string allocation
type MapKey struct {
	LogPathHash string
	Key         interface{}
}

type Normalizer struct {
	hash_to_file_path_map      *sync.Map
	file_path_to_file_hash_map *sync.Map
	file_path_to_file_size_map *sync.Map

	// Updated to use MapKey
	file_id_to_path_map       map[MapKey]string
	file_id_to_size_map       map[MapKey]int64
	input_set_id_to_files_map map[MapKey][]string

	// Runfiles tree ID -> the input set ID containing the runfiles
	runfiles_tree_id_to_input_set_id_map map[MapKey]int64

	hash_to_env_map *sync.Map

	ctx *context.Context
}

type normalizer interface {
	Run(types.ExecLogEntry) types.ExecLogEntry
}

func New(ctx *context.Context) *Normalizer {
	return &Normalizer{
		hash_to_file_path_map:                &sync.Map{},
		file_path_to_file_hash_map:           &sync.Map{},
		file_path_to_file_size_map:           &sync.Map{},
		file_id_to_path_map:                  make(map[MapKey]string),
		file_id_to_size_map:                  make(map[MapKey]int64),
		input_set_id_to_files_map:            make(map[MapKey][]string),
		runfiles_tree_id_to_input_set_id_map: make(map[MapKey]int64),
		hash_to_env_map:                      &sync.Map{},
		ctx:                                  ctx,
	}
}

func (n *Normalizer) Normalize(execLogChan <-chan types.ExecLogEntry, normalizedExecLogChan chan<- types.Spawn, waitGroup *sync.WaitGroup) {
	normalizers := []normalizer{}

	switch n.ctx.Mode {
	case context.Diff:
		normalizers = append(normalizers,
			HashEnvironmentValues(n),
			HashFilePaths(n),
			HashArgFilePaths(n),
			HashEnv(n),
			StoreFileSizes(n),
			StoreFiles(n),
			StoreRunfilesTrees(n),
			StoreInputSets(n),
			PopulateFiles(n),
			HashSpawn(n),
		)
	case context.ListLabels:
		normalizers = append(normalizers,
			HashSpawn(n),
		)
	}

	for execLog := range execLogChan {
		if !n.ctx.Continue() {
			break
		}

		for _, normalizer := range normalizers {
			execLog = normalizer.Run(execLog)
		}

		incomingSpawn := execLog.Spawn
		if incomingSpawn != nil {
			spawn := types.Spawn{
				LogPath:        incomingSpawn.LogPath,
				LogPathHash:    incomingSpawn.LogPathHash,
				Args:           incomingSpawn.Args,
				EnvHash:        incomingSpawn.EnvHash,
				Label:          incomingSpawn.Label,
				CacheHit:       incomingSpawn.CacheHit,
				Remoteable:     incomingSpawn.Remoteable,
				Cacheable:      incomingSpawn.Cacheable,
				RemoteCachable: incomingSpawn.RemoteCachable,
				Hash:           incomingSpawn.Hash,
				InputFiles:     incomingSpawn.InputFiles,
				OutputFiles:    incomingSpawn.OutputFiles,
			}

			switch n.ctx.Mode {
			case context.Diff:
				normalizedExecLogChan <- spawn
			case context.ListLabels:
				spawn = types.Spawn{
					Label:          incomingSpawn.Label,
					CacheHit:       incomingSpawn.CacheHit,
					Remoteable:     incomingSpawn.Remoteable,
					Cacheable:      incomingSpawn.Cacheable,
					RemoteCachable: incomingSpawn.RemoteCachable,
					Hash:           incomingSpawn.Hash,
				}
				normalizedExecLogChan <- spawn
			}
		}
	}

	waitGroup.Done()
}

func (n *Normalizer) GetFilePathFromFilePathHash(hash string) string {
	if val, ok := n.hash_to_file_path_map.Load(hash); ok {
		return val.(string)
	}

	return hash
}

func (n *Normalizer) GetFileHash(logPathHash string, key interface{}) string {
	hash := lookupId(logPathHash, key)
	if val, ok := n.file_path_to_file_hash_map.Load(hash); ok {
		return val.(string)
	}

	return ""
}

func (n *Normalizer) DeleteFileHash(logPathHash string, key interface{}) {
	lookupKey := lookupId(logPathHash, key)
	n.file_path_to_file_hash_map.Delete(lookupKey)
}

func (n *Normalizer) GetFileSize(logPathHash string, key interface{}) int64 {
	hash := lookupId(logPathHash, key)
	if val, ok := n.file_path_to_file_size_map.Load(hash); ok {
		return val.(int64)
	}

	return 0
}

func (n *Normalizer) DeleteFileSize(logPathHash string, key interface{}) {
	lookupKey := lookupId(logPathHash, key)
	n.file_path_to_file_size_map.Delete(lookupKey)
}

// OPTIMIZATION: Return struct instead of allocating string.
// Normalizes int→int64 so map lookups are type-consistent regardless of how the caller spells the literal.
func lookupId(logPathHash string, key interface{}) MapKey {
	if v, ok := key.(int); ok {
		key = int64(v)
	}
	return MapKey{
		LogPathHash: logPathHash,
		Key:         key,
	}
}

func nestedLookupId(logPathHash string, key interface{}) string {
	return fmt.Sprintf("%s%s:%v", nestedPrefix, logPathHash, key)
}

func (n *Normalizer) getInputsFromInputSet(key MapKey) []string {
	inputsMap := n.getInputsMap(key, make(map[string]struct{}), make(map[string]struct{}))

	inputsList := make([]string, 0, len(inputsMap))
	for input := range inputsMap {
		inputsList = append(inputsList, input)
	}

	sort.Strings(inputsList)

	return inputsList
}

func (n *Normalizer) getInputsMap(key MapKey, checked map[string]struct{}, inputs map[string]struct{}) map[string]struct{} {
	for _, input := range n.input_set_id_to_files_map[key] {
		if strings.HasPrefix(input, nestedPrefix) {
			// Format is "nested:LogPathHash:ID"
			// We need to parse this string back into a MapKey to recurse
			trimmed := input[len(nestedPrefix):]

			// Stop recursion if we've already checked this key
			if _, isChecked := checked[trimmed]; !isChecked {
				checked[trimmed] = struct{}{}

				// Split "LogPathHash:ID"
				if hash, idStr, found := strings.Cut(trimmed, ":"); found {
					if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
						nestedKey := MapKey{LogPathHash: hash, Key: id}
						n.getInputsMap(nestedKey, checked, inputs)
					}
				}
			}
		} else {
			inputs[input] = struct{}{}
		}
	}

	return inputs
}

func (n *Normalizer) getEnv(hash string) map[string]string {
	env, _ := n.hash_to_env_map.Load(hash)
	return env.(map[string]string)
}
