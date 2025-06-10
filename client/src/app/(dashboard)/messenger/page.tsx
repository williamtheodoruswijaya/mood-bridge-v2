"use client";

import axios from "axios";
import Cookies from "js-cookie";
import { useEffect, useRef, useState } from "react";
import {
  type MessageInterface,
  type FriendInterface,
  type FriendResponse,
  type User,
} from "~/types/types";
import profile_1 from "~/assets/profile/profile-picture-1.png";
import profile_2 from "~/assets/profile/profile-picture-2.png";
import profile_3 from "~/assets/profile/profile-picture-3.png";
import profile_4 from "~/assets/profile/profile-picture-4.png";
import profile_5 from "~/assets/profile/profile-picture-5.png";
import { DecodeUserFromToken } from "~/utils/utils";
import Image from "next/image";

export default function Page() {
  const token = Cookies.get("token");
  const [user, setUser] = useState<User>({
    id: 0,
    username: "",
    email: "",
    fullname: "",
    createdAt: "",
  });
  const [friends, setFriends] = useState<FriendInterface[]>([]);
  const [messages, setMessages] = useState<MessageInterface[]>([]);
  const [currentChatFriend, setCurrentChatFriend] =
    useState<FriendInterface | null>(null);
  const [newMessage, setNewMessage] = useState("");
  const [isLoadingHistory, setIsLoadingHistory] = useState(false); // eslint-disable-line
  const ws = useRef<WebSocket | null>(null);
  const chatContainerRef = useRef<HTMLDivElement | null>(null); // eslint-disable-line
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

  // step 1: Fetch user data yang sedang login
  useEffect(() => {
    if (token) {
      const decodedUser = DecodeUserFromToken(token);
      if (decodedUser?.user && decodedUser) {
        setUser({
          id: decodedUser.user.id,
          username: decodedUser.user.username,
          email: decodedUser.user.email,
          fullname: decodedUser.user.fullname,
          createdAt: decodedUser.user.created_at,
        });
      }
    }
  }, [token]);

  // step 2: Fetch daftar teman milik user
  useEffect(() => {
    const fetchFriends = async () => {
      if (!token || !user.id) return;
      if (token && user.id) {
        try {
          const response = await axios.get<FriendResponse>(
            `${process.env.NEXT_PUBLIC_API_URL}/api/friend/all/${user.id}`,
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
          // TODO: Handle error appropriately
          console.error("Error fetching friends:", error);
        }
      }
    };
    fetchFriends().catch((error) => {
      console.error("Error fetching friends:", error);
    });
  }, [token, user.id]);

  // step 3: Nyalakan koneksi WebSocket
  useEffect(() => {
    if (!token || !user.id) return;
    if (!ws.current || ws.current.readyState === WebSocket.CLOSED) {
      ws.current = new WebSocket(
        `${process.env.NEXT_PUBLIC_WEB_SOCKET_URL}/api/chat/ws?id=${user.id}`,
      );
      ws.current.onopen = () => {
        console.log("WebSocket connection established");
      };
      ws.current.onclose = (event) => {
        console.log("WebSocket connection closed", event.code, event.reason);
      };
      ws.current.onerror = (error) => {
        console.error("WebSocket error:", error);
      };
    }

    ws.current.onmessage = (event: MessageEvent<string>) => {
      try {
        // step 6: parse incoming message as JSON
        const receivedMsg = JSON.parse(event.data) as MessageInterface;
        console.log("Received message:", receivedMsg);

        // step 7: ambil pesan dari payload
        if (
          receivedMsg.type === "new_private_message" ||
          receivedMsg.type === "offline_message"
        ) {
          const chatMessage = receivedMsg.payload;
          if (
            currentChatFriend &&
            (chatMessage.senderid === currentChatFriend.userid ||
              chatMessage.senderid === user.id)
          ) {
            const newMessage: MessageInterface = {
              type: receivedMsg.type,
              payload: {
                id: chatMessage.id,
                senderid: chatMessage.senderid,
                recipientid: chatMessage.recipientid,
                content: chatMessage.content,
                timestamp: chatMessage.timestamp,
                status: chatMessage.status,
              },
            };
            setMessages((prevMessages) => [...prevMessages, newMessage]);
          } else {
            // Jika pesan bukan untuk chat saat ini, bisa simpan atau tampilkan notifikasi
            // TODO: Handle message not for current chat by showing a notification or updating UI
            console.log("Message not for current chat:", receivedMsg);
          }
        }
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
      }
    };
    return () => {
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        ws.current.close();
        console.log("WebSocket connection closed on unmount");
      }
    };
  }, [token, user.id, currentChatFriend]);

  // Fetch history messages when a friend is selected
  const fetchHistoryMessages = async (friend: FriendInterface) => {
    setCurrentChatFriend(friend);

    // Reset unread messages count for the selected friend
    setFriends((prevFriends) =>
      prevFriends.map((f) =>
        f.id === friend.id ? { ...f, unreadMessages: 0 } : f,
      ),
    );

    if (!token || !user.id) return;
    try {
      const response = await axios.get<MessageInterface[]>(
        `${process.env.NEXT_PUBLIC_WEB_SOCKET_URL}/api/chat/history?with_user_id=${friend.userid}&limit=50&offset=0`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        },
      );
      if (response.status === 200) {
        const historyMessages = response.data
          .map((msg) => ({
            type: msg.type,
            payload: {
              id: msg.payload.id,
              senderid: msg.payload.senderid,
              recipientid: msg.payload.recipientid,
              content: msg.payload.content,
              timestamp: msg.payload.timestamp,
              status: msg.payload.status,
            },
          }))
          .sort(
            (a, b) =>
              new Date(a.payload.timestamp).getTime() -
              new Date(b.payload.timestamp).getTime(),
          );
        setMessages(historyMessages);
      }
    } catch (error) {
      console.error("Error fetching history messages:", error);
      setMessages([]); // Clear messages on error
    }
  };

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (
      !newMessage.trim() ||
      !currentChatFriend ||
      !ws.current ||
      ws.current.readyState !== WebSocket.OPEN
    )
      return;
    const messagePayload = {
      recipientid: currentChatFriend.userid,
      content: newMessage.trim(),
    };
    try {
      ws.current.send(JSON.stringify(messagePayload));
      setNewMessage("");
    } catch (error) {
      console.error("Error sending message:", error);
      alert("Failed to send message. Please try again.");
    }
  };

  return (
    <main className="grid h-full w-full grid-cols-[1fr_300px] text-black">
      <section className="flex flex-col justify-between rounded-lg bg-gradient-to-tl from-white to-cyan-100 p-4">
        {currentChatFriend ? (
          <>
            <div className="mb-2 pb-2">
              <h3 className="text-2xl font-semibold">
                Chat with @{currentChatFriend.user.username}{" "}
              </h3>
            </div>
            <div className="flex-1 space-y-4 overflow-y-auto rounded border p-4 shadow-inner">
              {messages.map((msg) => (
                <div
                  key={msg.payload.id}
                  className={`flex ${msg.payload.senderid === user.id ? "justify-end" : "justify-start"}`}
                >
                  <div
                    className={`max-w-xs rounded-lg px-4 py-2 ${
                      msg.payload.senderid === user.id
                        ? "bg-blue-500 text-white"
                        : "bg-gray-200 text-black"
                    }`}
                  >
                    {msg.payload.content}
                  </div>
                </div>
              ))}
            </div>
            <form className="mt-4 flex gap-2" onSubmit={handleSendMessage}>
              <input
                type="text"
                placeholder="Type a message..."
                className="flex-1 rounded border px-3 py-2"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
              />
              <button
                type="submit"
                className="rounded bg-blue-500 px-4 py-2 text-white"
              >
                Send
              </button>
            </form>
          </>
        ) : (
          <div className="flex h-full items-center justify-center text-gray-500">
            Select a friend to start chatting.
          </div>
        )}
      </section>

      {/* Friends List - Tampilkan Unread Count */}
      <aside className="ml-4 overflow-y-auto rounded-lg bg-[#28b7be] p-4">
        <h2 className="mb-4 text-2xl font-bold text-white">Friends</h2>
        <ul className="space-y-2">
          {friends.map((friend) => (
            <button
              key={friend.id}
              className={`flex w-full cursor-pointer items-center rounded bg-white p-2 shadow hover:bg-sky-200 ${
                currentChatFriend?.id === friend.id ? "bg-sky-300" : ""
              }`}
              onClick={() => fetchHistoryMessages(friend)}
            >
              <Image
                src={getProfilePicture(friend.user.userid.toString())!.src}
                width={40}
                height={40}
                alt="Profile-Picture"
                className="mr-2 rounded-full object-cover"
              />
              <div className="text-sm font-medium">{friend.user.fullname}</div>
              <div className="ml-1 text-xs font-extralight">
                @{friend.user.username}
              </div>
            </button>
          ))}
        </ul>
      </aside>
    </main>
  );
}
