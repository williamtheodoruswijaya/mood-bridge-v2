import type { CommentInterface } from "~/types/types";
import { TimeAgo } from "~/utils/utils";

const Comment: React.FC<CommentInterface> = (props) => {
  return (
    <div className="w-full rounded-xl bg-white p-5 shadow-md">
      <div className="flex items-start justify-between">
        <div className="text-lg font-bold text-black">
          {props.user.fullname}{" "}
          <span className="text-md font-normal text-gray-800">
            @{props.user.username}
          </span>
          <span className="px-2 text-xs font-light text-gray-800">
            {TimeAgo(props.created_at)}
          </span>
        </div>
      </div>
      <p className="mt-2 text-sm text-black">{props.content}</p>
    </div>
  );
};

export default Comment;
