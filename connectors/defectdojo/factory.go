package defectdojo

import "ScopeGuardian/connectors/defectdojo/client"

// GetDefectDojoService constructs and returns a DefectDojoService configured with
// the provided HTTP client, API base URL, and access token.
func GetDefectDojoService(client client.Client, url string, accessToken string) DefectDojoService {
	return newDefectDojoService(client, url, accessToken)
}
