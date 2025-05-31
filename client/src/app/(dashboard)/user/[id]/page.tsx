"use client";
import axios from "axios";
import Image from "next/image";
import { useParams } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import {
  type PostResponse,
  type PostInterface,
  type RegisterResponse,
  type User,
  type FriendInterface,
  type FriendResponse,
} from "~/types/types";
import profile_1 from "~/assets/profile/profile-picture-1.png";
import profile_2 from "~/assets/profile/profile-picture-2.png";
import profile_3 from "~/assets/profile/profile-picture-3.png";
import profile_4 from "~/assets/profile/profile-picture-4.png";
import profile_5 from "~/assets/profile/profile-picture-5.png";
import Post from "~/components/post";

export default function Page() {
  const params = useParams();
  const userID = params.id;
  const [posts, setPosts] = useState<PostInterface[]>([]);
  const [friends, setFriends] = useState<FriendInterface[]>([]);
  const [user, setUser] = useState<User>({
    userID: 0,
    username: "",
    email: "",
    fullname: "",
    createdAt: "",
  });

  useEffect(() => {
    const fetchUserData = async (userID: string) => {
      try {
        const response = await axios.get<RegisterResponse>(
          `http://localhost:8080/api/user/by-id/${userID}`,
          {
            headers: {
              "Content-Type": "application/json",
            },
          },
        );
        if (response.status === 200) {
          const data = response.data.data;
          setUser({
            userID: data.userID,
            username: data.username,
            email: data.email,
            fullname: data.fullname,
            createdAt: data.createdAt,
          });
        }
      } catch (error) {
        // TODO: Handle error using toast
        console.error("Error fetching user data:", error);
      }
    };

    const fetchUserPosts = async (userID: string) => {
      try {
        const response = await axios.get<PostResponse>(
          `http://localhost:8080/api/post/by-userid/${userID}`,
          {
            headers: {
              "Content-Type": "application/json",
            },
          },
        );
        if (response.status === 200) {
          const data = response.data.data;
          setPosts(data);
        }
      } catch (error) {
        // TODO: Handle error using toast
        console.error("Error fetching user posts:", error);
      }
    };

    const fetchUserFriends = async (userID: string) => {
      try {
        const response = await axios.get<FriendResponse>(
          `http://localhost:8080/api/friend/all/${userID}`,
          {
            headers: {
              "Content-Type": "application/json",
            },
          },
        );
        if (response.status === 200) {
          const data = response.data.data;
          setFriends(data);
        }
      } catch (error) {
        // TODO: Handle error using toast
        console.error("Error fetching user friends:", error);
      }
    };

    fetchUserData(userID as string).catch((error) => {
      console.error("Failed to fetch user data:", error);
    });
    fetchUserPosts(userID as string).catch((error) => {
      console.error("Failed to fetch user posts:", error);
    });
    fetchUserFriends(userID as string).catch((error) => {
      console.error("Failed to fetch user friends:", error);
    });
  }, [userID]);

  const profilePictures = [
    profile_1,
    profile_2,
    profile_3,
    profile_4,
    profile_5,
  ];
  const randomIndex = useMemo(() => {
    return Math.floor(Math.random() * profilePictures.length);
  }, [profilePictures.length]);

  const overallMood = useMemo(() => {
    const moodCount: Record<string, number> = {};
    posts.forEach((post) => {
      moodCount[post.mood] = (moodCount[post.mood] ?? 0) + 1;
    });
    return Object.entries(moodCount).reduce(
      (prev, current) => (current[1] > prev[1] ? current : prev),
      ["neutral", 0],
    )[0];
  }, [posts]);

  const categoryColor: Record<string, string> = {
    Normal: "#219E2C",
    Anxiety: "#FFAE00",
    Depression: "#0D00FF",
    Suicidal: "#FF0000",
    Stress: "#FF00A0",
    Bipolar: "#8B00FF",
    "Personality Disorder": "#000000",
  };
  return (
    <main className="grid h-screen w-full">
      <section className="px-6">
        <div className="mx-auto mt-4 w-full">
          {/* Bagian Profile Atas */}
          <div className="flex flex-row justify-between rounded-xl bg-gradient-to-tr from-white to-blue-300 p-10 shadow-md">
            <div className="flex flex-row">
              <Image
                src={profilePictures[randomIndex] ?? profile_1}
                alt="profile-picture"
                className="h-40 w-40 rounded-full object-cover shadow-lg"
              />
              <div className="ml-10 flex flex-col justify-center">
                <h1 className="text-3xl font-bold">{user.fullname}</h1>
                <p className="text-lg text-gray-600">@{user.username}</p>
                <p className="text-md text-gray-500">{user.email}</p>
                <div className="mt-4 flex gap-4">
                  <div className="w-32 rounded-lg bg-white p-2 text-center shadow-md">
                    <p className="text-xl font-semibold">{friends.length}</p>
                    <p className="text-xs text-gray-500">Friends</p>
                  </div>
                  <div className="w-32 rounded-lg bg-white p-2 text-center shadow-md">
                    <p className="text-xl font-semibold">{posts.length}</p>
                    <p className="text-xs text-gray-500">Post</p>
                  </div>
                  <div
                    className="w-32 rounded-lg p-2 text-center text-white shadow-md"
                    style={{
                      backgroundColor: categoryColor[overallMood] ?? "#687669",
                    }}
                  >
                    <p className="text-xl font-semibold">{overallMood}</p>
                    <p className="text-xs">Overall Mood</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Bagian Postingan */}
          <div className="mt-6">
            {posts.length > 0 ? (
              posts.map((post) => (
                <div key={post.postid} className="py-1">
                  <Post {...post} />
                </div>
              ))
            ) : (
              <p className="text-center text-gray-500">No posts available</p>
            )}
          </div>
        </div>
      </section>
    </main>
  );
}
