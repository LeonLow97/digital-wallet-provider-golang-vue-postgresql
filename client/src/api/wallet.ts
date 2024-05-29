import axios from "axios";
import type { APIResponse, GENERIC_STATUS_RESPONSE } from "@/types/generic";
import type {
  Wallet,
  CreateWalletRequest,
  GetWalletTypesResponse,
  TopUpWalletRequest,
} from "@/types/wallet";

const API_URL = import.meta.env.VITE_APP_API_URL;

const GET_WALLETS_URL = `${API_URL}/wallet/all`;
const GET_WALLET_URL = `${API_URL}/wallet/`;
const GET_WALLET_TYPES_URL = `${API_URL}/wallet/types`;
const CREATE_WALLET_URL = `${API_URL}/wallet`;
const TOP_UP_WALLET_URL = `${API_URL}/wallet/topup/`;

export const GET_WALLETS = async (): Promise<APIResponse<Wallet[]>> => {
  const { data, status } = await axios.get<Wallet[]>(GET_WALLETS_URL, {
    withCredentials: true,
  });

  return { data, status };
};

export const GET_WALLET = async (
  paramId: number,
): Promise<APIResponse<Wallet>> => {
  const { data, status } = await axios.get<Wallet>(GET_WALLET_URL + paramId, {
    withCredentials: true,
  });

  return { data, status };
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

export const TOP_UP_WALLET = async (
  walletId: number,
  body: TopUpWalletRequest,
): Promise<GENERIC_STATUS_RESPONSE> => {
  const { status } = await axios.put<GENERIC_STATUS_RESPONSE>(
    TOP_UP_WALLET_URL + walletId,
    body,
    { withCredentials: true },
  );

  return { status };
};
