export interface User {
  firstName: string;
  lastName: string;
  email: string;
  username: string;
  mobileNumber: string;
}

export interface LOGIN_REQUEST {
  email: string;
  password: string;
}

export interface LOGIN_RESPONSE {
  data: LOGIN_RESPONSE_DATA;
  status: number;
}

interface LOGIN_RESPONSE_DATA {
  firstName: string;
  lastName: string;
  email: string;
  username: string;
  mobileNumber: string;
  isMfaConfigured?: boolean;
  mfaConfig?: MFA_CONFIG;
}

interface MFA_CONFIG {
  secret: string;
  url: string;
}

export interface SIGNUP_REQUEST {
  first_name: string | null;
  last_name: string | null;
  username: string;
  email: string;
  password: string;
  mobile_number: string;
}

export interface CHANGE_PASSWORD_REQUEST {
  current_password: string;
  new_password: string;
}

export interface SEND_PASSWORD_RESET_EMAIL_REQUEST {
  email: string;
}

export interface PASSWORD_RESET_REQUEST {
  token: string;
  password: string;
}

export interface UPDATE_USER_REQUEST {
  first_name: string | null;
  last_name: string | null;
  username: string;
  email: string;
  mobile_number: string;
}

export interface CONFIGURE_MFA_REQUEST {
  email: string;
  secret: string;
  mfa_code: string;
}

export interface VERIFY_MFA_REQUEST {
  email: string;
  mfa_code: string;
}