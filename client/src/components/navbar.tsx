"use client";

import { useState, useEffect } from "react";
import Cookies from "js-cookie";
import type { User } from "~/types/types";
import { CiSearch } from "react-icons/ci";
import { useRouter } from "next/navigation";
import { IoMdNotificationsOutline } from "react-icons/io";
import { FaUserCircle } from "react-icons/fa";

export default function Navbar() {
  const router = useRouter();
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [user, setUser] = useState<User>({
    username: "",
    email: "",
    fullname: "",
    createdAt: "",
  });
  const [searchQuery, setSearchQuery] = useState("");
  const handleSearchKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      e.preventDefault();
      // filter post by backend (nanti belom jadi)
      setSearchQuery("");
    }
  };

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

  return (
    <nav className="flex w-full items-center justify-between bg-[#28b7be] px-6 py-3 text-white shadow">
      <div className="ml-auto flex items-center gap-4">
        {isLoggedIn && (
          <>
            <div className="relative w-64">
              <CiSearch className="absolute top-1/2 left-3 -translate-y-1/2 transform text-gray-400" />
              <input
                type="text"
                placeholder="Search..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyDown={handleSearchKeyDown}
                className="w-full rounded-md border border-gray-300 bg-white px-10 py-2 text-gray-700 focus:border-blue-500 focus:outline-none"
              />
            </div>
            <button onClick={() => router.push("/notifications")}>
              <IoMdNotificationsOutline className="text-2xl text-white hover:text-gray-200" />
            </button>
            <button onClick={() => router.push("/profile")}>
              <FaUserCircle className="text-2xl text-white hover:text-gray-200" />
            </button>
          </>
        )}
      </div>
    </nav>
  );
}
