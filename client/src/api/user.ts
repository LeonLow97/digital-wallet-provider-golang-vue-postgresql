import axios from "axios";
import type {
  User,
  LoginRequest,
  SignUpRequest,
  UpdateUserRequest,
  ChangePasswordRequest,
  SendPasswordResetEmailRequest,
  PasswordResetRequest,
  LoginResponse,
  ConfigureMfaRequest,
  VerifyMfaRequest,
} from "@/types/user";
import type { HttpStatusResponse } from "@/types/generic";

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

export const LOGIN = async (body: LoginRequest): Promise<LoginResponse> => {
  const { data, status } = await axios.post<User>(
    LOGIN_URL,
    JSON.stringify(body),
    {
      withCredentials: true,
    },
  );

  return { data, status };
};

export const LOGOUT = async (): Promise<HttpStatusResponse> => {
  const { status } = await axios.post(LOGOUT_URL, `{}`, {
    withCredentials: true,
  });

  return { status };
};

export const SIGNUP = async (
  body: SignUpRequest,
): Promise<HttpStatusResponse> => {
  const { status } = await axios.post(SIGNUP_URL, JSON.stringify(body), {
    withCredentials: true,
  });

  return { status };
};

export const CHANGE_PASSWORD = async (
  body: ChangePasswordRequest,
): Promise<HttpStatusResponse> => {
  const { status } = await axios.patch(
    CHANGE_PASSWORD_URL,
    JSON.stringify(body),
    {
      withCredentials: true,
    },
  );

  return { status };
};

export const PASSWORD_RESET = async (
  body: PasswordResetRequest,
): Promise<HttpStatusResponse> => {
  const { status } = await axios.patch(
    PASSWORD_RESET_URL,
    JSON.stringify(body),
    undefined,
  );

  return { status };
};

export const SEND_PASSWORD_RESET_EMAIL = async (
  body: SendPasswordResetEmailRequest,
): Promise<HttpStatusResponse> => {
  const { status } = await axios.post(
    PASSWORD_RESET_SEND_EMAIL_URL,
    JSON.stringify(body),
    undefined,
  );

  return { status };
};

export const GET_USER = async () => {
  const { status } = await axios.get(ME_URL, { withCredentials: true });

  return status;
};

export const UPDATE_USER = async (
  body: UpdateUserRequest,
): Promise<HttpStatusResponse> => {
  const { status } = await axios.put(UPDATE_USER_URL, JSON.stringify(body), {
    withCredentials: true,
  });

  return { status };
};

export const CONFIGURE_MFA = async (
  body: ConfigureMfaRequest,
): Promise<HttpStatusResponse> => {
  const { status } = await axios.post(CONFIGURE_MFA_URL, JSON.stringify(body), {
    withCredentials: true,
  });

  return { status };
};

export const VERIFY_MFA = async (
  body: VerifyMfaRequest,
): Promise<HttpStatusResponse> => {
  const { status } = await axios.post(VERIFY_MFA_URL, JSON.stringify(body), {
    withCredentials: true,
  });

  return { status };
};
