'use client';

import React, { useState, useRef, useEffect, useCallback } from 'react';
import { Bot } from 'lucide-react';
import { MessageComponent, type Message } from '../../components/chat/Message';
import { ChatInput } from '../../components/chat/ChatInput';
import { LoadingSpinner } from '../../components/common/LoadingSpinner';
import { chatService, type ChatMessage } from '../../services/chatService';

export default function ChatPage() {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: 1,
      content: 'Hello! I\'m your album management assistant. I can help you view, create, update, and delete albums. What would you like to do?',
      role: 'assistant',
      timestamp: new Date(),
    },
  ]);
  const [streamingMessage, setStreamingMessage] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);
  const [isStreaming, setIsStreaming] = useState(false);

  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, streamingMessage]);

  const sendMessage = useCallback(async (content: string) => {
    // Add user message
    const userMessage: Message = {
      id: Date.now(),
      content,
      role: 'user',
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setIsLoading(true);

    try {
      // Convert messages to API format
      const apiMessages: ChatMessage[] = messages
        .filter((m) => m.role !== 'system')
        .map((m) => ({
          role: m.role,
          content: m.content,
        }));

      // Add current user message
      apiMessages.push({
        role: 'user',
        content,
      });

      // Call API
      const response = await chatService.sendMessage(apiMessages);

      // Add assistant response
      const assistantMessage: Message = {
        id: Date.now() + 1,
        content: response.message,
        role: 'assistant',
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, assistantMessage]);
    } catch (error) {
      console.error('Error sending message:', error);

      // Add error message
      const errorMessage: Message = {
        id: Date.now() + 1,
        content: 'Sorry, I encountered an error processing your request. Please try again.',
        role: 'assistant',
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  }, [messages]);

  const quickActions = [
    'Show me all albums',
    'Create a new album',
    'What albums do you have?',
    'Delete an album',
  ];

  return (
    <div
      className="min-h-screen py-12 px-4 sm:px-6 lg:px-8"
      style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
    >
      <div className="max-w-4xl mx-auto">
        <div className="flex flex-col h-[600px] border border-gray-200 rounded-lg shadow-lg bg-white">
          {/* Header */}
          <div className="flex items-center justify-between p-4 border-b border-gray-200 bg-gray-50 rounded-t-lg">
            <div className="flex items-center gap-3">
              <Bot className="text-purple-600" size={24} />
              <div>
                <h3 className="font-semibold text-gray-900">Album Assistant</h3>
                <p className="text-sm text-gray-600">
                  Ask me anything about your album collection
                </p>
              </div>
            </div>
          </div>

          {/* Messages */}
          <div className="flex-1 overflow-y-auto">
            {messages.map((message) => (
              <MessageComponent key={message.id} message={message} />
            ))}

            {/* Streaming message */}
            {isStreaming && streamingMessage && (
              <MessageComponent
                message={{
                  id: 0,
                  content: streamingMessage,
                  role: 'assistant',
                  timestamp: new Date(),
                }}
                isStreaming={true}
              />
            )}

            {/* Loading indicator */}
            {isLoading && (
              <div className="flex items-center gap-3 p-4">
                <div className="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center">
                  <Bot size={16} className="text-gray-600" />
                </div>
                <div className="flex items-center gap-2 bg-gray-100 rounded-lg p-3">
                  <LoadingSpinner size={16} className="text-gray-600" />
                  <span className="text-gray-600">Thinking...</span>
                </div>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>

          {/* Quick Actions */}
          {messages.length <= 1 && (
            <div className="p-4 border-t border-gray-100 bg-gray-50">
              <p className="text-sm text-gray-600 mb-2">Try asking:</p>
              <div className="flex flex-wrap gap-2">
                {quickActions.map((action, index) => (
                  <button
                    key={index}
                    onClick={() => sendMessage(action)}
                    disabled={isLoading || isStreaming}
                    className="text-sm px-3 py-1 bg-purple-100 text-purple-700 rounded-full hover:bg-purple-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {action}
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Input */}
          <ChatInput
            onSendMessage={sendMessage}
            isLoading={isLoading}
            isStreaming={isStreaming}
          />
        </div>
      </div>
    </div>
  );
}
