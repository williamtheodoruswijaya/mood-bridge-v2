export interface User {
  username: string;
  fullname: string;
  email: string;
  createdAt: string;
}

export interface LoginResponse {
  code: number;
  data: string; // JWT token
  message: string;
}

export interface RegisterResponse {
  code: number;
  data: User;
  message: string;
}
