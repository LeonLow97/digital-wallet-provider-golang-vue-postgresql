import axios from "axios";
import type { HttpStatusResponse, ApiResponse } from "@/types/generic";
import type {
  CreateBeneficiaryRequest,
  DeleteBeneficiaryRequest,
  GetBeneficiariesResponse,
} from "@/types/beneficiary";

const API_URL = import.meta.env.VITE_APP_API_URL;

const CREATE_BENEFICIARY_URL = `${API_URL}/beneficiary`;
const GET_BENEFICIARIES_URL = `${API_URL}/beneficiary`;
const DELETE_BENEFICIARY_URL = `${API_URL}/beneficiary`;

export const CREATE_BENEFICIARY = async (body: CreateBeneficiaryRequest) => {
  const { status } = await axios.post<HttpStatusResponse>(
    CREATE_BENEFICIARY_URL,
    body,
    { withCredentials: true },
  );
  return { status };
};

export const GET_BENEFICIARIES = async (): Promise<
  ApiResponse<GetBeneficiariesResponse>
> => {
  const { data, status } = await axios.get<GetBeneficiariesResponse>(
    GET_BENEFICIARIES_URL,
    { withCredentials: true },
  );

  return { data, status };
};

export const DELETE_BENEFICIARY = async (
  body: DeleteBeneficiaryRequest,
): Promise<HttpStatusResponse> => {
  const { status } = await axios.put<HttpStatusResponse>(
    DELETE_BENEFICIARY_URL,
    body,
    { withCredentials: true },
  );

  return { status };
};
