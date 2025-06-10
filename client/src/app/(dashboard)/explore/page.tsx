"use client";
import axios from "axios";
import Cookies from "js-cookie";
import { useEffect, useRef, useState } from "react";
import FriendRecommendation from "~/components/friend-recommendation";
import News from "~/components/news";
import Post from "~/components/post";
import type { PostInterface, PostResponse, User } from "~/types/types";
import { DecodeUserFromToken } from "~/utils/utils";

export default function Page() {
  const scrollRef = useRef<HTMLDivElement>(null);
  const [limit] = useState(10);
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [posts, setPosts] = useState<PostInterface[]>([]);
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState<User>({
    id: 0,
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
      setIsLoggedIn(false);

      if (token) {
        const user = DecodeUserFromToken(token);
        if (user) {
          setUser({
            id: user.user.id,
            username: user.user.username,
            email: user.user.email,
            fullname: user.user.fullname,
            createdAt: user.user.created_at,
          });
          setIsLoggedIn(true);
        }
      }
      setLoading(true);
      const fetchAllPosts = async () => {
        const apiUrl = `${process.env.NEXT_PUBLIC_API_URL}/api/post/all?limit=${limit}&offset=${offset}`;
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

      await fetchAllPosts();
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
  }, [loading, hasMore]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (offset === 0) return;

    const fetchMorePosts = async () => {
      setLoading(true);
      try {
        const apiUrl = `${process.env.NEXT_PUBLIC_API_URL}/api/post/all?limit=${limit}&offset=${offset}`;
        const response = await axios.get<PostResponse>(apiUrl);
        if (response.status === 200) {
          const newPosts = response.data.data;
          setPosts((prev) => mergePostsUnique(prev, newPosts));
          if (newPosts.length < limit || newPosts.length === 0) {
            setHasMore(false);
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
  }, [offset, isLoggedIn, user.id]); // eslint-disable-line react-hooks/exhaustive-deps

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
        <div className="mx-auto mt-4 max-w-4xl"></div>
        <div className="mx-auto mt-8 flex max-w-4xl flex-col gap-4">
          {renderPosts()}
        </div>
        {loading && offset !== 0 && (
          <p className="text-center text-black">Loading more posts...</p>
        )}
      </section>

      <aside className="border-l p-4 backdrop-blur-md">
        {isLoggedIn && (
          <div className="mx-auto mb-8 max-w-md">
            <FriendRecommendation
              _userID={user.id.toString()}
              _token={Cookies.get("token")}
            />
          </div>
        )}
        <News />
      </aside>
    </main>
  );
}
