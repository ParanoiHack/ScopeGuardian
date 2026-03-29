package defectdojo

type Product struct {
	Id int `json:"id"`
}

type Engagement struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Branch    string `json:"branch_tag"`
	TargetEnd string `json:"target_end"`
}

type GetProductByNameResponse struct {
	Count   int       `json:"count"`
	Results []Product `json:"results"`
}

type GetEngagementsResponse struct {
	Count   int          `json:"count"`
	Results []Engagement `json:"results"`
}

type CreateEngagementResponse struct {
	Id int `json:"id"`
}

type EngagementPayload struct {
	Tags                      []string `json:"tags,omitempty"`
	Name                      string   `json:"name,omitempty"`
	Description               string   `json:"description,omitempty"`
	TargetStart               string   `json:"target_start,omitempty"`
	TargetEnd                 string   `json:"target_end,omitempty"`
	Status                    string   `json:"status,omitempty"`
	EngagementType            string   `json:"engagement_type,omitempty"`
	Branch                    string   `json:"branch_tag,omitempty"`
	DeduplicationOnEngagement bool     `json:"deduplication_on_engagement,omitempty"`
	Lead                      int      `json:"lead,omitempty"`
	Product                   int      `json:"product,omitempty"`
}

type Finding struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Severity string `json:"severity"`
}

type GetFindingsResponse struct {
	Count   int       `json:"count"`
	Results []Finding `json:"results"`
}

type ScanPayload struct {
	Timestamp         string   `json:"scan_date" form:"scan_date"`
	SeverityThreshold string   `json:"minimum_severity" form:"minimum_severity"`
	File              []byte   `json:"file" form:"file"`
	Branch            string   `json:"branch_tag" form:"branch_tag"`
	Tags              []string `json:"tags" form:"tags"`
	GroupBy           string   `json:"group_by" form:"group_by"`
	FindingGroup      bool     `json:"create_finding_groups_for_all_findings" form:"create_finding_groups_for_all_findings"`
	FindingTag        bool     `json:"apply_tags_to_findings" form:"apply_tags_to_findings"`
	ScanType          string   `json:"scan_type" form:"scan_type"`
	EngagementId      int      `json:"engagement" form:"engagement"`
	CloseOldFinding   bool     `json:"close_old_findings_product_scope" form:"close_old_findings_product_scope"`
}
