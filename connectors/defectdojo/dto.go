package defectdojo

type Product struct {
	Id uint `json:"id"`
}

type GetProductByNameResponse struct {
	Count   uint      `json:"count"`
	Results []Product `json:"results"`
}
