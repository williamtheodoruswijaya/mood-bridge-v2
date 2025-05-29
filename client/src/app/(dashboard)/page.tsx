"use client";
import axios from "axios";
import Cookies from "js-cookie";
import { useEffect, useState } from "react";
import CreatePost from "~/components/create-post";
import Post from "~/components/post";
import type { PostInterface, PostResponse, User } from "~/types/types";

export default function HomePage() {
  const [posts, setPosts] = useState<PostInterface[]>([]);
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState<User>({
    userID: 0,
    username: "",
    email: "",
    fullname: "",
    createdAt: "",
  });
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
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
            });
            setIsLoggedIn(true);
          }
        } catch (err) {
          console.error("Token parsing failed:", err);
          setIsLoggedIn(false);
        }
      }

      setLoading(true);

      const fetchAllPosts = async () => {
        const apiUrl = "http://localhost:8080/api/post/all";
        try {
          const response = await axios.get<PostResponse>(apiUrl);
          if (response.status === 200) {
            setPosts(response.data.data);
          } else {
            console.error("Failed to fetch all posts:", response.statusText);
          }
        } catch (error) {
          console.error("Error fetching all posts:", error);
        } finally {
          setLoading(false);
        }
      };

      const fetchFriendPosts = async (userID: number) => {
        const apiUrl = `http://localhost:8080/api/post/friend-posts/${userID}`;
        try {
          const response = await axios.get<PostResponse>(apiUrl, {
            headers: {
              Authorization: `Bearer ${Cookies.get("token")}`,
              "Content-Type": "application/json",
            },
          });
          if (response.status === 200) {
            setPosts(response.data.data);
          } else {
            console.error("Failed to fetch friend posts:", response.statusText);
            await fetchAllPosts(); // fallback
          }
        } catch (error) {
          console.error("Error fetching friend posts:", error);
          await fetchAllPosts(); // fallback
        } finally {
          setLoading(false);
        }
      };

      if (parsedUser && parsedUser.user.id > 0) {
        await fetchFriendPosts(parsedUser.user.id);
      } else {
        await fetchAllPosts();
      }
    };

    fetchUserAndPosts().catch((error) => {
      console.error("Error in fetchUserAndPosts:", error);
      setLoading(false);
    });
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  function renderPosts() {
    if (loading) return <p className="text-black">Loading posts...</p>;
    if (posts.length === 0)
      return <p className="text-black">No posts available</p>;
    return posts.map((post) => <Post key={post.postID} {...post} />);
  }
  return (
    <main className="flex min-h-screen w-full text-white">
      <section className="flex flex-1 flex-col items-center px-6">
        {/* Bagian Atas */}
        <div className="mt-4 w-full max-w-4xl">
          <CreatePost
            onPostCreated={(post: PostInterface) =>
              setPosts((prev) => [post, ...prev])
            }
          />
        </div>
        <div className="mt-8 flex w-full max-w-4xl flex-col gap-4">
          {renderPosts()}
        </div>
      </section>
      <aside className="flex w-[250px] flex-col gap-4 border-r p-4 backdrop-blur-md">
        {/* What's happening Section */}
      </aside>
    </main>
  );
}
