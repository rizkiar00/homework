package model

type CreateTestDBRequest struct {
	DescTest *string `json:"desc_test"`
}

type UpdateTestDBRequest struct {
	DescTest *string `json:"desc_test"`
}

type TestDBResponse struct {
	IDTest   string  `json:"id_test"`
	DescTest *string `json:"desc_test"`
}
