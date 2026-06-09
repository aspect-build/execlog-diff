package sawmill

import (
	"github.com/aspect-build/execlog-diff/result/proto"
	"github.com/aspect-build/execlog-diff/types"
)

func InputFiles(sawmill *SawMill) differ {
	return &input_file_differ{
		sawmill: sawmill,
	}
}

type input_file_differ struct {
	sawmill *SawMill
}

func (d *input_file_differ) Run(spawn1 *types.Spawn, spawn2 *types.Spawn) []*proto.InputItem {
	diffResults := make([]*proto.InputItem, 0)

	spawn1Files := make([]string, 0)
	for _, path := range spawn1.InputFiles {
		hash := d.sawmill.normalizer.GetFileHash(spawn1.LogPathHash, path)
		if hash != "" {
			spawn1Files = append(spawn1Files, path)
		}
	}

	spawn2Files := make([]string, 0)
	for _, path := range spawn2.InputFiles {
		hash := d.sawmill.normalizer.GetFileHash(spawn2.LogPathHash, path)
		if hash != "" {
			spawn2Files = append(spawn2Files, path)
		}
	}

	for i1, filepath := range spawn1Files {
		i2 := stringListFind(spawn2Files, filepath, i1)
		if i2 == -1 {
			fileSize1 := d.sawmill.normalizer.GetFileSize(spawn1.LogPathHash, filepath)
			fileSize2 := d.sawmill.normalizer.GetFileSize(spawn2.LogPathHash, filepath)
			diffResults = append(diffResults, &proto.InputItem{
				Item: &proto.InputItem_InputFileListChanges{
					InputFileListChanges: &proto.InputFileChange{
						Path:      d.sawmill.normalizer.GetFilePathFromFilePathHash(filepath),
						SizeBytes: max(fileSize1, fileSize2),
					},
				},
				Reason: proto.Reason_REMOVED,
			})
		} else if i1 == i2 {
			// File is found in both spawns. Check that they have the same content hash
			file1hash := d.sawmill.normalizer.GetFileHash(spawn1.LogPathHash, filepath)
			file2hash := d.sawmill.normalizer.GetFileHash(spawn2.LogPathHash, filepath)
			if file1hash != file2hash {
				fileSize1 := d.sawmill.normalizer.GetFileSize(spawn1.LogPathHash, filepath)
				fileSize2 := d.sawmill.normalizer.GetFileSize(spawn2.LogPathHash, filepath)
				diffResults = append(diffResults, &proto.InputItem{
					Item: &proto.InputItem_InputFileListChanges{
						InputFileListChanges: &proto.InputFileChange{
							Path:      d.sawmill.normalizer.GetFilePathFromFilePathHash(filepath),
							SizeBytes: max(fileSize1, fileSize2),
						},
					},
					Reason: proto.Reason_CHANGED,
				})
			}
		}
	}

	for i2, filepath := range spawn2Files {
		i1 := stringListFind(spawn1Files, filepath, i2)
		if i1 == -1 {
			fileSize1 := d.sawmill.normalizer.GetFileSize(spawn1.LogPathHash, filepath)
			fileSize2 := d.sawmill.normalizer.GetFileSize(spawn2.LogPathHash, filepath)
			diffResults = append(diffResults, &proto.InputItem{
				Item: &proto.InputItem_InputFileListChanges{
					InputFileListChanges: &proto.InputFileChange{
						Path:      d.sawmill.normalizer.GetFilePathFromFilePathHash(filepath),
						SizeBytes: max(fileSize1, fileSize2),
					},
				},
				Reason: proto.Reason_ADDED,
			})
		}
	}

	return diffResults
}
