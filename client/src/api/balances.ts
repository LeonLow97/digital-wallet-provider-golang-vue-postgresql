import axios from "axios";
import type { APIResponse } from "@/types/generic";
import type { GetBalancesResponse } from "@/types/balances";

const API_URL = import.meta.env.VITE_APP_API_URL;

const GET_BALANCES_URL = `${API_URL}/balances`;

export const GET_BALANCES = async (): Promise<
  APIResponse<GetBalancesResponse>
> => {
  try {
    const { data, status } = await axios.get<GetBalancesResponse>(
      GET_BALANCES_URL,
      {
        withCredentials: true,
      },
    );

    return { data, status };
  } catch (error: unknown) {
    throw error;
  }
};
