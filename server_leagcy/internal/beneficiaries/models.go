package beneficiaries

type Beneficiaries struct {
	Beneficiaries []Beneficiary `json:"beneficiaries"`
}

type Beneficiary struct {
	BeneficiaryId   int    `json:"beneficiary_id"`
	BeneficiaryName string `json:"beneficiary_name"`
	MobileNumber    string `json:"mobile_number"`
	Currency        string `json:"currency"`
}
