package defectdojo

const (
	APIPrefix       = "/api/v2/"
	GetProductsPath = "products?name_exact="
)

const (
	logErrorRetrieveProducts = "Cannot retrieve defectdojo products"
	logErrorAuthorization    = "Wrong API Key"
	logErrorDecodingToken    = "Cannot unmarshall json [%s]"
)

const (
	errRetrieveProducts = "cannot retrieve defectdojo products"
	errAuthtorization   = "wrong api key"
	errDuplicateProduct = "two products with the same name already exist"
	errProductNotExist  = "product does not exist"
	errUnmarshal        = "cannot unmarshal json"
)
