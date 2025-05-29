export interface User {
  userID: number;
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

export interface MoodPredictionResponse {
  prediction: string;
}
