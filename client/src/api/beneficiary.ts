import axios from "axios";
import type { GENERIC_STATUS_RESPONSE, APIResponse } from "@/types/generic";
import type {
  CREATE_BENEFICIARY_REQUEST,
  DELETE_BENEFICIARY_REQUEST,
  GET_BENEFICIARIES_RESPONSE,
} from "@/types/beneficiary";

const API_URL = import.meta.env.VITE_APP_API_URL;

const CREATE_BENEFICIARY_URL = `${API_URL}/beneficiary`;
const GET_BENEFICIARIES_URL = `${API_URL}/beneficiary`;
const DELETE_BENEFICIARY_URL = `${API_URL}/beneficiary`;

export const CREATE_BENEFICIARY = async (body: CREATE_BENEFICIARY_REQUEST) => {
    const { status } = await axios.post<GENERIC_STATUS_RESPONSE>(
      CREATE_BENEFICIARY_URL,
      body,
      { withCredentials: true },
    );
    return { status };
};

export const GET_BENEFICIARIES = async (): Promise<
  APIResponse<GET_BENEFICIARIES_RESPONSE>
> => {
    const { data, status } = await axios.get<GET_BENEFICIARIES_RESPONSE>(
      GET_BENEFICIARIES_URL,
      { withCredentials: true },
    );

    return { data, status };
};

export const DELETE_BENEFICIARY = async (
  body: DELETE_BENEFICIARY_REQUEST,
): Promise<GENERIC_STATUS_RESPONSE> => {
    const { status } = await axios.put<GENERIC_STATUS_RESPONSE>(
      DELETE_BENEFICIARY_URL,
      body,
      { withCredentials: true },
    );

    return { status };
};
