"use client";
import axios from "axios";
import Image from "next/image";
import { useParams } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import type { RegisterResponse, User } from "~/types/types";
import profile_1 from "~/assets/profile/profile-picture-1.png";
import profile_2 from "~/assets/profile/profile-picture-2.png";
import profile_3 from "~/assets/profile/profile-picture-3.png";
import profile_4 from "~/assets/profile/profile-picture-4.png";
import profile_5 from "~/assets/profile/profile-picture-5.png";

export default function Page() {
  const params = useParams();
  const userID = params.id;
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
    fetchUserData(userID as string).catch((error) => {
      console.error("Failed to fetch user data:", error);
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

  return (
    <main className="grid h-screen w-full">
      <section className="px-6">
        <div className="mx-auto mt-4 w-full">
          <div className="flex flex-row justify-between rounded-lg bg-gradient-to-t from-cyan-100 to-white p-10 shadow-md">
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
                  <div className="w-24 rounded border p-2 text-center shadow">
                    <p className="text-xl font-semibold">30</p>
                    <p className="text-sm text-gray-500">Friends</p>
                  </div>
                  <div className="w-24 rounded border p-2 text-center shadow">
                    <p className="text-xl font-semibold">15</p>
                    <p className="text-sm text-gray-500">Post</p>
                  </div>
                </div>
              </div>
            </div>

            {/* Right: Overall Mood */}
            <div className="flex flex-row items-end justify-start">
              <p className="text-xl font-medium">Overall mood:</p>
              <div className="mt-2 h-10 w-10 rounded border-2 border-gray-500" />
            </div>
          </div>
        </div>
      </section>
    </main>
  );
}
