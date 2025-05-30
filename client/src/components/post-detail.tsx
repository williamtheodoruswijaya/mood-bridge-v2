import type { PostInterface } from "~/types/types";

const PostDetail: React.FC<PostInterface> = (props) => {
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
    <div className="w-full rounded-xl border border-gray-200 bg-white p-5 shadow-md">
      <div className="flex items-start justify-between">
        <div className="flex flex-col">
          <span className="text-lg font-semibold text-gray-900">
            {props.user.fullname}
          </span>
          <span className="text-sm text-gray-500">@{props.user.username}</span>
        </div>
        <div
          className="rounded-full px-3 py-1 text-xs font-bold tracking-wide text-white uppercase"
          style={{ backgroundColor: categoryColor[props.mood] ?? "#687669" }}
        >
          {props.mood}
        </div>
      </div>
      <p className="mt-3 text-base text-gray-800">{props.content}</p>
    </div>
  );
};

export default PostDetail;
