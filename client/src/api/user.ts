import axios from "axios";
import type {
  User,
  LOGIN_REQUEST,
  SIGNUP_REQUEST,
  UPDATE_USER_REQUEST,
  CHANGE_PASSWORD_REQUEST,
  SEND_PASSWORD_RESET_EMAIL_REQUEST,
  PASSWORD_RESET_REQUEST,
  LOGIN_RESPONSE,
  CONFIGURE_MFA_REQUEST,
  VERIFY_MFA_REQUEST,
} from "@/types/user";
import type { GENERIC_STATUS_RESPONSE } from "@/types/generic";

const API_URL = import.meta.env.VITE_APP_API_URL;

const LOGIN_URL = `${API_URL}/login`;
const LOGOUT_URL = `${API_URL}/logout`;
const SIGNUP_URL = `${API_URL}/signup`;
const CHANGE_PASSWORD_URL = `${API_URL}/change-password`;
const PASSWORD_RESET_URL = `${API_URL}/password-reset/reset`;
const PASSWORD_RESET_SEND_EMAIL_URL = `${API_URL}/password-reset/send`;
const CONFIGURE_MFA_URL = `${API_URL}/configure-mfa`;
const VERIFY_MFA_URL = `${API_URL}/verify-mfa`;

const ME_URL = `${API_URL}/users/me`;
const UPDATE_USER_URL = `${API_URL}/users/profile`;

export const LOGIN = async (body: LOGIN_REQUEST): Promise<LOGIN_RESPONSE> => {
  try {
    const { data, status } = await axios.post<User>(
      LOGIN_URL,
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

export const LOGOUT = async (): Promise<GENERIC_STATUS_RESPONSE> => {
  try {
    const { status } = await axios.post(LOGOUT_URL, `{}`, {
      withCredentials: true,
    });

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};

export const SIGNUP = async (
  body: SIGNUP_REQUEST,
): Promise<GENERIC_STATUS_RESPONSE> => {
  try {
    const { status } = await axios.post(SIGNUP_URL, JSON.stringify(body), {
      withCredentials: true,
    });

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};

export const CHANGE_PASSWORD = async (
  body: CHANGE_PASSWORD_REQUEST,
): Promise<GENERIC_STATUS_RESPONSE> => {
  try {
    const { status } = await axios.patch(
      CHANGE_PASSWORD_URL,
      JSON.stringify(body),
      {
        withCredentials: true,
      },
    );

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};

export const PASSWORD_RESET = async (
  body: PASSWORD_RESET_REQUEST,
): Promise<GENERIC_STATUS_RESPONSE> => {
  try {
    const { status } = await axios.patch(
      PASSWORD_RESET_URL,
      JSON.stringify(body),
      undefined,
    );

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};

export const SEND_PASSWORD_RESET_EMAIL = async (
  body: SEND_PASSWORD_RESET_EMAIL_REQUEST,
): Promise<GENERIC_STATUS_RESPONSE> => {
  try {
    const { status } = await axios.post(
      PASSWORD_RESET_SEND_EMAIL_URL,
      JSON.stringify(body),
      undefined,
    );

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};

export const GET_USER = async () => {
  try {
    const { status } = await axios.get(ME_URL, { withCredentials: true });

    return status;
  } catch (error: unknown) {
    throw error;
  }
};

export const UPDATE_USER = async (
  body: UPDATE_USER_REQUEST,
): Promise<GENERIC_STATUS_RESPONSE> => {
  try {
    const { status } = await axios.put(UPDATE_USER_URL, JSON.stringify(body), {
      withCredentials: true,
    });

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};

export const CONFIGURE_MFA = async (
  body: CONFIGURE_MFA_REQUEST,
): Promise<GENERIC_STATUS_RESPONSE> => {
  try {
    const { status } = await axios.post(
      CONFIGURE_MFA_URL,
      JSON.stringify(body),
      { withCredentials: true },
    );

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};

export const VERIFY_MFA = async (
  body: VERIFY_MFA_REQUEST,
): Promise<GENERIC_STATUS_RESPONSE> => {
  try {
    const { status } = await axios.post(VERIFY_MFA_URL, JSON.stringify(body), {
      withCredentials: true,
    });

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};
