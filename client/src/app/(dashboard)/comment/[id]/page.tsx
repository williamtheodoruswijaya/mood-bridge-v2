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
          <div className="mt-6">
            <label
              htmlFor="comment"
              className="block text-sm font-medium text-gray-700"
            >
              Write a comment
            </label>
            <textarea
              id="comment"
              name="comment"
              rows={4}
              className="focus:ring-opacity-50 mt-1 block w-full rounded-lg border border-gray-300 p-3 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
              placeholder="Share your thoughts..."
            />
            <button
              type="button"
              className="mt-2 rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
            >
              Post Comment
            </button>
          </div>
        </div>
      </section>
    </main>
  );
}
