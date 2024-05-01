import axios from "axios";
import type { Wallets } from "@/types/wallet";

const API_URL = import.meta.env.VITE_APP_API_URL;

export const GET_WALLETS = async () => {
  try {
    const apiURL = `${API_URL}/wallet/all`;
    const { data, status } = await axios.get<Wallets>(apiURL, {
      withCredentials: true,
    });

    return { data, status };
  } catch (error: unknown) {
    throw error;
  }
};
