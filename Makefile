.PHONY: mockgen
mockgen:
	@find ~+ -type f -name "*.go" -print0 | xargs -0 grep -l "^//go:generate" | sort -u | xargs -L1 -P $$(nproc) go generate
