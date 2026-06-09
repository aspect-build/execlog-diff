# Tooling

Always use bazel for running tools or testing.

# Post Changes Instructions

After making any changes to the repo, ensure that the following commands pass:
- bazel test //...
- bazel run //result/proto:write_pb_go
- bazel run //diagnostic:write_pb_go

Disabling tests is not an option. The tests should always pass, if you broke a test, you broke something. If you made a legitimate change that should result in a change in a test, that is fine, but an unexpected test change is not acceptable