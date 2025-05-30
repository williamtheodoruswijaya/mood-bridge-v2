"use client";
import axios from "axios";
import Cookies from "js-cookie";
import { useEffect, useRef, useState } from "react";
import CreatePost from "~/components/create-post";
import News from "~/components/news";
import Post from "~/components/post";
import type { PostInterface, PostResponse, User } from "~/types/types";

export default function HomePage() {
  const scrollRef = useRef<HTMLDivElement>(null);
  const [limit] = useState(10);
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);
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

  // Helper function untuk gabungkan posts tanpa duplikat berdasarkan postid
  function mergePostsUnique(
    prevPosts: PostInterface[],
    newPosts: PostInterface[],
  ) {
    const existingIds = new Set(prevPosts.map((p) => p.postid));
    const filteredNewPosts = newPosts.filter((p) => !existingIds.has(p.postid));
    return [...prevPosts, ...filteredNewPosts];
  }

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
        const apiUrl = `http://localhost:8080/api/post/all?limit=${limit}&offset=${offset}`;
        try {
          const response = await axios.get<PostResponse>(apiUrl);
          if (response.status === 200) {
            const newPosts = response.data.data;
            setPosts((prev) => mergePostsUnique(prev, newPosts));
            if (newPosts.length < limit || newPosts.length === 0) {
              setHasMore(false);
            }
          } else {
            console.error("Failed to fetch all posts:", response.statusText);
          }
        } catch (error) {
          console.error("Error fetching all posts:", error);
          setHasMore(false);
        } finally {
          setLoading(false);
        }
      };

      const fetchFriendPosts = async (
        userID: number,
        limit: number,
        offset: number,
      ) => {
        const apiUrl = `http://localhost:8080/api/post/friend-posts/${userID}?limit=${limit}&offset=${offset}`;
        try {
          const response = await axios.get<PostResponse>(apiUrl, {
            headers: {
              Authorization: `Bearer ${Cookies.get("token")}`,
              "Content-Type": "application/json",
            },
          });
          if (response.status === 200) {
            const friendPosts = response.data.data;
            setPosts((prev) => mergePostsUnique(prev, friendPosts));
            if (friendPosts.length < limit || friendPosts.length === 0) {
              setHasMore(false);
            }
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
        await fetchFriendPosts(parsedUser.user.id, limit, offset);
      } else {
        await fetchAllPosts();
      }
    };

    fetchUserAndPosts().catch((error) => {
      console.error("Error in fetchUserAndPosts:", error);
      setLoading(false);
    });
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    const container = scrollRef.current;
    if (!container) return;

    const handleScroll = () => {
      const { scrollTop, scrollHeight, clientHeight } = container;
      if (
        scrollTop + clientHeight >= scrollHeight - 200 &&
        !loading &&
        hasMore
      ) {
        setLoading(true);
        setOffset((prev) => prev + limit);
      }
    };

    container.addEventListener("scroll", handleScroll);
    return () => container.removeEventListener("scroll", handleScroll);
  }, [loading, hasMore]);

  useEffect(() => {
    if (offset === 0) return;

    const fetchMorePosts = async () => {
      setLoading(true);
      try {
        if (isLoggedIn && user.userID > 0) {
          const apiUrl = `http://localhost:8080/api/post/friend-posts/${user.userID}?limit=${limit}&offset=${offset}`;
          const response = await axios.get<PostResponse>(apiUrl, {
            headers: {
              Authorization: `Bearer ${Cookies.get("token")}`,
              "Content-Type": "application/json",
            },
          });
          if (response.status === 200) {
            const friendPosts = response.data.data;
            setPosts((prev) => mergePostsUnique(prev, friendPosts));
            if (friendPosts.length < limit || friendPosts.length === 0) {
              setHasMore(false);
            }
          }
        } else {
          const apiUrl = `http://localhost:8080/api/post/all?limit=${limit}&offset=${offset}`;
          const response = await axios.get<PostResponse>(apiUrl);
          if (response.status === 200) {
            const newPosts = response.data.data;
            setPosts((prev) => [...prev, ...newPosts]);
            if (newPosts.length < limit || newPosts.length === 0) {
              setHasMore(false);
            }
          }
        }
      } catch (error) {
        console.error(error);
        setHasMore(false);
      } finally {
        setLoading(false);
      }
    };

    fetchMorePosts().catch((error) => {
      console.error("Error fetching more posts:", error);
      setLoading(false);
      setHasMore(false);
    });
  }, [offset, isLoggedIn, user.userID]); // eslint-disable-line react-hooks/exhaustive-deps

  function renderPosts() {
    if (loading && posts.length === 0)
      return <p className="text-black">Loading posts...</p>;
    if (posts.length === 0)
      return <p className="text-black">No posts available</p>;
    return posts.map((post) => <Post key={post.postid} {...post} />);
  }

  return (
    <main className="grid h-screen w-full grid-cols-[1fr_400px] text-white">
      <section ref={scrollRef} className="overflow-y-auto px-6">
        <div className="mx-auto mt-4 max-w-4xl">
          <CreatePost
            onPostCreated={(post: PostInterface) =>
              setPosts((prev) => [post, ...prev])
            }
          />
        </div>
        <div className="mx-auto mt-8 flex max-w-4xl flex-col gap-4">
          {renderPosts()}
        </div>
        {loading && offset !== 0 && (
          <p className="text-center text-black">Loading more posts...</p>
        )}
      </section>

      <aside className="overflow-y-auto border-l p-4 backdrop-blur-md">
        <News />
      </aside>
    </main>
  );
}
