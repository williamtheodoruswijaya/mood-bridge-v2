"use client";

import axios from "axios";
import Cookies from "js-cookie";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";
import { DecodeUserFromToken } from "~/utils/utils";
import PostDetail from "~/components/post-detail";
import { type PostInterface, type PostResponseDetail } from "~/types/types";

export default function Page() {
  const params = useParams();
  const postID = params.id;
  const [loading, setLoading] = useState(true);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [user, setUser] = useState({
    userID: 0,
    username: "",
    email: "",
    fullname: "",
    createdAt: "",
  });
  const [post, setPost] = useState<PostInterface>({
    postid: 0,
    userid: 0,
    user: {
      userid: 0,
      username: "",
      fullname: "",
    },
    content: "",
    mood: "",
    createdAt: "",
  });
  const [rows, setRows] = useState(1);
  const [value, setValue] = useState("");
  const [focused, setFocused] = useState(false);

  const handleFocus = () => {
    setRows(5);
    setFocused(true);
  };

  const handleBlur = () => {
    if (value.trim() === "") {
      setRows(1);
      setFocused(false);
    }
  };

  useEffect(() => {
    const fetchPost = async (postID: string) => {
      try {
        const response = await axios.get<PostResponseDetail>(
          `http://localhost:8080/api/post/by-id/${postID}`,
          {
            headers: {
              Authorization: `Bearer ${Cookies.get("token")}`,
              "Content-Type": "application/json",
            },
          },
        );
        if (response.status === 200) {
          setPost(response.data.data);
        } else {
          console.error("Failed to fetch post:", response.statusText);
        }
      } catch (error) {
        console.error("Error fetching post:", error);
      }
    };
    const fetchUserAndPosts = async () => {
      const token = Cookies.get("token");
      setIsLoggedIn(false);
      if (token) {
        const user = DecodeUserFromToken(token);
        if (user) {
          setUser({
            userID: user.user.id,
            username: user.user.username,
            email: user.user.email,
            fullname: user.user.fullname,
            createdAt: user.user.created_at,
          });
          setIsLoggedIn(true);
          await fetchPost(postID as string);
        }
      }
    };
    fetchUserAndPosts().catch((error) => {
      console.error("Error fetching user and posts:", error);
      setLoading(false);
    });
  }, [postID]);

  return (
    <main className="grid h-screen w-full">
      <section className="px-6">
        <div className="mx-auto mt-4 w-full">
          <PostDetail {...post} />
          <div className="relative mt-4">
            <textarea
              id="comment"
              name="comment"
              rows={rows}
              value={value}
              onChange={(e) => setValue(e.target.value)}
              onFocus={handleFocus}
              onBlur={handleBlur}
              className="block w-full resize-none rounded-lg border border-gray-300 p-3 pr-24 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
              placeholder="Share your thoughts..."
            />
            {(focused || value.trim() !== "") && (
              <button
                type="button"
                className="absolute right-3 bottom-3 rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
              >
                Post Comment
              </button>
            )}
          </div>
        </div>
      </section>
    </main>
  );
}
