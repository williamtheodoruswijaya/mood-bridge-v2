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
  // ini tak pake juga buat retrieve user data
  code: number;
  data: User;
  message: string;
}

export interface MoodPredictionResponse {
  prediction: string;
}

export interface PostResponse {
  code: number;
  data: PostInterface[];
  message: string;
}

export interface PostResponseDetail {
  code: number;
  data: PostInterface;
  message: string;
}

export interface PostInterface {
  postid: number;
  userid: number;
  user: {
    userid: number;
    username: string;
    fullname: string;
  };
  content: string;
  mood: string;
  createdat: string;
}

export interface CommentResponse {
  code: number;
  data: CommentInterface[];
  message: string;
}

export interface CommentDetailResponse {
  code: number;
  data: CommentInterface;
  message: string;
}

export interface CommentInterface {
  commentid: number;
  postid: number;
  userid: number;
  user: {
    userid: number;
    username: string;
    fullname: string;
  };
  content: string;
  created_at: string;
}

export interface FriendResponse {
  code: number;
  data: FriendInterface[];
  message: string;
}

export interface FriendInterface {
  id: number;
  userid: number;
  frienduserid: number;
  friendstatus: boolean;
  createdat: string;
  user: {
    userid: number;
    username: string;
    fullname: string;
  };
}
