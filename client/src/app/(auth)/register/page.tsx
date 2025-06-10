"use client";
import Image from "next/image";
import RegisterBackground from "~/assets/register-background.png";
import Icon from "~/assets/icon.png";
import { useState } from "react";
import { useRouter } from "next/navigation";
import axios from "axios";
import type { RegisterResponse } from "~/types/types";

export default function RegisterPage() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await axios.post<RegisterResponse>(
        `${process.env.NEXT_PUBLIC_API_URL}/api/user/register`,
        {
          username: username,
          fullname: fullName,
          email: email,
          password: password,
        },
      );
      if (response.status === 200) {
        // TODO: part ini nanti diganti sama toast
        console.log(response);
        console.log("Registration successful:", response.data.message);
        alert("Registration successful: " + response.data.message);
        router.push("/login");
      } else {
        // TODO: part ini nanti diganti sama toast
        console.error("Registration failed:", response.data.message);
        alert("Registration failed: " + response.data.message);
      }
    } catch (error) {
      // part ini nanti diganti sama toast
      console.error("Registration error:", error);
      alert("Registration error: " + (error as Error).message);
    }
  };
  return (
    <div className="flex min-h-screen">
      <div className="relative flex w-full items-center justify-center bg-[#28b7be] md:w-2/5">
        <div className="absolute top-6 left-6 flex items-center gap-2">
          <Image src={Icon} alt="Logo" className="h-15 w-15" />
          <h1 className="text-4xl font-bold text-white">Mood Bridge</h1>
        </div>
        <div className="w-full max-w-md rounded-xl bg-white p-8 shadow-lg">
          <p className="mb-6 text-center text-sm text-gray-600">
            Already have an account?{" "}
            <a href="/login" className="font-semibold text-[#28b7be]">
              Sign in
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
              htmlFor="fullname"
              className="ck text-md font-medium text-gray-700"
            >
              Fullname
            </label>
            <input
              type="text"
              placeholder="Fullname"
              className="mb-4 w-full rounded border p-3 focus:outline-none"
              value={fullName}
              onChange={(e) => setFullName(e.target.value)}
            />
            <label
              htmlFor="email"
              className="ck text-md font-medium text-gray-700"
            >
              Email
            </label>
            <input
              type="email"
              placeholder="Email"
              className="mb-4 w-full rounded border p-3 focus:outline-none"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
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
              Register
            </button>
          </form>
        </div>
      </div>
      <div className="hidden w-3/5 md:block">
        <Image
          src={RegisterBackground}
          alt="Background"
          className="h-screen w-full object-cover"
        />
      </div>
    </div>
  );
}
