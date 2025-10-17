import React from 'react';
import { Bot, User } from 'lucide-react';

export interface Message {
  id: number;
  content: string;
  role: string;
  timestamp: Date;
}

interface MessageProps {
  message: Message;
  isStreaming?: boolean;
}

const formatDate = (date: Date) => {
  return new Intl.DateTimeFormat('en-US', {
    hour: 'numeric',
    minute: 'numeric',
    hour12: true
  }).format(date);
};

export const MessageComponent: React.FC<MessageProps> = ({
  message,
  isStreaming = false
}) => {
  const isUser = message.role === 'user';

  return (
    <div className={`flex items-start gap-3 p-4 ${isUser ? 'flex-row-reverse' : ''}`}>
      <div
        className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center ${
          isUser ? 'bg-purple-600 text-white' : 'bg-gray-200 text-gray-600'
        }`}
      >
        {isUser ? <User size={16} /> : <Bot size={16} />}
      </div>

      <div className={`max-w-[70%] ${isUser ? 'text-right' : ''}`}>
        <div
          className={`inline-block p-3 rounded-lg whitespace-pre-wrap ${
            isUser
              ? 'bg-purple-600 text-white rounded-br-sm'
              : 'bg-gray-100 text-gray-800 rounded-bl-sm'
          }`}
        >
          {message.content}
          {isStreaming && (
            <span className="inline-block w-2 h-5 bg-current ml-1 animate-pulse" />
          )}
        </div>

        <div className="text-xs text-gray-500 mt-1">
          {formatDate(message.timestamp)}
        </div>
      </div>
    </div>
  );
};
