"use client";
import { useRouter } from "next/navigation";
import { MdComment } from "react-icons/md";
import type { PostInterface } from "~/types/types";

const Post: React.FC<PostInterface> = (props) => {
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
  return (
    <div className="w-full rounded-xl bg-[#84E7EE] p-4 shadow-lg backdrop-blur-md">
      <div className="flex items-start justify-between">
        <div className="text-lg font-bold text-black">
          {props.user.fullname}{" "}
          <span className="text-sm font-normal text-gray-800">
            @{props.user.username}
          </span>
        </div>
        <div
          className="min-w-20 rounded-md px-3 py-1 text-center text-sm font-semibold text-white"
          style={{ backgroundColor: categoryColor[props.mood] ?? "#687669" }}
        >
          {props.mood}
        </div>
      </div>
      <p className="mt-2 text-sm text-black">{props.content}</p>
      <button
        className="mt-4 flex items-center text-sm text-gray-700 hover:text-blue-600"
        onClick={() => router.push(`/comment/${props.postid}`)}
      >
        <span className="mr-2">
          <MdComment />
        </span>{" "}
        View Comments
      </button>
    </div>
  );
};

export default Post;
