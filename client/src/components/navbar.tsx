"use client";

import { useState, useEffect } from "react";
import Cookies from "js-cookie";
import type { User } from "~/types/types";

export default function Navbar() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [user, setUser] = useState<User>({
    username: "",
    email: "",
    fullname: "",
    createdAt: "",
  });

  useEffect(() => {
    const token = Cookies.get("token");

    if (token) {
      try {
        const base64Payload = token.split(".")[1];
        if (!base64Payload) throw new Error("Invalid token structure");
        const decodedPayload = atob(base64Payload);
        const parsed = JSON.parse(decodedPayload) as {
          user: {
            username: string;
            fullname: string;
            email: string;
            created_at: string;
          };
          exp: number;
        };
        const user = parsed.user;
        if (user) {
          setUser({
            username: user.username,
            email: user.email,
            fullname: user.fullname,
            createdAt: user.created_at,
          });
        }
        setIsLoggedIn(true);
      } catch (err) {
        console.error("Token parsing failed:", err);
        setIsLoggedIn(false);
      }
    } else {
      setIsLoggedIn(false);
    }
  }, []);
}
