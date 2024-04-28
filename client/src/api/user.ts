import axios from "axios";
import type { User } from "@/types/user";
import type {
  LOGIN_REQUEST,
  SIGNUP_REQUEST,
  SIGNUP_RESPONSE,
  LOGOUT_RESPONSE,
  UPDATE_USER_REQUEST,
  UPDATE_USER_RESPONSE,
} from "@/types/user";

const API_URL = import.meta.env.VITE_APP_API_URL;

export const LOGIN = async (body: LOGIN_REQUEST) => {
  try {
    const apiURL = `${API_URL}/login`;
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
    const apiURL = `${API_URL}/logout`;
    const { status } = await axios.post(apiURL, `{}`, {
      withCredentials: true,
    });

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};

export const SIGNUP = async (
  body: SIGNUP_REQUEST,
): Promise<SIGNUP_RESPONSE> => {
  try {
    const apiURL = `${API_URL}/signup`;
    const { status } = await axios.post(apiURL, JSON.stringify(body), {
      withCredentials: true,
    });

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};

export const UPDATE_USER = async (
  body: UPDATE_USER_REQUEST,
): Promise<UPDATE_USER_RESPONSE> => {
  try {
    const apiURL = `${API_URL}/users/profile`;
    const { status } = await axios.post(apiURL, JSON.stringify(body), {
      withCredentials: true,
    });

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};
