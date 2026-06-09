package terminal

import (
	"os"
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/result/proto"
	"google.golang.org/protobuf/encoding/protojson"
	"gotest.tools/v3/golden"

	. "github.com/onsi/gomega"
)

func TestContainsFile(t *testing.T) {
	t.Run("generate golden files", func(t *testing.T) {
		g := NewGomegaWithT(t)

		files := []string{
			"cli_proto_result.txt",
			"cli_list_proto_result.txt",
		}

		for _, file := range files {
			// Run it & check the golden files many times to ensure they are consistent given consistent input
			for i := 0; i < 100; i++ {
				b, err := os.ReadFile("./testdata/" + file)
				if err != nil {
					t.Error(err)
				}

				result := &proto.ExecLogDiffResult{}
				protojson.Unmarshal(b, result)

				var stdout strings.Builder
				var stderr strings.Builder
				term := New(&stdout, &stderr)

				diagnostics := term.Print(result.Result)

				golden.Assert(t, stdout.String(), "golden_files/"+file)

				g.Expect(stderr.String()).To(Equal(""))
				g.Expect(len(diagnostics)).To(Equal(0))
			}
		}
	})
}
