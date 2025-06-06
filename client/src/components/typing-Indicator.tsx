import React from 'react';

interface TypingIndicatorProps {
  isTyping: boolean;
}

const TypingIndicator: React.FC<TypingIndicatorProps> = ({ isTyping }) => {
  if (!isTyping) return null;

  return (
    <div className="italic text-gray-500 px-3 py-1">AI is typing...</div>
  );
};

export default TypingIndicator;
