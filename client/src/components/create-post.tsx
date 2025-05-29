"use client";

import axios from "axios";
import { useEffect, useState } from "react";
import type { MoodPredictionResponse } from "~/types/types";

export default function CreatePost() {
  const [content, setContent] = useState("");
  const [category, setCategory] = useState("Normal");
  const [debouncedContent, setDebouncedContent] = useState("");
  const [loading, setLoading] = useState(false);

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
  }, [debouncedContent]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim()) return;
    setLoading(true);
    try {
      await axios.post("http://localhost:8080/api/post/create", {
        {
          userid: 
        }
      })
    }
  }
  return (
    <div className="w-full max-w-2xl rounded-xl bg-white/20 p-4 shadow-lg backdrop-blur-md">
      <form className="flex flex-col gap-4">
        <textarea
          className="w-full resize-none rounded-md bg-white/10 p-3 text-white placeholder-white/70 focus:ring-2 focus:ring-cyan-300 focus:outline-none"
          rows={3}
          placeholder="What's happening..."
          value={content}
          onChange={(e) => setContent(e.target.value)}
        />
      </form>
    </div>
  );
}
