import axios from "axios";
import type { APIResponse, GENERIC_STATUS_RESPONSE } from "@/types/generic";
import type {
  CreateWalletRequest,
  GetWalletTypesResponse,
  GetWalletsResponse,
} from "@/types/wallet";

const API_URL = import.meta.env.VITE_APP_API_URL;

const GET_WALLETS_URL = `${API_URL}/wallet/all`;
const GET_WALLET_TYPES_URL = `${API_URL}/wallet/types`;
const CREATE_WALLET_URL = `${API_URL}/wallet`;

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

export const GET_WALLET_TYPES = async (): Promise<
  APIResponse<GetWalletTypesResponse[]>
> => {
  const { data, status } = await axios.get<GetWalletTypesResponse[]>(
    GET_WALLET_TYPES_URL,
    { withCredentials: true },
  );

  return { data, status };
};

export const CREATE_WALLET = async (
  body: CreateWalletRequest,
): Promise<GENERIC_STATUS_RESPONSE> => {
  const { status } = await axios.post<GENERIC_STATUS_RESPONSE>(
    CREATE_WALLET_URL,
    body,
    { withCredentials: true },
  );

  return { status };
};
