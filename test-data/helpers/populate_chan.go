package helpers

import (
	"github.com/aspect-build/execlog-diff/types"
)

type populate_chan struct {
	execlogChan chan types.ExecLogEntry
	counter     int64
}

func (t *populate_chan) GetLog1Path() string {
	return "path1"
}

func (t *populate_chan) GetLog2Path() string {
	return "path2"
}

func (t *populate_chan) Log1File(file string, hash string) {
	t.execlogChan <- types.ExecLogEntry{
		File: &types.File{
			Id:          t.counter,
			Path:        file,
			LogPathHash: t.GetLog1Path(),
			Hash:        hash,
		},
	}

	t.counter++
}

func (t *populate_chan) Log2File(file string, hash string) {
	t.execlogChan <- types.ExecLogEntry{
		File: &types.File{
			Id:          t.counter,
			Path:        file,
			LogPathHash: t.GetLog2Path(),
			Hash:        hash,
		},
	}

	t.counter++
}

func NewPopulateHelper(execlogChan chan types.ExecLogEntry) *populate_chan {
	return &populate_chan{
		execlogChan: execlogChan,
		counter:     1,
	}
}
