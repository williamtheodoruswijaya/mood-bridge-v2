// components/Sidebar.tsx
"use client";

import Link from "next/link";
import Image from "next/image";
import Icon from "~/assets/icon.png";
import { MdDashboard, MdLogout } from "react-icons/md";
import { FaCompass } from "react-icons/fa";
import { BsChatSquareDots } from "react-icons/bs";
import { PiStarFourFill } from "react-icons/pi";
import { usePathname, useRouter } from "next/navigation";
import Cookies from "js-cookie";

const items = [
  { title: "Dashboard", icon: MdDashboard, url: "/" },
  { title: "Explore", icon: FaCompass, url: "/explore" },
  { title: "Messenger", icon: BsChatSquareDots, url: "/messenger" },
  { title: "Check your mood", icon: PiStarFourFill, url: "/check-mood" },
  { title: "Log out", icon: MdLogout, url: "/logout" },
];

export default function Sidebar() {
  const router = useRouter();
  const pathname = usePathname();
  const isActive = (url: string) => pathname === url;
  const handleLogout = () => {
    // clear cookies
    Cookies.remove("token");
    Cookies.remove("user");
    // refresh the page
    window.location.href = "/";
    router.push("/");
  };

  return (
    <aside className="hidden w-64 flex-col justify-between bg-[#28b7be] p-4 md:block">
      <div>
        <div className="mb-8 flex items-center gap-2">
          <Image src={Icon} alt="Icon" className="h-15 w-15" />
          <h1 className="text-2xl font-bold text-white">Mood Bridge</h1>
        </div>
        <nav className="flex flex-col gap-4">
          {items.map((item) => {
            const isLogout = item.title === "Log out";
            const activeClass = isActive(item.url) ? "bg-[#00A6FF]" : "";

            return isLogout ? (
              <button
                key={item.title}
                onClick={handleLogout}
                className={`flex items-center gap-5 rounded-md p-2 font-bold text-white transition-colors hover:bg-[#1a9ea0] ${activeClass}`}
              >
                <item.icon className="text-2xl" />
                <span>{item.title}</span>
              </button>
            ) : (
              <Link
                key={item.title}
                href={item.url}
                className={`flex items-center gap-5 rounded-md p-2 font-bold text-white transition-colors hover:bg-[#1a9ea0] ${activeClass}`}
              >
                <item.icon className="text-2xl" />
                <span>{item.title}</span>
              </Link>
            );
          })}
        </nav>
      </div>
    </aside>
  );
}
