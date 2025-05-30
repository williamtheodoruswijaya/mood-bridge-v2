"use client";

import axios from "axios";
import Cookies from "js-cookie";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";
import Post from "~/components/post";
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
  })
  useEffect(() => {
    const fetchPost = async (postID: string) => {
      try {
        const response = await axios.get<PostResponseDetail>(`http://localhost:8080/api/post/by-id/${postID}`,
          {
            headers: {
              Authorization: `Bearer ${Cookies.get("token")}`,
              "Content-Type": "application/json",
            }
          }
        )
        if (response.status === 200) {
          const fetchedPost = response.data.data;
          setPost(fetchedPost);
        }
      } catch (error) {
        console.error("Error fetching post:", error);
      } finally {
        setLoading(false);
      }
    }

    const fetchUserAndPosts = async () => {
      const token = Cookies.get("token");
      let parsedUser = null;
      setIsLoggedIn(false);
      if (token) {
        try {
          const base64Payload = token.split(".")[1];
          if (!base64Payload) throw new Error("Invalid token format");
          const decodedPayload = atob(base64Payload);
          parsedUser = JSON.parse(decodedPayload) as {
            user: {
              id: number;
              username: string;
              fullname: string;
              email: string;
              created_at: string;
            };
            exp: number;
          };
          const user = parsedUser.user;
          if (user) {
            setUser({
              userID: user.id,
              username: user.username,
              email: user.email,
              fullname: user.fullname,
              createdAt: user.created_at,
            })
            fetchPost(postID as string);
            setIsLoggedIn(true);
          }
        } catch (err) {
          console.error("Token parsing failed:", err);
          setIsLoggedIn(false);
        }
      }
    }
    fetchUserAndPosts().catch((error) => {
      console.error("Error fetching user and posts:", error);
      setLoading(false);
    })
  }, [])

  return (
    <main className="grid h-screen w-full">
      <section className=" px-6">
        <div className="mx-auto mt-4 w-full">
          <PostDetail {...post}/>
        </div>
      </section>
    </main>
  );
}
