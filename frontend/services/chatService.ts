const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface ChatMessage {
  role: string;
  content: string;
}

export interface ChatRequest {
  messages: ChatMessage[];
}

export interface ChatResponse {
  message: string;
  tool_calls?: any[];
  tool_results?: any[];
}

class ChatServiceError extends Error {
  constructor(
    public status: number,
    message: string
  ) {
    super(message);
    this.name = 'ChatServiceError';
  }
}

export const chatService = {
  async sendMessage(messages: ChatMessage[]): Promise<ChatResponse> {
    try {
      const response = await fetch(`${API_URL}/chat`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ messages }),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new ChatServiceError(
          response.status,
          errorData.error || `HTTP error! status: ${response.status}`
        );
      }

      return await response.json();
    } catch (error) {
      if (error instanceof ChatServiceError) {
        throw error;
      }
      throw new ChatServiceError(
        0,
        error instanceof Error ? error.message : 'An unknown error occurred'
      );
    }
  },
};
