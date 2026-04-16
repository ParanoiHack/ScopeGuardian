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
// vulnId is the vulnerability identifier (e.g. CVE-2021-1234 for Grype findings).
// For scanners that do not emit a per-finding vulnerability ID (KICS, Opengrep) it
// should be passed as an empty string. On the DefectDojo side, the finding Title is
// used as the vulnId because DefectDojo stores the CVE ID as the title for Grype
// findings, while storing a category-prefixed rule name for KICS findings. By
// computing two hashes per DD finding — one with the title as vulnId and one with an
// empty vulnId — both Grype findings and KICS/Opengrep findings are matched correctly
// without any title-parsing logic.
func ComputeFindingHash(vulnId, severity, sinkFile string, sinkLine int, recommendation string) string {
	input := strings.ToLower(strings.TrimSpace(vulnId)) + "|" +
		strings.ToLower(strings.TrimSpace(severity)) + "|" +
		strings.ToLower(strings.TrimSpace(sinkFile)) + "|" +
		strconv.Itoa(sinkLine) + "|" +
		strings.ToLower(strings.TrimSpace(recommendation))
	h := sha256.Sum256([]byte(input))
	return hex.EncodeToString(h[:])
}
