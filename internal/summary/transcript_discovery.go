package summary

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DiscoverTranscriptCandidates returns transcript candidate paths sorted by filename.
func DiscoverTranscriptCandidates(meetingDir string) ([]string, error) {
	entries, err := os.ReadDir(meetingDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read meeting directory %s: %w", meetingDir, err)
	}

	candidates := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.EqualFold(filepath.Ext(entry.Name()), ".txt") {
			candidates = append(candidates, filepath.Join(meetingDir, entry.Name()))
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return filepath.Base(candidates[i]) < filepath.Base(candidates[j])
	})

	return candidates, nil
}

// FindSingleTranscriptCandidate resolves transcript source with the 0/1/many contract.
func FindSingleTranscriptCandidate(meetingDir string) (string, error) {
	candidates, err := DiscoverTranscriptCandidates(meetingDir)
	if err != nil {
		return "", err
	}

	switch len(candidates) {
	case 0:
		return "", fmt.Errorf("no transcript candidate found in %s; expected exactly one .txt file", meetingDir)
	case 1:
		return candidates[0], nil
	default:
		names := make([]string, 0, len(candidates))
		for _, candidate := range candidates {
			names = append(names, filepath.Base(candidate))
		}
		return "", fmt.Errorf(
			"multiple transcript candidates found in %s: %s",
			meetingDir,
			strings.Join(names, ", "),
		)
	}
}
