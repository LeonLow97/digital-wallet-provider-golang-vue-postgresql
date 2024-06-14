export interface CreateBeneficiaryRequest {
  mobile_country_code: string;
  mobile_number: string;
}

export interface DeleteBeneficiaryRequest {
  is_deleted: number;
  beneficiary_id: number;
}

// Interface for a single beneficiary response
export interface GetBeneficiaryResponse {
  beneficiaryID: number;
  isDeleted: number;
  beneficiaryFirstName: string;
  beneficiaryLastName: string;
  beneficiaryEmail: string;
  beneficiaryUsername: string;
  active: number;
  beneficiaryMobileCountryCode: string;
  beneficiaryMobileNumber: string;
}

// Interface for multiple beneficiaries response
export interface GetBeneficiariesResponse {
  beneficiaries: GetBeneficiaryResponse[];
}
