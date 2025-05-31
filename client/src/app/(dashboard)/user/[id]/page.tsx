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
  type AddOrAcceptFriendResponse,
} from "~/types/types";
import profile_1 from "~/assets/profile/profile-picture-1.png";
import profile_2 from "~/assets/profile/profile-picture-2.png";
import profile_3 from "~/assets/profile/profile-picture-3.png";
import profile_4 from "~/assets/profile/profile-picture-4.png";
import profile_5 from "~/assets/profile/profile-picture-5.png";
import Post from "~/components/post";
import Cookies from "js-cookie";
import { DecodeUserFromToken } from "~/utils/utils";

export default function Page() {
  const params = useParams();
  const userID = params.id;
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [posts, setPosts] = useState<PostInterface[]>([]);
  const [friends, setFriends] = useState<FriendInterface[]>([]);
  const [friendRequests, setFriendRequests] = useState<FriendInterface[]>([]);
  const [loggedInUserFriendRequests, setLoggedInUserFriendRequests] = useState<
    FriendInterface[]
  >([]);
  const [loggedInUser, setLoggedInUser] = useState<User>({
    id: 0,
    username: "",
    email: "",
    fullname: "",
    createdAt: "",
  });
  const [user, setUser] = useState<User>({
    id: 0,
    username: "",
    email: "",
    fullname: "",
    createdAt: "",
  });

  useEffect(() => {
    const token = Cookies.get("token");
    if (token) {
      const user = DecodeUserFromToken(token);
      if (user) {
        setLoggedInUser({
          id: user.user.id,
          username: user.user.username,
          email: user.user.email,
          fullname: user.user.fullname,
          createdAt: user.user.created_at,
        });
        setIsLoggedIn(true);
      }
    }
  }, []);

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
            id: data.id,
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

    const fetchUserFriendRequests = async (userID: string) => {
      try {
        const response = await axios.get<FriendResponse>(
          `http://localhost:8080/api/friend/requests/${userID}`,
          {
            headers: {
              "Content-Type": "application/json",
            },
          },
        );
        if (response.status === 200) {
          const data = response.data.data;
          setFriendRequests(data);
        }
      } catch (error) {
        // TODO: Handle error using toast
        console.error("Error fetching user friend requests:", error);
      }
    };

    const fetchMyFriendRequests = async (userID: string) => {
      try {
        const response = await axios.get<FriendResponse>(
          `http://localhost:8080/api/friend/requests/${userID}`,
          {
            headers: {
              "Content-Type": "application/json",
            },
          },
        );
        if (response.status === 200) {
          const data = response.data.data;
          setLoggedInUserFriendRequests(data);
        }
      } catch (error) {
        console.error("Error fetching logged-in user friend requests:", error);
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
    fetchUserFriendRequests(userID as string).catch((error) => {
      console.error("Failed to fetch user friend requests:", error);
    });
    if (isLoggedIn) {
      fetchMyFriendRequests(loggedInUser.id.toString()).catch((error) => {
        console.error("Failed to fetch logged-in user friend requests:", error);
      });
    }
  }, [userID, isLoggedIn, loggedInUser.id]);

  const profilePictures = [
    profile_1,
    profile_2,
    profile_3,
    profile_4,
    profile_5,
  ];
  const [randomIndex] = useState(() =>
    Math.floor(Math.random() * profilePictures.length),
  );

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

  const addFriend = async () => {
    if (!isLoggedIn || loggedInUser.username === user.username) return;

    try {
      await axios.post<AddOrAcceptFriendResponse>(
        `http://localhost:8080/api/friend/add`,
        {
          userid: loggedInUser.id,
          frienduserid: user.id,
        },
        {
          headers: {
            Authorization: `Bearer ${Cookies.get("token")}`,
            "Content-Type": "application/json",
          },
        },
      );
    } catch (error) {
      // TODO: Handle error using toast
      console.error("Error adding friend:", error);
    } finally {
      location.reload();
    }
  };

  const acceptFriendRequest = async (friendID: number) => {
    if (!isLoggedIn || loggedInUser.username === user.username) return;

    try {
      await axios.post<AddOrAcceptFriendResponse>(
        `http://localhost:8080/api/friend/accept`,
        {
          userid: loggedInUser.id,
          frienduserid: friendID, // ini orang yang mau di-accept
        },
        {
          headers: {
            Authorization: `Bearer ${Cookies.get("token")}`,
            "Content-Type": "application/json",
          },
        },
      );
    } catch (error) {
      // TODO: Handle error using toast
      console.error("Error accepting friend request:", error);
    } finally {
      location.reload();
    }
  };

  const removeFriend = async (userID: number) => {
    if (!isLoggedIn || loggedInUser.username === user.username) return;

    // step 1: cari friendID dari friends (kalau userid = userID dan frienduserid = loggedInUser.id)
    const friend = friends.find((f) => {
      return (
        (f.userid === userID && f.frienduserid === loggedInUser.id) ||
        (f.userid === loggedInUser.id && f.frienduserid === userID)
      );
    });

    // step 2: ambil id-nya
    const friendID = friend?.id;
    if (!friendID) return;

    // step 3: panggil API untuk menghapus friend
    try {
      await axios.delete<AddOrAcceptFriendResponse>(
        `http://localhost:8080/api/friend/delete/${friendID}`,
        {
          headers: {
            Authorization: `Bearer ${Cookies.get("token")}`,
            "Content-Type": "application/json",
          },
        },
      );
    } catch (error) {
      // TODO: Handle error using toast
      console.error("Error removing friend:", error);
    } finally {
      location.reload();
    }
  };

  return (
    <main className="grid h-screen w-full">
      <section className="px-6">
        <div className="mx-auto mt-4 w-full">
          {/* Bagian Profile Atas */}
          <div className="relative flex flex-row justify-between rounded-xl bg-gradient-to-tr from-white to-blue-300 p-10 shadow-md">
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

            {/* Bottom-right Add Friend Button */}
            {isLoggedIn &&
              loggedInUser.id !== user.id &&
              !friends.some(
                (f) =>
                  (f.userid === loggedInUser.id &&
                    f.frienduserid === user.id) ||
                  (f.userid === user.id && f.frienduserid === loggedInUser.id),
              ) &&
              (friendRequests.some((f) => f.userid === loggedInUser.id) ? (
                <button
                  disabled
                  className="absolute right-4 bottom-4 items-center justify-center rounded-lg bg-yellow-400 px-6 py-3 font-bold text-white shadow-xl"
                >
                  Pending Request
                </button>
              ) : (
                <button
                  onClick={addFriend}
                  className="absolute right-4 bottom-4 items-center justify-center rounded-lg bg-blue-600 px-6 py-3 font-bold text-white shadow-xl transition-colors duration-300 hover:bg-blue-700"
                >
                  Add Friend
                </button>
              ))}

            {isLoggedIn &&
              loggedInUser.id !== user.id &&
              loggedInUserFriendRequests.some((f) => f.userid === user.id) && (
                <button
                  onClick={() => acceptFriendRequest(user.id)}
                  className="absolute right-4 bottom-4 items-center justify-center rounded-lg bg-green-600 px-6 py-3 font-bold text-white shadow-xl transition-colors duration-300 hover:bg-green-700"
                >
                  Accept Request
                </button>
              )}

            {isLoggedIn &&
              loggedInUser.id !== user.id &&
              friends.some(
                (f) =>
                  (f.userid === loggedInUser.id &&
                    f.frienduserid === user.id) ||
                  (f.userid === user.id && f.frienduserid === loggedInUser.id),
              ) && (
                <button
                  onClick={() => removeFriend(user.id)}
                  className="absolute right-4 bottom-4 items-center justify-center rounded-lg bg-red-600 px-6 py-3 font-bold text-white shadow-xl transition-colors duration-300 hover:bg-red-700"
                >
                  Remove Friend
                </button>
              )}
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
