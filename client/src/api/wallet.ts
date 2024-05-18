import axios from "axios";
import type { APIResponse } from "@/types/generic";
import type { GetWalletsResponse } from "@/types/wallet";

const API_URL = import.meta.env.VITE_APP_API_URL;

const GET_WALLETS_URL = `${API_URL}/wallet/all`;

export const GET_WALLETS = async (): Promise<
  APIResponse<GetWalletsResponse>
> => {
  try {
    const { data, status } = await axios.get<GetWalletsResponse>(
      GET_WALLETS_URL,
      {
        withCredentials: true,
      },
    );

    return { data, status };
  } catch (error: unknown) {
    throw error;
  }
};
