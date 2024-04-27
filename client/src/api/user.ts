import axios from "axios";
import type { User } from "@/types/user";
import type {
  LOGIN_BODY,
  SIGNUP_BODY,
  SIGNUP_RESPONSE,
  LOGOUT_RESPONSE,
} from "@/types/user";

const BASE_URL = import.meta.env.VITE_APP_API_URL;

export const LOGIN = async (body: LOGIN_BODY) => {
  try {
    const apiURL = `${BASE_URL}/login`;
    const { data, status } = await axios.post<User>(
      apiURL,
      JSON.stringify(body),
      {
        withCredentials: true,
      },
    );

    // Return an object containing both data and status
    return { data, status };
  } catch (error: unknown) {
    throw error;
  }
};

export const LOGOUT = async (): Promise<LOGOUT_RESPONSE> => {
  try {
    const apiURL = `${BASE_URL}/logout`;
    const { status } = await axios.post(
      apiURL,
      {},
      {
        withCredentials: true,
      },
    );

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};

export const SIGNUP = async (body: SIGNUP_BODY): Promise<SIGNUP_RESPONSE> => {
  try {
    const apiURL = `${BASE_URL}/signup`;
    const { status } = await axios.post(apiURL, JSON.stringify(body), {
      withCredentials: true,
    });

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};
