"use client";
import React, { useState, useRef, useEffect} from "react";
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

    try {
      const response = await fetch("http://localhost:8080/api/ai/chat", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ user_id: 1, message: msg }),
      });

      const data = await response.json();

      const aiMsg: ChatMessage = {
        id: Date.now() + 1,
        sender: "ai",
        text: data.response || "Sorry, I couldn't understand that.",
      };
      setMessages((prev) => [...prev, aiMsg]);
    } catch (error) {
      setMessages((prev) => [
        ...prev,
        {
          id: Date.now() + 1,
          sender: "ai",
          text: "Error connecting to AI service.",
        },
      ]);
    } finally {
      setIsTyping(false);
    }
  };

  return (
    <div className="flex flex-col h-screen bg-transparent text-white">
      <header className="p-6 border-b border-[#28b7be] text-center font-bold text-4xl">
        AI Companion
      </header>

<main className="flex-grow overflow-auto p-4 flex flex-col space-y-6">
  {messages.map((msg) => (
<div
  key={msg.id}
  className={`p-3 rounded-3xl break-words whitespace-normal text-xl 
    ${msg.sender === "user"
      ? "bg-[#28b7be] text-white self-end text-left"
      : "bg-gray-700 text-white self-start text-left"
    }
    max-w-[40ch]`}
>
  {msg.text}
</div>


  ))}
  {isTyping && (
    <div className="p-3 rounded-lg bg-gray-700 self-start text-left italic opacity-70 max-w-[70%]">
      AI is typing...
    </div>
  )}
  <div ref={bottomRef} />
</main>


      <ChatInputBox onSend={sendMessage} />
    </div>
  );
}
