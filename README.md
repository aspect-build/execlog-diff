# execlog-diff

![status](https://img.shields.io/badge/status-pre--release-orange)

A tool for detecting non-hermetic builds by comparing two Bazel execution logs.

## What it does

Bazel's remote cache relies on hermetic, reproducible builds. When the same action produces different outputs across runs, cache misses occur and build times suffer. execlog-diff helps you find the root cause by comparing two compact execution logs and surfacing actions whose inputs, outputs, or environment changed between runs.

## Generating execution logs

You need two compact execution logs from separate clean builds. Run both builds with the following flags to avoid polluting your remote cache during this analysis:

```
--noremote_accept_cached
--noremote_upload_local_results
```

Along with one of these (depending on your Bazel version):

```
--experimental_execution_log_compact_file=compact.exec.log
--execution_log_compact_file=compact.exec.log
```

A full example for a single build:

```sh
bazel clean --expunge
bazel build //... \
  --noremote_accept_cached \
  --noremote_upload_local_results \
  --execution_log_compact_file=compact.exec.log
```

Run this twice (saving each log separately, e.g. `compact.exec.1.log` and `compact.exec.2.log`) to get the pair of logs needed for comparison.

Using `--expunge` is recommended so that any external dependencies that vary across machines or warming jobs are caught rather than masked by a warm local cache.

## Running the tool

```sh
bazel run //:bin -- explain \
  /path/to/compact.exec.1.log \
  /path/to/compact.exec.2.log \
  --cache-misses \
  --html \
  --html-output /path/to/results.html
```

Open the generated HTML file to see which actions changed between runs and why.

### Flags

| Flag | Description |
|------|-------------|
| `--cache-misses` | Find actions that changed between the two logs |
| `--html` | Generate HTML output (requires `--cache-misses`) |
| `--html-output <path>` | Path for the HTML file (implies `--html`) |
| `--list` | List all labels in a single execlog |
