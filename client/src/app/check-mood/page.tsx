"use client";
import React, { useState, useRef, useEffect } from "react";
import ChatInputBox from "~/components/chat-box";

type ChatMessage = {
  id: number;
  sender: "user" | "ai";
  text: string;
};

export default function AIChatPage() {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [isTyping, setIsTyping] = useState(false);
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" }); // 👈 Auto-scroll
  }, [messages, isTyping]);

  const sendMessage = async (msg: string) => {
    const userMsg: ChatMessage = {
      id: Date.now(),
      sender: "user",
      text: msg,
    };
    setMessages((prev) => [...prev, userMsg]);
    setIsTyping(true);

    // try {
    //   const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/ai/chat`, {
    //     method: "POST",
    //     headers: { "Content-Type": "application/json" },
    //     body: JSON.stringify({ user_id: 1, message: msg }),
    //   });

    //   // const data = await response.json();

    //   // const aiMsg: ChatMessage = {
    //   //   id: Date.now() + 1,
    //   //   sender: "ai",
    //   //   text: data.response ?? "Sorry, I couldn't understand that.",
    //   // };
    //   // setMessages((prev) => [...prev, aiMsg]);
    // } catch (error) {
    //   setMessages((prev) => [
    //     ...prev,
    //     {
    //       id: Date.now() + 1,
    //       sender: "ai",
    //       text: "Error connecting to AI service.",
    //     },
    //   ]);
    // } finally {
    //   setIsTyping(false);
    // }
  };

  return (
    <div className="flex h-screen flex-col bg-transparent text-white">
      <header className="border-b border-[#28b7be] p-6 text-center text-4xl font-bold">
        AI Companion
      </header>

      <main className="flex flex-grow flex-col space-y-6 overflow-auto p-4">
        {messages.map((msg) => (
          <div
            key={msg.id}
            className={`rounded-3xl p-3 text-xl break-words whitespace-normal ${
              msg.sender === "user"
                ? "self-end bg-[#28b7be] text-left text-white"
                : "self-start bg-gray-700 text-left text-white"
            } max-w-[40ch]`}
          >
            {msg.text}
          </div>
        ))}
        {isTyping && (
          <div className="max-w-[70%] self-start rounded-lg bg-gray-700 p-3 text-left italic opacity-70">
            AI is typing...
          </div>
        )}
        <div ref={bottomRef} />
      </main>

      <ChatInputBox onSend={sendMessage} />
    </div>
  );
}
