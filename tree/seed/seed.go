package seed

import (
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/types"
)

func NewSeed(ctx *context.Context, normalizer *normalizer.Normalizer) *Seed {
	return &Seed{
		ctx:        ctx,
		Normalizer: normalizer,

		log_path_hash_to_output_file_and_vertex_map: &sync.Map{},
		spawn_hash_to_vertex_map:                    &sync.Map{},
	}
}

type Seed struct {
	ctx        *context.Context
	Normalizer *normalizer.Normalizer

	// Each incoming execlog needs its own "output file" to "vertex" mapping.
	// Store these individual maps in a different map so we can easily select
	// the map we need to use based on the incoming execlog path
	// Key: log path hash
	// Value: another map
	//     Nested Key: output file path
	//     Nested Value: reference to the vertex that created the output file
	log_path_hash_to_output_file_and_vertex_map *sync.Map

	// key: "Hash" value from a given spawn
	// value: a given "vertex"
	spawn_hash_to_vertex_map *sync.Map
}

func (s *Seed) SetVertex(spawnHash string, val *types.Vertex) {
	s.spawn_hash_to_vertex_map.Store(spawnHash, val)
}

func (s *Seed) GetVertex(spawnHash string) *types.Vertex {
	val, ok := s.spawn_hash_to_vertex_map.Load(spawnHash)
	if ok {
		return val.(*types.Vertex)
	} else {
		return nil
	}
}

func (s *Seed) GetHashToVertexMap() *sync.Map {
	return s.spawn_hash_to_vertex_map
}

func (s *Seed) DeleteVertex(spawnHash string) {
	s.spawn_hash_to_vertex_map.Delete(spawnHash)
}

func (s *Seed) SetOutputFileVertex(logPathHash string, outputFilePath string, vertex *types.Vertex) {
	temp, present := s.log_path_hash_to_output_file_and_vertex_map.Load(logPathHash)
	var outputFileMap *sync.Map

	if present {
		outputFileMap = temp.(*sync.Map)
	} else {
		outputFileMap = &sync.Map{}
		s.log_path_hash_to_output_file_and_vertex_map.Store(logPathHash, outputFileMap)
	}

	if _, ok := outputFileMap.Load(outputFilePath); !ok {
		outputFileMap.Store(outputFilePath, vertex)
	}
}

func (s *Seed) GetOutputFileVertex(logPathHash string, outputFilePath string) *types.Vertex {
	temp, present := s.log_path_hash_to_output_file_and_vertex_map.Load(logPathHash)
	if !present {
		return nil
	}

	outputFileMap := temp.(*sync.Map)

	val, ok := outputFileMap.Load(outputFilePath)
	if ok {
		return val.(*types.Vertex)
	} else {
		return nil
	}
}

func (s *Seed) DeleteOutputFileVertex(logPathHash string, outputFilePath string) {
	temp, present := s.log_path_hash_to_output_file_and_vertex_map.Load(logPathHash)
	if !present {
		return
	}

	outputFileMap := temp.(*sync.Map)

	outputFileMap.Delete(outputFilePath)
}
