import React, { useState, type KeyboardEvent, type ChangeEvent } from "react";
import { FiSend } from "react-icons/fi";

interface ChatInputBoxProps {
  onSend: (message: string) => void;
}

const ChatInputBox: React.FC<ChatInputBoxProps> = ({ onSend }) => {
  const [message, setMessage] = useState<string>("");

  const hasText = message.trim().length > 0;

  const handleSend = () => {
    if (!hasText) return;
    onSend(message.trim());
    setMessage("");
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const onChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
    setMessage(e.target.value);
  };

  return (
    <div className="fixed bottom-4 left-0 right-0 flex justify-center px-4 z-50">
      <div
        className={`
          flex items-center rounded-3xl px-4 py-2 shadow-lg transition-all duration-300 ease-in-out
          w-full max-w-[400px]
          hover:max-w-[600px]
          ${hasText ? "max-w-[600px]" : ""}
          bg-[#28b7be]
          translate-x-27
        `}
      >
        <textarea
          rows={1}
          value={message}
          onChange={onChange}
          onKeyDown={handleKeyDown}
          placeholder="Type your message..."
          className="flex-grow resize-none bg-transparent px-2 py-2 text-white placeholder-white focus:outline-none"
        />
        <button
          onClick={handleSend}
          disabled={!hasText}
          className={`
            ml-2 flex h-10 w-10 items-center justify-center rounded-full
            transition-colors duration-200
            ${hasText ? "bg-white" : "bg-[#1a9ea0] hover:bg-[#1692a5]"}
            disabled:opacity-50
          `}
          aria-label="Send message"
        >
          <FiSend className={`${hasText ? "text-black" : "text-white"}`} />
        </button>
      </div>
    </div>
  );

};

export default ChatInputBox;
