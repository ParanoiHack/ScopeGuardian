package models

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

const (
	// FindingStatusActive indicates a finding that is active in DefectDojo
	// (not a duplicate, not suppressed). Without --sync all findings default to ACTIVE.
	FindingStatusActive = "ACTIVE"
	// FindingStatusInactive indicates a finding that DefectDojo has suppressed,
	// marked as a false positive, or accepted as a risk. The local scanner still
	// reports it but it is excluded from security-gate evaluation.
	FindingStatusInactive = "INACTIVE"
	// FindingStatusDuplicate indicates a finding that DefectDojo's deduplication
	// engine has identified as a duplicate of another finding in the product.
	// Duplicate findings are excluded from security-gate evaluation.
	FindingStatusDuplicate = "DUPLICATE"
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
	// hash derived from the DefectDojo finding fields in MarkFindingsByDDFindings.
	Hash string
	// Status reflects the DefectDojo state of this finding: ACTIVE (not suppressed,
	// not a duplicate), INACTIVE (suppressed / false-positive / accepted risk), or
	// DUPLICATE (DefectDojo's deduplication engine identified it as a duplicate of
	// another finding in the product). Without --sync all findings are ACTIVE.
	Status string
}

// FilterFindingsByStatus returns a new slice containing only findings whose
// Status is present in the allowed set. The comparison is case-insensitive.
// If statuses is empty, the original slice is returned unchanged.
func FilterFindingsByStatus(findings []Finding, statuses []string) []Finding {
	if len(statuses) == 0 {
		return findings
	}
	allowed := make(map[string]bool, len(statuses))
	for _, s := range statuses {
		allowed[strings.ToUpper(s)] = true
	}
	result := make([]Finding, 0, len(findings))
	for _, f := range findings {
		if allowed[strings.ToUpper(f.Status)] {
			result = append(result, f)
		}
	}
	return result
}

// ComputeFindingHash returns a deterministic SHA-256 hex hash over the finding
// fields that are reliably preserved when a scan result is imported into and then
// read back from DefectDojo. The hash is therefore computable independently from
// both the local (scanner) side and the DefectDojo API side and used as the primary
// matching key in MarkFindingsByDDFindings.
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
