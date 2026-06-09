package types

type IncomingSpawn struct {
	LogPath        string
	LogPathHash    string
	Id             int64
	Args           []string
	Env            map[string]string
	EnvHash        string
	Label          string
	CacheHit       bool
	Remoteable     bool
	Cacheable      bool
	RemoteCachable bool
	Hash           string
	InputFiles     []string
	OutputFiles    []string
	InputSetId     int64
	OutputFileIds  []int64
	OutputSetIds   []int64

	// TODO
	// - platform
	// - runner: Have seen "disk cache hit" && "darwin sandbox"
	// - mnemonic: Might give a clue as to what was run
	// - metrics: Could be useful to show how long an invalidation branch took to run
}

// Have a trimmed down type to save space when saving a fully formed spawn
type Spawn struct {
	LogPath        string
	LogPathHash    string
	Args           []string
	EnvHash        string
	Label          string
	CacheHit       bool
	Remoteable     bool
	Cacheable      bool
	RemoteCachable bool
	Hash           string
	InputFiles     []string
	OutputFiles    []string
}
