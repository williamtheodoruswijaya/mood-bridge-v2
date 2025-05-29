"use client";
import axios from "axios";
import { useEffect, useState } from "react";
import CreatePost from "~/components/create-post";
import Post from "~/components/post";
import type { PostInterface, PostResponse } from "~/types/types";

export default function HomePage() {
  const [posts, setPosts] = useState<PostInterface[]>([]);
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const response = await axios.get<PostResponse>(
          "http://localhost:8080/api/post/all",
        );
        if (response.status === 200) {
          setPosts(response.data.data);
        } else {
          console.error("Failed to fetch posts:", response.data.message);
        }
      } catch (error) {
        console.error("Error fetching posts:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchPosts().catch((error) => {
      console.error("Error in useEffect:", error);
      setLoading(false);
    });
  }, []);
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
