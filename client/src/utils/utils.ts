export function DecodeUserFromToken(token: string) {
  try {
    const base64Payload = token.split(".")[1];
    if (!base64Payload) throw new Error("Invalid token format");
    const decoded = atob(base64Payload);
    return JSON.parse(decoded) as {
      user: {
        id: number;
        username: string;
        fullname: string;
        email: string;
        created_at: string;
      };
      exp: number;
    };
  } catch {
    return null;
  }
}
