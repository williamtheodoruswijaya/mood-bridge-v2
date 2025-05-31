import { useRouter } from "next/navigation";
import { MdComment } from "react-icons/md";
import type { PostInterface } from "~/types/types";
import { TimeAgo } from "~/utils/utils";

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
    <div className="w-full rounded-xl border-gray-200 bg-white p-5 shadow-md">
      <div className="flex items-start justify-between">
        <button
          className="text-black flex items-center text-lg font-bold text-black"
          onClick={() => router.push(`/user/${props.userid}`)}
        >
          <div className="hover:text-blue-800 hover:underline">
            {props.user.fullname}{" "}
            <span className="text-sm font-normal text-gray-800">
              @{props.user.username}
            </span>
          </div>
          <span className="px-2 text-xs font-medium text-gray-600">
            {TimeAgo(props.createdat)}
          </span>
        </button>
        <div
          className="min-w-20 rounded-md px-3 py-1 text-center text-sm font-semibold text-white"
          style={{ backgroundColor: categoryColor[props.mood] ?? "#687669" }}
        >
          {props.mood}
        </div>
      </div>
      <p className="mt-2 text-sm text-black">{props.content}</p>
      <button
        className="mt-4 flex items-center text-sm text-gray-700 hover:text-cyan-600"
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
