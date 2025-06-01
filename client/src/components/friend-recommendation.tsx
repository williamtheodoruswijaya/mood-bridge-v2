"use client";
import axios from "axios";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import profile_1 from "~/assets/profile/profile-picture-1.png";
import profile_2 from "~/assets/profile/profile-picture-2.png";
import profile_3 from "~/assets/profile/profile-picture-3.png";
import profile_4 from "~/assets/profile/profile-picture-4.png";
import profile_5 from "~/assets/profile/profile-picture-5.png";
import {
  type FriendRecommendationInterface,
  type FriendRecommendationResponse,
} from "~/types/types";

interface FriendRecommendationProps {
  _userID: string;
  _token?: string;
}

const FriendRecommendation = ({
  _userID,
  _token,
}: FriendRecommendationProps) => {
  const [friendRecommendations, setFriendRecommendations] = useState<
    FriendRecommendationInterface[]
  >([]);
  const router = useRouter();
  const categoryColor: Record<string, string> = {
    Normal: "#219E2C",
    Anxiety: "#FFAE00",
    Depression: "#0D00FF",
    Suicidal: "#FF0000",
    Stress: "#FF00A0",
    Bipolar: "#8B00FF",
    "Personality Disorder": "#000000",
  };
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
  function mergeFriendRecommendations(
    prevRecommendations: FriendRecommendationInterface[],
    newRecommendations: FriendRecommendationInterface[],
  ) {
    const existingIds = new Set(prevRecommendations.map((p) => p.userid));
    const filteredNewRecommendations = newRecommendations.filter(
      (p) => !existingIds.has(p.userid),
    );
    return [...prevRecommendations, ...filteredNewRecommendations];
  }

  useEffect(() => {
    if (!_userID || !_token) {
      console.error("User ID or token is not provided.");
      return;
    }
    const fetchFriendRecommendations = async () => {
      try {
        const response = await axios.get<FriendRecommendationResponse>(
          `http://localhost:8080/api/friend/recommendation/${_userID}`,
          {
            headers: {
              Authorization: `Bearer ${_token}`,
              "Content-Type": "application/json",
            },
          },
        );
        if (response.data.code === 200) {
          const data = response.data.data;
          setFriendRecommendations((prev) =>
            mergeFriendRecommendations(prev, data),
          );
        } else {
          console.error(
            "Failed to fetch friend recommendations:",
            response.data.message,
          );
        }
      } catch (error) {
        console.error("Error fetching friend recommendations:", error);
      }
    };
    fetchFriendRecommendations().catch((error) => {
      console.error("Error in fetchFriendRecommendations:", error);
      setFriendRecommendations([]);
    });
  }, [_userID, _token]);
  return (
    <div className="mx-auto max-w-md space-y-4 rounded-xl bg-[#ffffff] p-6 shadow-lg">
      <h1 className="text-xl font-semibold text-gray-800">Who to Follow</h1>
      {friendRecommendations.length > 0 ? (
        <>
          {friendRecommendations.map((recommendation) => (
            <button
              onClick={() => router.push(`/user/${recommendation.userid}`)}
              key={recommendation.userid}
              className="flex w-full items-center gap-4 rounded-lg p-2 text-left transition hover:bg-[#f3faff]"
            >
              <Image
                src={getProfilePicture(recommendation.userid.toString())!.src}
                alt={recommendation.fullname}
                width={64}
                height={64}
                className="flex-shrink-0 rounded-full object-cover"
              />
              <div className="grid w-full grid-rows-[auto_auto]">
                <div className="flex w-full items-center justify-between">
                  <h2 className="truncate text-sm font-semibold text-gray-800">
                    {recommendation.fullname}
                  </h2>
                  <div
                    className="ml-2 rounded-md px-3 py-0.5 text-xs font-semibold whitespace-nowrap text-white"
                    style={{
                      backgroundColor:
                        categoryColor[recommendation.overall_mood] ?? "#687669",
                    }}
                  >
                    {recommendation.overall_mood}
                  </div>
                </div>
                <p className="truncate text-xs text-gray-600">
                  @{recommendation.username}
                </p>
              </div>
            </button>
          ))}
        </>
      ) : (
        <div className="text-center text-gray-500">
          No friend recommendations available.
        </div>
      )}
    </div>
  );
};

export default FriendRecommendation;
