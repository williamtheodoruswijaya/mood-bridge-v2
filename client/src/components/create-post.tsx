"use client";

import axios from "axios";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import type { MoodPredictionResponse, Post, User } from "~/types/types";

export default function CreatePost() {
  const router = useRouter();
  const [content, setContent] = useState("");
  const [category, setCategory] = useState("Normal");
  const [debouncedContent, setDebouncedContent] = useState("");
  const [loading, setLoading] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [user, setUser] = useState<User>({
    userID: 0,
    username: "",
    email: "",
    fullname: "",
    createdAt: "",
  });

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedContent(content);
    }, 500);
    return () => clearTimeout(handler);
  }, [content]);

  useEffect(() => {
    if (!debouncedContent) return;

    const getCategory = async () => {
      try {
        const response = await axios.post<MoodPredictionResponse>(
          "https://adamantix-ensemble-model-mental-illness-classification.hf.space/mic-predict",
          { input: debouncedContent },
        );
        if (response.status === 200) {
          setCategory(response.data.prediction);
        }
      } catch (error) {
        console.error("Error fetching mood prediction:", error);
      }
    };
    getCategory().catch((error) => {
      console.error("Error in useEffect:", error);
      setCategory("Normal");
    });
  }, [category, debouncedContent]);

  useEffect(() => {
    const token = Cookies.get("token");
    if (token) {
      try {
        const base64Payload = token.split(".")[1];
        if (!base64Payload) throw new Error("Invalid token structure");
        const decodedPayload = atob(base64Payload);
        const parsed = JSON.parse(decodedPayload) as {
          user: {
            id: number;
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
            userID: user.id,
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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isLoggedIn) router.push("/login");
    if (!content.trim()) return;
    setLoading(true);
    try {
      console.log("user", user);
      const response = await axios.post<Post>(
        "http://localhost:8080/api/post/create",
        {
          userID: user.userID,
          content: content,
        },
        {
          headers: {
            Authorization: `Bearer ${Cookies.get("token")}`,
          },
        },
      );
      if (response.status === 200) {
        // TODO: Ganti sama toast
        console.log("Post created successfully:", response.data.message);
        alert("Post created successfully: " + response.data.message);

        // reset content and category
        setContent("");
        setCategory("Normal");
      }
    } catch (error) {
      console.error("Error creating post:", error);
      // Ganti sama toast
      alert("Error creating post: " + (error as Error).message);
    }
  };

  return (
    <div className="w-full rounded-xl bg-[#84E7EE] p-4 shadow-lg backdrop-blur-md">
      <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
        <div className="relative w-full">
          <div className="absolute top-1 right-2 z-10 min-w-20 rounded-md bg-white px-2 py-1 text-center text-sm font-semibold text-black shadow">
            {category}
          </div>
          <textarea
            className="focus-outline-none w-full resize-none rounded-md bg-[#73CFD5] p-5 text-black"
            rows={3}
            placeholder="What's happening..."
            value={content}
            onChange={(e) => setContent(e.target.value)}
          />
        </div>
        <button
          type="submit"
          className="min-w-20 self-end rounded-md bg-[#4DC0D9] px-4 py-1 font-semibold text-white hover:bg-[#3AAFC9] disabled:opacity-50"
          disabled={loading}
        >
          {loading ? "Posting..." : "Post"}
        </button>
      </form>
    </div>
  );
}
