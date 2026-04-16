package models

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

const (
	// FindingStatusNew indicates a finding that has not been seen in DefectDojo
	// before and is therefore newly discovered in the current scan run.
	FindingStatusNew = "NEW"
	// FindingStatusDuplicated indicates a finding that already exists as an active
	// finding in DefectDojo from a previous scan run (a known vulnerability).
	FindingStatusDuplicated = "DUPLICATED"
)

// Finding represents a single security finding produced by a scanner.
type Finding struct {
	Engine         string
	Severity       string
	Name           string
	VulnId         string
	Cwe            string
	Description    string
	SinkFile       string
	SinkLine       int
	Recommendation string
	// Hash is a deterministic content hash used for stable cross-scanner matching
	// against DefectDojo findings. It is computed by ComputeFindingHash when the
	// finding is loaded from a scanner output and compared against the equivalent
	// hash derived from the DefectDojo finding fields in MarkFindingsByActiveFindings.
	Hash string
	// Status classifies the finding as either NEW (not previously seen in DefectDojo)
	// or DUPLICATED (already tracked as an active finding in DefectDojo). Without
	// --sync all findings default to NEW.
	Status string
}

// ComputeFindingHash returns a deterministic SHA-256 hex hash over the finding
// fields that are reliably preserved when a scan result is imported into and then
// read back from DefectDojo. The hash is therefore computable independently from
// both the local (scanner) side and the DefectDojo API side and used as the primary
// matching key in MarkFindingsByActiveFindings.
//
// The same formula is used for all scanners:
//
//	hash(lower(severity) | lower(sinkFile) | sinkLine | lower(recommendation))
//
// Scanner-specific notes:
//   - Grype:    recommendation is the "Upgrade to X" string derived from fix.versions.
//   - OpenGrep: recommendation is always "" because DefectDojo's Semgrep parser stores
//               extra.message in description, not mitigation. The hash is additionally
//               injected into extra.fingerprint before upload so that DefectDojo stores
//               it as unique_id_from_tool, enabling a direct lookup without recomputation.
//   - KICS:     recommendation is the expected_value from each file entry.
func ComputeFindingHash(severity, sinkFile string, sinkLine int, recommendation string) string {
	input := strings.ToLower(strings.TrimSpace(severity)) + "|" +
		strings.ToLower(strings.TrimSpace(sinkFile)) + "|" +
		strconv.Itoa(sinkLine) + "|" +
		strings.ToLower(strings.TrimSpace(recommendation))
	h := sha256.Sum256([]byte(input))
	return hex.EncodeToString(h[:])
}
