export interface User {
  firstName: string;
  lastName: string;
  email: string;
  username: string;
  mobileNumber: string;
}

export interface LOGIN_BODY {
  email: string;
  password: string;
}

export interface SIGNUP_BODY {
  first_name: string | null;
  last_name: string | null;
  username: string;
  email: string;
  password: string;
  mobile_number: string;
}

export interface SIGNUP_RESPONSE {
  status: number;
}

export interface LOGOUT_RESPONSE {
  status: number;
}
