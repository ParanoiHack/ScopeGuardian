package models

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
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
	// hash derived from the DefectDojo finding fields in FilterByActiveFindings.
	Hash string
}

// ComputeFindingHash returns a deterministic SHA-256 hex hash over the finding
// fields that are reliably preserved when a scan result is imported into and then
// read back from DefectDojo. The hash is therefore computable independently from
// both the local (scanner) side and the DefectDojo API side and used as the primary
// matching key in FilterByActiveFindings.
//
// vulnId is the vulnerability identifier stored in DefectDojo's vulnerability_ids
// array for this finding. Each scanner populates it differently:
//   - Grype:    the CVE or GHSA identifier (e.g. "CVE-2021-1234"), which DefectDojo's
//               Anchore Grype parser reads directly from vulnerability.id.
//   - OpenGrep: the Semgrep check_id (e.g. "go.lang.security.injection.sql"), injected
//               into extra.metadata.cve before upload so DefectDojo's Semgrep parser
//               stores it in vulnerability_ids.
//   - KICS:     an empty string, because DefectDojo's KICS parser does not populate
//               vulnerability_ids; KICS findings are matched via the empty-vulnId hash.
func ComputeFindingHash(vulnId, severity, sinkFile string, sinkLine int, recommendation string) string {
	input := strings.ToLower(strings.TrimSpace(vulnId)) + "|" +
		strings.ToLower(strings.TrimSpace(severity)) + "|" +
		strings.ToLower(strings.TrimSpace(sinkFile)) + "|" +
		strconv.Itoa(sinkLine) + "|" +
		strings.ToLower(strings.TrimSpace(recommendation))
	h := sha256.Sum256([]byte(input))
	return hex.EncodeToString(h[:])
}
