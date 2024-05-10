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

export interface GENERIC_STATUS_RESPONSE {
  status: number;
}

export interface UPDATE_USER_REQUEST {
  first_name: string | null;
  last_name: string | null;
  username: string;
  email: string;
  mobile_number: string;
}
