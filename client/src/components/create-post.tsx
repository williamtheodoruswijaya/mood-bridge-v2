"use client";

import axios from "axios";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import type { MoodPredictionResponse, PostResponse, User } from "~/types/types";

export default function CreatePost() {
  const router = useRouter();
  const [content, setContent] = useState("");
  const [category, setCategory] = useState("Your mood will appear here...");
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
  const categoryColor: Record<string, string> = {
    Normal: "#219E2C",
    Anxiety: "#FFAE00",
    Depression: "#0D00FF",
    Suicidal: "#FF0000",
    Stress: "#FF00A0",
    Bipolar: "#8B00FF",
    "Personality Disorder": "#000000",
  };

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedContent(content);
    }, 500);
    return () => clearTimeout(handler);
  }, [content]);

  useEffect(() => {
    if (debouncedContent.length <= 0) {
      setCategory("Write something here...");
      return;
    }
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
      const response = await axios.post<PostResponse>(
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
      }
    } catch (error) {
      console.error("Error creating post:", error);
      // TODO: Ganti sama toast
      alert("Error creating post: " + (error as Error).message);
    } finally {
      setLoading(false);
      setContent("");
      setCategory("Your mood will appear here...");
      setDebouncedContent("");
      router.refresh();
    }
  };

  return (
    <div className="w-full rounded-xl bg-[#84E7EE] p-4 shadow-lg backdrop-blur-md">
      <form className="flex flex-col gap-1" onSubmit={handleSubmit}>
        <div className="relative flex w-full flex-col">
          <div
            className="text-md absolute top-2 right-2 z-10 rounded-md px-2 py-1 text-center font-semibold text-white shadow"
            style={{ backgroundColor: categoryColor[category] ?? "#687669" }}
          >
            {category}
          </div>

          <textarea
            className="w-full resize-none rounded-md bg-[#73CFD5] p-5 pr-32 pb-16 text-black focus:outline-none"
            rows={2}
            placeholder="How are you feeling today? Share your thoughts..."
            value={content}
            onChange={(e) => setContent(e.target.value)}
          />

          <button
            type="submit"
            className="text-md absolute right-2 bottom-2 min-w-24 rounded-md bg-[#30ACFF] px-4 py-1 font-bold text-white hover:bg-[#0085DE] disabled:opacity-50"
            disabled={loading}
          >
            {loading ? "Posting..." : "Post"}
          </button>
        </div>
      </form>
    </div>
  );
}
