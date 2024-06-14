import axios from "axios";
import type { ApiResponse, HttpStatusResponse } from "@/types/generic";
import type {
  Wallet,
  CreateWalletRequest,
  GetWalletTypesResponse,
  WalletExchangesRequest,
} from "@/types/wallet";

const API_URL = import.meta.env.VITE_APP_API_URL;

const GET_WALLETS_URL = `${API_URL}/wallet/all`;
const GET_WALLET_URL = `${API_URL}/wallet/`;
const GET_WALLET_TYPES_URL = `${API_URL}/wallet/types`;
const CREATE_WALLET_URL = `${API_URL}/wallet`;
const UPDATE_WALLET_URL = `${API_URL}/wallet/update/`;

export const GET_WALLETS = async (): Promise<ApiResponse<Wallet[]>> => {
  const { data, status } = await axios.get<Wallet[]>(GET_WALLETS_URL, {
    withCredentials: true,
  });

  return { data, status };
};

export const GET_WALLET = async (
  paramId: number,
): Promise<ApiResponse<Wallet>> => {
  const { data, status } = await axios.get<Wallet>(GET_WALLET_URL + paramId, {
    withCredentials: true,
  });

  return { data, status };
};

export const GET_WALLET_TYPES = async (): Promise<
  ApiResponse<GetWalletTypesResponse[]>
> => {
  const { data, status } = await axios.get<GetWalletTypesResponse[]>(
    GET_WALLET_TYPES_URL,
    { withCredentials: true },
  );

  return { data, status };
};

export const CREATE_WALLET = async (
  body: CreateWalletRequest,
): Promise<HttpStatusResponse> => {
  const { status } = await axios.post<HttpStatusResponse>(
    CREATE_WALLET_URL,
    body,
    { withCredentials: true },
  );

  return { status };
};

export const UPDATE_WALLET = async (
  walletId: number,
  operation: string,
  body: WalletExchangesRequest,
): Promise<HttpStatusResponse> => {
  const { status } = await axios.put<HttpStatusResponse>(
    `${UPDATE_WALLET_URL}/${walletId}/${operation}`,
    body,
    { withCredentials: true },
  );

  return { status };
};
