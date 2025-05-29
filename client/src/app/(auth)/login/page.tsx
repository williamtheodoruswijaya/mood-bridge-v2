"use client";
import LoginBackground from "~/assets/login-background.png";
import Icon from "~/assets/icon.png";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useState } from "react";
import axios from "axios";
import Cookies from "js-cookie";
import type { LoginResponse } from "~/types/types";

export default function LoginPage() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await axios.post<LoginResponse>(
        "http://localhost:8080/api/user/login",
        {
          username: username,
          password: password,
        },
      );
      if (response.status === 200) {
        const token = response.data.data;
        Cookies.set("token", token, { expires: 7 });
        router.push("/");
      } else {
        // part ini nanti diganti sama toast
        console.error("Login failed:", response.data.message);
        alert("Login failed: " + response.data.message);
      }
    } catch (error) {
      // part ini nanti diganti sama toast
      console.error("Login error:", error);
      alert("Login error: " + (error as Error).message);
    }
  };
  return (
    <div className="flex min-h-screen">
      <div className="hidden w-1/2 md:block">
        <Image
          src={LoginBackground}
          alt="Login Background"
          className="h-screen w-full object-cover"
        />
      </div>
      <div className="relative flex w-full items-center justify-center bg-[#28b7be] md:w-1/2">
        <div className="absolute top-6 left-6 flex items-center gap-2">
          <Image src={Icon} alt="Logo" className="h-15 w-15" />
          <h1 className="text-4xl font-bold text-white">Mood Bridge</h1>
        </div>
        <div className="w-full max-w-md rounded-xl bg-white p-8 shadow-lg">
          <p className="mb-6 text-center text-sm text-gray-600">
            Don&apos;t have an account?{" "}
            <a href="/register" className="font-semibold text-[#28b7be]">
              Sign up
            </a>
          </p>
          <form>
            <label
              htmlFor="username"
              className="ck text-md font-medium text-gray-700"
            >
              Username
            </label>
            <input
              type="text"
              placeholder="Username"
              className="mb-4 w-full rounded border p-3 focus:outline-none"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
            <label
              htmlFor="password"
              className="ck text-md font-medium text-gray-700"
            >
              Password
            </label>
            <input
              type="password"
              placeholder="Password"
              className="mb-8 w-full rounded border p-3 focus:outline-none"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <button
              className="w-full rounded bg-[#28b7be] py-3 font-semibold text-white"
              type="submit"
              onClick={onSubmit}
            >
              Login
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
