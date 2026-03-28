package defectdojo

import "scope-guardian/connectors/defectdojo/client"

func GetDefectDojoService(client client.Client, url string, accessToken string) DefectDojoService {
	return newDefectDojoService(client, url, accessToken)
}
