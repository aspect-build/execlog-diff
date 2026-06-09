package types

type ExecLogEntry struct {
	File         *File
	InputSet     *InputSet
	Spawn        *IncomingSpawn
	RunfilesTree *RunfilesTree
}
