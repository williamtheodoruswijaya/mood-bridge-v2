import { useRouter } from "next/navigation";
import { MdComment } from "react-icons/md";
import type { PostInterface } from "~/types/types";
import { TimeAgo } from "~/utils/utils";
import profile_1 from "~/assets/profile/profile-picture-1.png";
import profile_2 from "~/assets/profile/profile-picture-2.png";
import profile_3 from "~/assets/profile/profile-picture-3.png";
import profile_4 from "~/assets/profile/profile-picture-4.png";
import profile_5 from "~/assets/profile/profile-picture-5.png";
import Image from "next/image";

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
  return (
    <div className="w-full rounded-xl border-gray-200 bg-white p-5 shadow-md">
      <div className="flex items-start justify-between">
        <div className="flex items-center">
          <Image
            src={getProfilePicture(props.userid.toString()).src}
            width={40}
            height={40}
            alt="Profile-Picture"
            className="mr-3 rounded-full object-cover"
          />
          <button
            className="flex items-center text-lg font-bold text-black"
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
