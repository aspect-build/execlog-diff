package types

// RunfilesTree represents the .runfiles directory of an executable.
// It contains references to input sets that hold the actual files.
type RunfilesTree struct {
	Id          int64
	LogPathHash string
	Path        string
	// The input set containing files at canonical locations in the tree
	InputSetId int64
	// The symlink entry set for custom symlinks (relative to workspace dir)
	SymlinksId int64
	// The symlink entry set for root symlinks
	RootSymlinksId int64
}
