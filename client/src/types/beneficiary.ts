export interface CREATE_BENEFICIARY_REQUEST {
  mobile_number: string;
}

export interface DELETE_BENEFICIARY_REQUEST {
  is_deleted: number;
  beneficiary_id: number;
}

// Interface for a single beneficiary response
export interface GET_BENEFICIARY_RESPONSE {
  beneficiaryID: number;
  isDeleted: number;
  beneficiaryFirstName: string;
  beneficiaryLastName: string;
  beneficiaryEmail: string;
  beneficiaryUsername: string;
  active: number;
  beneficiaryMobileNumber: string;
}

// Interface for multiple beneficiaries response
export interface GET_BENEFICIARIES_RESPONSE {
  beneficiaries: GET_BENEFICIARY_RESPONSE[];
}
