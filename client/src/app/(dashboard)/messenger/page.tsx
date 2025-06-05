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
import { DecodeUserFromToken, TimeAgo } from "~/utils/utils";

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
  const ws = useRef<WebSocket | null>(null);

  // Fetch user data and friends list
  useEffect(() => {
    const fetchUserData = async () => {
      if (token) {
        const user = DecodeUserFromToken(token);
        if (user) {
          setUser({
            id: user.user.id,
            username: user.user.username,
            email: user.user.email,
            fullname: user.user.fullname,
            createdAt: user.user.created_at,
          });
        }
      }
    };

    fetchUserData().catch((error) => {
      console.error("Error fetching user data:", error);
    });
  }, [token]);

  useEffect(() => {
    const fetchFriends = async () => {
      if (token) {
        try {
          const response = await axios.get<FriendResponse>(
            `http://localhost:8080/api/friend/all/${user.id}`,
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

  // Establish WebSocket connection
  useEffect(() => {
    if (!token) return;

    // step 1: pastiin cuman ada satu koneksi WebSocket yang terbentuk
    if (ws.current && ws.current.readyState !== WebSocket.CLOSED) {
      console.log("WebSocket already connected");
    } else {
      // step 2: buat koneksi WebSocket baru
      ws.current = new WebSocket(
        `ws://localhost:8080/api/chat/ws?id=${user.id}`,
      );

      // step 3: make sure WebSocket sudah open
      ws.current.onopen = () => {
        console.log("WebSocket connection established");

        // step 4: kirim authentication message
        ws.current?.send(
          JSON.stringify({
            type: "Authorization",
            token: token,
          }),
        );
      };

      // step 5: handle incoming messages
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

      // step 8: handle close event
      ws.current.onclose = (event) => {
        console.log("WebSocket connection closed:", event);
        ws.current = null; // Reset WebSocket reference
      };

      // step 9: handle WebSocket errors
      ws.current.onerror = (error) => {
        console.error("WebSocket error:", error);
      };
    }

    // Cleanup function to close WebSocket connection on unmount
    return () => {
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        ws.current.close();
        console.log("WebSocket connection closed on unmount");
      }
    };
  }, [token, currentChatFriend, user.id]);

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
        `ws://localhost:8080/api/chat/history?with_user_id=${friend.userid}&limit=50&offset=0`,
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

  console.log("Current user:", user);
  console.log("Friends list:", friends);

  return (
    <main className="grid h-screen w-full grid-cols-[1fr_300px] bg-white text-black">
      {/* Chat Box */}
      <section className="flex flex-col justify-between bg-gradient-to-br from-blue-50 to-white p-4">
        {/* ... (bagian chat box tidak berubah signifikan, kecuali tampilan pesan mungkin) ... */}
        {currentChatFriend ? (
          <>
            <div className="mb-2 border-b pb-2">
              <h3 className="text-xl font-semibold">
                Chat with {currentChatFriend.user.username}{" "}
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
      <aside className="overflow-y-auto border-l bg-sky-100 p-4">
        <h2 className="mb-4 text-xl font-bold">Friends</h2>
        <ul className="space-y-2">
          {friends.map((friend) => (
            <li
              key={friend.id}
              className={`flex cursor-pointer items-center justify-between rounded bg-white p-2 shadow hover:bg-sky-200 ${
                currentChatFriend?.id === friend.id ? "bg-sky-300" : ""
              }`}
              onClick={() => fetchHistoryMessages(friend)}
            >
              <span>{friend.user.username}</span>
              {/* {unreadMessage > 0 && (
                <span className="ml-2 rounded-full bg-red-600 px-2 py-0.5 text-xs font-semibold text-white">
                  {friend.unreadCount}
                </span>
              )} */}
            </li>
          ))}
        </ul>
      </aside>
    </main>
  );
}
