import React from 'react';

interface ChatBubbleProps {
  sender: 'user' | 'bot';
  message: string;
}

const ChatBubble: React.FC<ChatBubbleProps> = ({ sender, message }) => {
  const isUser = sender === 'user';
  return (
    <div className={`flex ${isUser ? 'justify-end' : 'justify-start'} my-2`}>
      <div
        className={`
          inline-block
          p-3 rounded-full
          ${isUser ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-900'}
          max-w-[80%]
          break-words
        `}
      >
        {message}
      </div>
    </div>
  );
};

export default ChatBubble;
