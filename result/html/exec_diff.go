package html

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/aspect-build/execlog-diff/result/proto"
)

type nodeData struct {
	nodes        []*proto.ActionGraphDiffNode
	targetLabel  string
	totalSize    int64
	rebuildCount int
	isCacheable  bool
}

func (h *html) printExecDiff(execDiff *proto.ExecDiff) {
	sort.Slice(execDiff.Nodes, func(i, j int) bool {
		return execDiff.Nodes[i].TargetLabel < execDiff.Nodes[j].TargetLabel
	})

	// Group nodes by target label and aggregate data
	targetMap := make(map[string]*nodeData)
	for _, node := range execDiff.Nodes {
		size := h.calculateOutputSize(node)
		rebuilds := TransitiveInfluenceCount(node.ActionHash, execDiff.Nodes)

		h.totalBytes += size
		if node.RemoteCachable {
			h.totalCacheableBytes += size
		}

		if data, exists := targetMap[node.TargetLabel]; exists {
			// Aggregate data for existing target
			data.nodes = append(data.nodes, node)
			data.totalSize += size
			if rebuilds > data.rebuildCount {
				data.rebuildCount = rebuilds
			}
			if node.RemoteCachable {
				data.isCacheable = true
			}
		} else {
			// Create new entry for target
			targetMap[node.TargetLabel] = &nodeData{
				nodes:        []*proto.ActionGraphDiffNode{node},
				targetLabel:  node.TargetLabel,
				totalSize:    size,
				rebuildCount: rebuilds,
				isCacheable:  node.RemoteCachable,
			}
		}
	}

	// Convert map to sorted slice
	nodesData := make([]nodeData, 0, len(targetMap))
	for _, data := range targetMap {
		nodesData = append(nodesData, *data)
	}
	sort.Slice(nodesData, func(i, j int) bool {
		return nodesData[i].targetLabel < nodesData[j].targetLabel
	})

	// Generate HTML
	html := h.generateHTML(nodesData, execDiff.Nodes)

	// Write to file
	file, err := os.Create(h.outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating HTML file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(html)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing HTML file: %v\n", err)
		return
	}

	fmt.Printf("HTML report generated: %s\n", h.outputPath)
}

func (h *html) calculateOutputSize(node *proto.ActionGraphDiffNode) int64 {
	var total int64
	for _, change := range node.Changes {
		switch item := change.Item.(type) {
		case *proto.InputItem_OutputFileListChanges:
			total += item.OutputFileListChanges.SizeBytes
		}
	}
	return total
}

func (h *html) generateHTML(nodesData []nodeData, allNodes []*proto.ActionGraphDiffNode) string {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ExecLog Diff Cache Misses Report</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: #1e1e1e;
            color: #d4d4d4;
            padding: 20px;
        }
        .container {
            max-width: 1400px;
            margin: 0 auto;
        }
        h1 {
            color: #4ec9b0;
            margin-bottom: 10px;
            font-size: 28px;
        }
        .summary {
            background: #252526;
            border: 1px solid #3e3e42;
            border-radius: 4px;
            padding: 15px 20px;
            margin-bottom: 20px;
        }
        .summary-item {
            margin: 8px 0;
            font-size: 14px;
        }
        .summary-label {
            color: #4ec9b0;
            font-weight: 500;
        }
        .filter-box {
            margin-bottom: 20px;
        }
        .filter-input {
            width: 100%;
            padding: 12px 16px;
            background: #252526;
            border: 1px solid #3e3e42;
            border-radius: 4px;
            color: #d4d4d4;
            font-size: 14px;
            font-family: inherit;
        }
        .filter-input:focus {
            outline: none;
            border-color: #4ec9b0;
        }
        .filter-input::placeholder {
            color: #858585;
        }
        .target-list {
            list-style: none;
        }
        .target-item {
            background: #252526;
            border: 1px solid #3e3e42;
            border-radius: 4px;
            margin-bottom: 10px;
            overflow: hidden;
            transition: border-color 0.2s;
        }
        .target-item:hover {
            border-color: #4ec9b0;
        }
        .target-header {
            padding: 15px 20px;
            cursor: pointer;
            display: flex;
            align-items: center;
            gap: 12px;
            user-select: none;
        }
        .target-header:hover {
            background: #2a2d2e;
        }
        .expand-icon {
            color: #858585;
            font-size: 12px;
            transition: transform 0.2s;
            min-width: 12px;
        }
        .target-item.expanded .expand-icon {
            transform: rotate(90deg);
        }
        .target-label {
            color: #4ec9b0;
            font-weight: 500;
            flex-grow: 1;
        }
        .target-meta {
            display: flex;
            gap: 16px;
            font-size: 13px;
        }
        .meta-item {
            color: #858585;
        }
        .meta-size {
            color: #569cd6;
        }
        .meta-rebuilds {
            color: #ce9178;
        }
        .meta-cacheable {
            color: #c586c0;
            font-weight: 500;
        }
        .target-details {
            display: none;
            padding: 0 20px 20px 20px;
            border-top: 1px solid #3e3e42;
        }
        .target-item.expanded .target-details {
            display: block;
        }
        .section {
            margin-top: 16px;
        }
        .section-title {
            color: #d4d4d4;
            font-weight: 500;
            margin-bottom: 8px;
            font-size: 14px;
        }
        .file-list {
            list-style: none;
            margin-left: 20px;
        }
        .file-item {
            padding: 4px 0;
            font-size: 13px;
            color: #d4d4d4;
        }
        .file-path {
            color: #d4d4d4;
        }
        .file-size {
            color: #569cd6;
            margin-left: 8px;
        }
        .dependency-item {
            padding: 8px 0;
            border-bottom: 1px solid #3e3e42;
        }
        .dependency-item:last-child {
            border-bottom: none;
        }
        .dependency-label {
            color: #dcdcaa;
            font-weight: 500;
            margin-bottom: 4px;
        }
        .no-results {
            text-align: center;
            padding: 40px 20px;
            color: #858585;
            font-size: 14px;
        }
        .warning {
            color: #c586c0;
            font-weight: 500;
            margin-left: 8px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>ExecLog Diff Cache Misses Report</h1>

        <div class="summary">
            <div class="summary-item">
                <span class="summary-label">Total invalidated output size:</span> ` + h.formatBytes(h.totalBytes) + `
            </div>
            <div class="summary-item">
                <span class="summary-label">Total invalidated cacheable output size:</span> ` + h.formatBytes(h.totalCacheableBytes) + `
            </div>
            <div class="summary-item">
                <span class="summary-label">Total targets:</span> ` + strconv.Itoa(len(nodesData)) + `
            </div>
        </div>

        <div class="filter-box">
            <input type="text" class="filter-input" id="filterInput" placeholder="Filter targets by label...">
        </div>

        <ul class="target-list" id="targetList">
`

	for _, data := range nodesData {
		html += h.generateTargetHTML(data, allNodes)
	}

	html += `        </ul>
        <div class="no-results" id="noResults" style="display: none;">
            No targets match your filter
        </div>
    </div>

    <script>
        // Filter functionality
        const filterInput = document.getElementById('filterInput');
        const targetList = document.getElementById('targetList');
        const noResults = document.getElementById('noResults');
        const targetItems = document.querySelectorAll('.target-item');

        filterInput.addEventListener('input', (e) => {
            const filter = e.target.value.toLowerCase();
            let visibleCount = 0;

            targetItems.forEach(item => {
                const label = item.dataset.label.toLowerCase();
                if (label.includes(filter)) {
                    item.style.display = '';
                    visibleCount++;
                } else {
                    item.style.display = 'none';
                }
            });

            if (visibleCount === 0) {
                targetList.style.display = 'none';
                noResults.style.display = 'block';
            } else {
                targetList.style.display = '';
                noResults.style.display = 'none';
            }
        });

        // Expand/collapse functionality
        document.querySelectorAll('.target-header').forEach(header => {
            header.addEventListener('click', () => {
                header.parentElement.classList.toggle('expanded');
            });
        });
    </script>
</body>
</html>`

	return html
}

func (h *html) generateTargetHTML(data nodeData, allNodes []*proto.ActionGraphDiffNode) string {
	html := `            <li class="target-item" data-label="` + h.escapeHTML(data.targetLabel) + `">
                <div class="target-header">
                    <span class="expand-icon">▶</span>
                    <span class="target-label">` + h.escapeHTML(data.targetLabel) + `</span>
                    <div class="target-meta">
`

	if data.totalSize > 0 {
		html += `                        <span class="meta-item meta-size">` + h.formatBytes(data.totalSize) + `</span>
`
	}

	if data.rebuildCount > 0 {
		html += `                        <span class="meta-item meta-rebuilds">` + strconv.Itoa(data.rebuildCount) + ` rebuilds</span>
`
	}

	if data.isCacheable {
		html += `                        <span class="meta-item meta-cacheable">cacheable</span>
`
	}

	html += `                    </div>
                </div>
                <div class="target-details">
`

	// Output changes - aggregate from all nodes
	html += h.generateOutputChangesHTML(data.nodes)

	// Input changes - aggregate from all nodes
	html += h.generateInputChangesHTML(data.nodes)

	// Changed dependencies - aggregate from all nodes
	html += h.generateDependenciesHTML(data.nodes, allNodes)

	html += `                </div>
            </li>
`

	return html
}

func (h *html) generateOutputChangesHTML(nodes []*proto.ActionGraphDiffNode) string {
	// Collect all output changes from all nodes
	type outputFile struct {
		path string
		size int64
	}
	filesMap := make(map[string]int64)

	for _, node := range nodes {
		for _, change := range node.Changes {
			switch item := change.Item.(type) {
			case *proto.InputItem_OutputFileListChanges:
				// Aggregate sizes if the same file appears multiple times
				filesMap[item.OutputFileListChanges.Path] += item.OutputFileListChanges.SizeBytes
			}
		}
	}

	if len(filesMap) == 0 {
		return ""
	}

	// Convert to slice and sort
	files := make([]outputFile, 0, len(filesMap))
	for path, size := range filesMap {
		files = append(files, outputFile{path: path, size: size})
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].path < files[j].path
	})

	html := `                    <div class="section">
                        <div class="section-title">Changed Output Files</div>
                        <ul class="file-list">
`

	for _, file := range files {
		html += `                            <li class="file-item">
                                <span class="file-path">` + h.escapeHTML(file.path) + `</span>
                                <span class="file-size">(` + h.formatBytes(file.size) + `)</span>
                            </li>
`
	}

	html += `                        </ul>
                    </div>
`

	return html
}

func (h *html) generateInputChangesHTML(nodes []*proto.ActionGraphDiffNode) string {
	// Collect all input changes from all nodes
	type inputFile struct {
		path string
		size int64
	}
	filesMap := make(map[string]int64)

	for _, node := range nodes {
		for _, change := range node.Changes {
			switch item := change.Item.(type) {
			case *proto.InputItem_InputFileListChanges:
				// Aggregate sizes if the same file appears multiple times
				filesMap[item.InputFileListChanges.Path] += item.InputFileListChanges.SizeBytes
			}
		}
	}

	if len(filesMap) == 0 {
		return ""
	}

	// Convert to slice and sort
	files := make([]inputFile, 0, len(filesMap))
	for path, size := range filesMap {
		files = append(files, inputFile{path: path, size: size})
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].path < files[j].path
	})

	html := `                    <div class="section">
                        <div class="section-title">Changed Input Files</div>
                        <ul class="file-list">
`

	for _, file := range files {
		html += `                            <li class="file-item">
                                <span class="file-path">` + h.escapeHTML(file.path) + `</span>
                                <span class="file-size">(` + h.formatBytes(file.size) + `)</span>
                            </li>
`
	}

	html += `                        </ul>
                    </div>
`

	return html
}

func (h *html) generateDependenciesHTML(nodes []*proto.ActionGraphDiffNode, allNodes []*proto.ActionGraphDiffNode) string {
	// Collect all dependencies from all nodes
	type dependencyInfo struct {
		targetLabel  string
		changedFiles map[string]bool
	}
	depsMap := make(map[string]*dependencyInfo)

	for _, node := range nodes {
		for _, child := range node.Children {
			targetLabel := GetTarget(child.Child, allNodes)
			if targetLabel == "" {
				continue
			}

			if dep, exists := depsMap[targetLabel]; exists {
				// Add changed files to existing dependency
				for _, file := range child.ChangedFiles {
					switch item := file.Item.(type) {
					case *proto.InputItem_InputFileListChanges:
						dep.changedFiles[item.InputFileListChanges.Path] = true
					}
				}
			} else {
				// Create new dependency entry
				dep := &dependencyInfo{
					targetLabel:  targetLabel,
					changedFiles: make(map[string]bool),
				}
				for _, file := range child.ChangedFiles {
					switch item := file.Item.(type) {
					case *proto.InputItem_InputFileListChanges:
						dep.changedFiles[item.InputFileListChanges.Path] = true
					}
				}
				depsMap[targetLabel] = dep
			}
		}
	}

	if len(depsMap) == 0 {
		return ""
	}

	// Convert to slice and sort
	deps := make([]*dependencyInfo, 0, len(depsMap))
	for _, dep := range depsMap {
		deps = append(deps, dep)
	}
	sort.Slice(deps, func(i, j int) bool {
		return deps[i].targetLabel < deps[j].targetLabel
	})

	html := `                    <div class="section">
                        <div class="section-title">Changed Dependencies</div>
`

	for _, dep := range deps {
		html += `                        <div class="dependency-item">
                            <div class="dependency-label">` + h.escapeHTML(dep.targetLabel) + `</div>
`

		if len(dep.changedFiles) > 0 {
			// Sort the changed files
			files := make([]string, 0, len(dep.changedFiles))
			for file := range dep.changedFiles {
				files = append(files, file)
			}
			sort.Strings(files)

			html += `                            <ul class="file-list">
`
			for _, file := range files {
				html += `                                <li class="file-item">
                                    <span class="file-path">` + h.escapeHTML(file) + `</span>
                                </li>
`
			}
			html += `                            </ul>
`
		}

		html += `                        </div>
`
	}

	html += `                    </div>
`

	return html
}

func (h *html) formatBytes(b int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case b >= GB:
		return fmt.Sprintf("%.2f GB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.2f MB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.2f KB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d Bytes", b)
	}
}

func (h *html) escapeHTML(s string) string {
	s = replaceAll(s, "&", "&amp;")
	s = replaceAll(s, "<", "&lt;")
	s = replaceAll(s, ">", "&gt;")
	s = replaceAll(s, "\"", "&quot;")
	s = replaceAll(s, "'", "&#39;")
	return s
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		if i <= len(s)-len(old) && s[i:i+len(old)] == old {
			result += new
			i += len(old) - 1
		} else {
			result += string(s[i])
		}
	}
	return result
}

// How many times is this node depended on directly or transitively
func TransitiveInfluenceCount(nodeActionHash string, nodes []*proto.ActionGraphDiffNode) int {
	count := 0
	for _, node := range nodes {
		if node.ActionHash != nodeActionHash {
			for _, child := range node.Children {
				if child.Child == nodeActionHash {
					count++
					count = count + TransitiveInfluenceCount(node.ActionHash, nodes)
				}
			}
		}
	}

	return count
}

func GetTarget(nodeActionHash string, nodes []*proto.ActionGraphDiffNode) string {
	for _, node := range nodes {
		if node.ActionHash == nodeActionHash {
			return node.TargetLabel
		}
	}

	return ""
}
