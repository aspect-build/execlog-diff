package execlogdiff

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/result/proto"
	"github.com/aspect-build/execlog-diff/result/terminal"
	"github.com/aspect-build/execlog-diff/tree"
	"github.com/aspect-build/execlog-diff/types"
	"github.com/aspect-build/execlog-diff/diagnostic"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/matchers"
)

func TestExecLogDiff(t *testing.T) {
	t.Run("surfaces errors opening the first file and stops the conversion of the second file", func(t *testing.T) {
		g := NewGomegaWithT(t)

		// Run this 1000 times to try and weed out any potential async issues
		for i := 0; i < 1000; i++ {
			var stderr strings.Builder
			ctx := context.New(&stderr, context.Diff)

			n := normalizer.New(ctx)
			f := execLogDiff{
				ctx:         ctx,
				execlogChan: make(chan types.ExecLogEntry, 100),
				spawnChan:   make(chan types.Spawn, 100),
				normalizer:  n,
				diffTree:    tree.NewDiffTree(ctx, n),
			}

			filepath1 := "./some_nonexistant_file_123456.log"
			filepath2 := "./test-data/changed_file_7.3.0_1.compact.log"

			execDiff := f.diffExecLogs(filepath1, filepath2)

			g.Expect(stderr.String()).To(HavePrefix("ERROR: Failed to open file"))
			g.Expect(stderr.String()).To(ContainSubstring(filepath1))
			g.Expect(stderr.String()).To(ContainSubstring("no such file or directory"))
			g.Expect(ctx.Continue()).To(BeFalse())
			g.Expect(len(execDiff.RuntimeDiagnostics)).To(Equal(1))
			g.Expect(execDiff.RuntimeDiagnostics[0].Severity).To(Equal(diagnostic.Severity_ERROR))
		}
	})

	t.Run("surfaces errors opening the second file and stops the conversion of the first file", func(t *testing.T) {
		g := NewGomegaWithT(t)

		// Run this 1000 times to try and weed out any potential async issues
		for i := 0; i < 1000; i++ {
			var stderr strings.Builder
			ctx := context.New(&stderr, context.Diff)

			n := normalizer.New(ctx)
			f := execLogDiff{
				ctx:         ctx,
				execlogChan: make(chan types.ExecLogEntry, 100),
				spawnChan:   make(chan types.Spawn, 100),
				normalizer:  n,
				diffTree:    tree.NewDiffTree(ctx, n),
			}

			filepath1 := "./test-data/changed_file_7.3.0_1.compact.log"
			filepath2 := "./some_nonexistant_file_123456.log"

			execDiff := f.diffExecLogs(filepath1, filepath2)

			g.Expect(stderr.String()).To(HavePrefix("ERROR: Failed to open file"))
			g.Expect(stderr.String()).To(ContainSubstring(filepath2))
			g.Expect(stderr.String()).To(ContainSubstring("no such file or directory"))
			g.Expect(ctx.Continue()).To(BeFalse())
			g.Expect(len(execDiff.RuntimeDiagnostics)).To(Equal(1))
			g.Expect(execDiff.RuntimeDiagnostics[0].Severity).To(Equal(diagnostic.Severity_ERROR))
		}
	})

	t.Run("surfaces errors opening both files", func(t *testing.T) {
		g := NewGomegaWithT(t)

		// Run this 1000 times to try and weed out any potential async issues
		for i := 0; i < 1000; i++ {
			var stderr strings.Builder
			ctx := context.New(&stderr, context.Diff)

			n := normalizer.New(ctx)
			f := execLogDiff{
				ctx:         ctx,
				execlogChan: make(chan types.ExecLogEntry),
				spawnChan:   make(chan types.Spawn),
				normalizer:  n,
				diffTree:    tree.NewDiffTree(ctx, n),
			}

			filepath1 := "./some_nonexistant_file_123456.log"
			filepath2 := "./some_other_nonexistant_file_123456.log"

			execDiff := f.diffExecLogs(filepath1, filepath2)

			g.Expect(stderr.String()).To(HavePrefix("ERROR: Failed to open file"))
			g.Expect(stderr.String()).To(ContainSubstring(filepath1))
			g.Expect(stderr.String()).To(ContainSubstring(filepath2))
			g.Expect(stderr.String()).To(ContainSubstring("no such file or directory"))
			g.Expect(ctx.Continue()).To(BeFalse())
			g.Expect(len(execDiff.RuntimeDiagnostics)).To(Equal(2))
			g.Expect(execDiff.RuntimeDiagnostics[0].Severity).To(Equal(diagnostic.Severity_ERROR))
			g.Expect(execDiff.RuntimeDiagnostics[1].Severity).To(Equal(diagnostic.Severity_ERROR))
			g.Expect(len(f.execlogChan)).To(Equal(0))
		}
	})

	t.Run("smoketest", func(t *testing.T) {
		g := NewGomegaWithT(t)

		// Run this 1000 times to try and weed out any potential async issues
		for i := 0; i < 10000; i++ {
			fmt.Println("Running number:", i)
			var stderr strings.Builder
			ctx := context.New(&stderr, context.Diff)

			execlogChan := make(chan types.ExecLogEntry, 5)
			spawnChan := make(chan types.Spawn, 1000)
			n := normalizer.New(ctx)
			f := execLogDiff{
				ctx:         ctx,
				execlogChan: execlogChan,
				spawnChan:   spawnChan,
				normalizer:  n,
				diffTree:    tree.NewDiffTree(ctx, n),
			}

			filepath1 := "./test-data/changed_file_7.3.0_1.compact.log"
			filepath2 := "./test-data/changed_file_7.3.0_2.compact.log"

			execDiff := f.diffExecLogs(filepath1, filepath2)

			g.Expect(stderr.String()).To(Equal(""))
			g.Expect(len(execDiff.RuntimeDiagnostics)).To(Equal(0))

			g.Expect(reflect.TypeOf(execDiff.Result.Result)).To(Equal(reflect.TypeOf(&proto.Result_ExecDiff{})))

			nodes := make([]proto.ActionGraphDiffNode, 0)
			for _, node := range execDiff.Result.GetExecDiff().Nodes {
				// Unfortunately, the matchers used below do not like pointers
				nodes = append(nodes, proto.ActionGraphDiffNode{
					ActionHash:  node.ActionHash,
					TargetLabel: node.TargetLabel,
				})
			}

			g.Expect(len(nodes)).To(Equal(2))

			g.Expect(nodes).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"ActionHash":  &matchers.EqualMatcher{Expected: "9366b90f37a5eeaed3528b11539487d8a23781e87522601f5bb26f3cd02350916405d622b2400f333ff95931aaf3eaf3"},
				"TargetLabel": &matchers.EqualMatcher{Expected: "//foo:foo"},
			})))

			g.Expect(nodes).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"ActionHash":  &matchers.EqualMatcher{Expected: "887da361e121526fc152cbf27796a099b4521636d63984bf849533ab492399237ad227397567501982f4a2af6d435cab"},
				"TargetLabel": &matchers.EqualMatcher{Expected: "//bar:bar"},
			})))

			g.Expect(execDiff.Result.GetExecDiff().RootNodes).To(Equal([]string{"887da361e121526fc152cbf27796a099b4521636d63984bf849533ab492399237ad227397567501982f4a2af6d435cab"}))
		}
	})

	t.Run("results are the same across bazel versions for the same changes", func(t *testing.T) {
		g := NewGomegaWithT(t)

		execlog_set := [][][]string{
			[][]string{
				[]string{
					"./test-data/changed_file_7.3.0_1.compact.log",
					"./test-data/changed_file_7.3.0_2.compact.log",
				},
				[]string{
					"./test-data/changed_file_7.3.1_1.compact.log",
					"./test-data/changed_file_7.3.1_2.compact.log",
				},
				[]string{
					"./test-data/changed_file_8.0.0_1.compact.log",
					"./test-data/changed_file_8.0.0_2.compact.log",
				},
				[]string{
					"./test-data/changed_file_9.0.0-pre.20250307.1_1.compact.log",
					"./test-data/changed_file_9.0.0-pre.20250307.1_2.compact.log",
				},
			},
			[][]string{
				[]string{
					"./test-data/changed_directory_7.3.0_1.compact.log",
					"./test-data/changed_directory_7.3.0_2.compact.log",
				},
				[]string{
					"./test-data/changed_directory_7.3.1_1.compact.log",
					"./test-data/changed_directory_7.3.1_2.compact.log",
				},
				[]string{
					"./test-data/changed_directory_8.0.0_1.compact.log",
					"./test-data/changed_directory_8.0.0_2.compact.log",
				},
				[]string{
					"./test-data/changed_directory_9.0.0-pre.20250307.1_1.compact.log",
					"./test-data/changed_directory_9.0.0-pre.20250307.1_2.compact.log",
				},
			},
		}

		// Run this 10 times to ensure we dont get a false positive on the execlogs matching
		for i := 0; i < 10; i++ {
			for _, change := range execlog_set {
				diffs := make([]*proto.ExecLogDiffResult, 0)
				for _, execlogs := range change {

					var stderr strings.Builder
					ctx := context.New(&stderr, context.Diff)

					execlogChan := make(chan types.ExecLogEntry, 5)
					spawnChan := make(chan types.Spawn, 1000)
					n := normalizer.New(ctx)
					f := execLogDiff{
						ctx:         ctx,
						execlogChan: execlogChan,
						spawnChan:   spawnChan,
						normalizer:  n,
						diffTree:    tree.NewDiffTree(ctx, n),
					}

					execDiff := f.diffExecLogs(execlogs[0], execlogs[1])

					diffs = append(diffs, execDiff)
				}

				for _, diff := range diffs[1:] {
					var stdout1 strings.Builder
					var stderr1 strings.Builder
					term1 := terminal.New(&stdout1, &stderr1)
					term1.Print(diffs[0].Result)

					var stdout2 strings.Builder
					var stderr2 strings.Builder
					term2 := terminal.New(&stdout2, &stderr2)
					term2.Print(diff.Result)

					g.Expect(stdout1.String()).To(Equal(stdout2.String()))
				}
			}
		}
	})
}

func TestListLabels(t *testing.T) {
	t.Run("surfaces errors opening file", func(t *testing.T) {
		g := NewGomegaWithT(t)

		filepath := "./some_nonexistant_file_123456.log"

		result := ListLabels(filepath)

		g.Expect(len(result.RuntimeDiagnostics)).To(Equal(1))
		g.Expect(result.RuntimeDiagnostics[0].Severity).To(Equal(diagnostic.Severity_ERROR))
		g.Expect(result.RuntimeDiagnostics[0].Message).To(ContainSubstring(filepath))
		g.Expect(result.RuntimeDiagnostics[0].Message).To(ContainSubstring("no such file or directory"))
	})

	t.Run("lists labels from execlog", func(t *testing.T) {
		g := NewGomegaWithT(t)

		filepath := "./test-data/changed_file_7.3.0_1.compact.log"

		result := ListLabels(filepath)

		g.Expect(len(result.RuntimeDiagnostics)).To(Equal(0))
		g.Expect(reflect.TypeOf(result.Result.Result)).To(Equal(reflect.TypeOf(&proto.Result_LabelQuery{})))

		labelQuery := result.Result.GetLabelQuery()
		g.Expect(labelQuery).NotTo(BeNil())
		g.Expect(len(labelQuery.Labels)).To(Equal(4))

		// Convert to non-pointer slice for matchers
		labels := make([]proto.Label, 0)
		for _, label := range labelQuery.Labels {
			labels = append(labels, proto.Label{
				TargetLabel: label.TargetLabel,
				CacheHit:    label.CacheHit,
				Remoteable:  label.Remoteable,
				Cacheable:   label.Cacheable,
			})
		}

		// Check that expected labels are present
		g.Expect(labels).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"TargetLabel": &matchers.EqualMatcher{Expected: "//foo:foo"},
		})))
		g.Expect(labels).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"TargetLabel": &matchers.EqualMatcher{Expected: "//bar:bar"},
		})))
		g.Expect(labels).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"TargetLabel": &matchers.EqualMatcher{Expected: "//baz:baz"},
		})))
		g.Expect(labels).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"TargetLabel": &matchers.EqualMatcher{Expected: "//:static"},
		})))
	})
}
