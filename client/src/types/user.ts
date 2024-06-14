export interface User {
  firstName: string;
  lastName: string;
  email: string;
  username: string;
  sourceCurrency: string;
  mobileCountryCode: string;
  mobileNumber: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  data: LoginResponseData;
  status: number;
}

interface LoginResponseData {
  firstName: string;
  lastName: string;
  email: string;
  username: string;
  sourceCurrency: string;
  mobileCountryCode: string;
  mobileNumber: string;
  isMfaConfigured?: boolean;
  mfaConfig?: MfaConfig;
}

interface MfaConfig {
  secret: string;
  url: string;
}

export interface SignUpRequest {
  first_name: string | null;
  last_name: string | null;
  username: string;
  email: string;
  password: string;
  mobile_country_code: string;
  mobile_number: string;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

export interface SendPasswordResetEmailRequest {
  email: string;
}

export interface PasswordResetRequest {
  token: string;
  password: string;
}

export interface UpdateUserRequest {
  first_name: string | null;
  last_name: string | null;
  username: string;
  email: string;
  mobile_country_code: string;
  mobile_number: string;
}

export interface ConfigureMfaRequest {
  email: string;
  secret: string;
  mfa_code: string;
}

export interface VerifyMfaRequest {
  email: string;
  mfa_code: string;
}
