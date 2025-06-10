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
  type LoginResponse,
} from "~/types/types";
import profile_1 from "~/assets/profile/profile-picture-1.png";
import profile_2 from "~/assets/profile/profile-picture-2.png";
import profile_3 from "~/assets/profile/profile-picture-3.png";
import profile_4 from "~/assets/profile/profile-picture-4.png";
import profile_5 from "~/assets/profile/profile-picture-5.png";
import Post from "~/components/post";
import Cookies from "js-cookie";
import { DecodeUserFromToken } from "~/utils/utils";
import { TbPencil, TbPencilCancel } from "react-icons/tb";

export default function Page() {
  const params = useParams();
  const userID = params.id;
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [posts, setPosts] = useState<PostInterface[]>([]);
  const [friends, setFriends] = useState<FriendInterface[]>([]);
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
  const [isEditing, setIsEditing] = useState(false);
  const [editUser, setEditUser] = useState({
    username: loggedInUser.username,
    fullname: loggedInUser.fullname,
    email: loggedInUser.email,
    password: "",
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setEditUser((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const reLogin = async (
    loggedInUsername: string,
    loggedInPassword: string,
  ) => {
    // buat update JWT Token setelah update
    try {
      const response = await axios.post<LoginResponse>(
        `${process.env.NEXT_PUBLIC_API_URL}/api/user/login`,
        {
          username: loggedInUsername,
          password: loggedInPassword,
        },
        {
          headers: {
            "Content-Type": "application/json",
          },
        },
      );
      if (response.status === 200) {
        const token = response.data.data;
        Cookies.set("token", token, { expires: 7 });
        location.reload(); // Reload the page to reflect changes
      }
    } catch (error) {
      // TODO: Handle error using toast
      console.error("Error re-logging in:", error);
    }
  };

  const handleSaveChanges = async () => {
    if (
      !editUser.username ||
      !editUser.fullname ||
      !editUser.email ||
      !editUser.password
    ) {
      alert("Please fill in all fields."); // TODO: GANTI SAMA TOAST
      setIsEditing(false);
    }
    try {
      const response = await axios.put<RegisterResponse>(
        `${process.env.NEXT_PUBLIC_API_URL}/api/user/update/${loggedInUser.id}`,
        {
          username: editUser.username,
          fullname: editUser.fullname,
          email: editUser.email,
          password: editUser.password,
          profile: "",
        },
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${Cookies.get("token")}`,
          },
        },
      );
      if (response.status === 200) {
        const data = response.data.data;
        setLoggedInUser({
          id: data.id,
          username: data.username,
          email: data.email,
          fullname: data.fullname,
          createdAt: data.createdAt,
        });
        setEditUser({
          username: data.username,
          fullname: data.fullname,
          email: data.email,
          password: "",
        });
        // Update user token as well:
        await reLogin(editUser.username, editUser.password);
        location.reload();
      }
    } catch (error) {
      // TODO: Handle error using toast
      console.error("Error updating user:", error);
    } finally {
      setIsEditing(false);
    }
  };

  const handleCancel = () => {
    setEditUser({
      username: loggedInUser.username,
      fullname: loggedInUser.fullname,
      email: loggedInUser.email,
      password: "",
    });
    setIsEditing(false);
  };

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
        setEditUser({
          username: user.user.username,
          fullname: user.user.fullname,
          email: user.user.email,
          password: "",
        });
        setIsLoggedIn(true);
      }
    }
  }, []);

  useEffect(() => {
    const fetchUserData = async (userID: string) => {
      try {
        const response = await axios.get<RegisterResponse>(
          `${process.env.NEXT_PUBLIC_API_URL}/api/user/by-id/${userID}`,
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
          `${process.env.NEXT_PUBLIC_API_URL}/api/post/by-userid/${userID}`,
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
          `${process.env.NEXT_PUBLIC_API_URL}/api/friend/all/${userID}`,
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
  }, [userID, isLoggedIn, loggedInUser.id]);

  const profilePictures = [
    profile_1,
    profile_2,
    profile_3,
    profile_4,
    profile_5,
  ];
  const getProfilePicture = (userid: string) => {
    const hash = Array.from(userid).reduce(
      (acc, char) => acc + char.charCodeAt(0),
      0,
    );
    const index = hash % profilePictures.length;
    return profilePictures[index];
  };

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
          <div className="relative flex flex-row justify-between rounded-xl bg-gradient-to-tr from-white to-blue-300 p-10 shadow-md">
            {/* Tombol Edit */}
            {loggedInUser.id === user.id && (
              <button
                onClick={() => setIsEditing(!isEditing)}
                disabled={!isLoggedIn}
                type="button"
                aria-label="Edit Profile"
                className="absolute top-4 right-4"
              >
                {isEditing ? (
                  <TbPencilCancel className="h-8 w-8" />
                ) : (
                  <TbPencil className="h-8 w-8" />
                )}
              </button>
            )}
            <div className="flex flex-row items-start gap-10">
              <Image
                src={getProfilePicture(user.id.toString())!.src}
                width={160}
                height={160}
                alt="profile-picture"
                className="rounded-full object-cover shadow-lg"
              />
              <div className="flex flex-col justify-center">
                {isEditing ? (
                  <>
                    <input
                      type="text"
                      name="fullname"
                      placeholder="Full Name"
                      value={editUser.fullname}
                      onChange={handleInputChange}
                      className="mb-2 rounded bg-white px-2 py-1 text-3xl font-bold"
                    />
                    <input
                      type="text"
                      name="username"
                      placeholder="Username"
                      value={editUser.username}
                      onChange={handleInputChange}
                      className="mb-2 rounded bg-white px-2 py-1 text-lg text-gray-600"
                    />
                    <input
                      type="email"
                      name="email"
                      placeholder="Email"
                      value={editUser.email}
                      onChange={handleInputChange}
                      className="text-md mb-2 rounded bg-white px-2 py-1 text-gray-500"
                    />
                    <input
                      type="password"
                      name="password"
                      placeholder="Password"
                      value={editUser.password}
                      onChange={handleInputChange}
                      className="mb-4 rounded bg-white px-2 py-1 text-sm"
                    />
                    <div className="flex gap-2">
                      <button
                        onClick={handleSaveChanges}
                        className="rounded bg-green-500 px-5 py-1 text-white"
                      >
                        Save
                      </button>
                      <button
                        onClick={handleCancel}
                        className="rounded bg-red-500 px-3 py-1 text-white"
                      >
                        Cancel
                      </button>
                    </div>
                  </>
                ) : (
                  <>
                    <h1 className="text-3xl font-bold">{user.fullname}</h1>
                    <p className="text-lg text-gray-600">@{user.username}</p>
                    <p className="text-md text-gray-500">{user.email}</p>
                  </>
                )}

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
