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

export function TimeAgo(isoString: string): string {
  const now = new Date();
  const then = new Date(isoString);
  const diffInSeconds = Math.floor((now.getTime() - then.getTime()) / 1000);

  if (diffInSeconds < 60) {
    return `${diffInSeconds}s ago`;
  } else if (diffInSeconds < 3600) {
    const minutes = Math.floor(diffInSeconds / 60);
    return `${minutes}m ago`;
  } else if (diffInSeconds < 86400) {
    const hours = Math.floor(diffInSeconds / 3600);
    return `${hours}h ago`;
  } else {
    const days = Math.floor(diffInSeconds / 86400);
    return `${days}d ago`;
  }
}
