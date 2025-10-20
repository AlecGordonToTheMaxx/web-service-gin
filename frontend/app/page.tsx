'use client';

import Link from 'next/link';
import { Music, MessageSquare } from 'lucide-react';

export default function Home() {
  return (
    <div
      className="min-h-screen flex items-center justify-center px-4"
      style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
    >
      <div className="max-w-4xl w-full">
        <div className="text-center mb-12">
          <h1 className="text-6xl font-bold text-white mb-4">
            Album Manager
          </h1>
          <p className="text-xl text-purple-200">
            Manage your music collection with AI assistance
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-8">
          {/* Album Manager Card */}
          <Link href="/albums">
            <div className="bg-white rounded-2xl shadow-2xl p-8 hover:shadow-purple-500/50 transition-all duration-300 transform hover:scale-105 cursor-pointer group">
              <div className="flex flex-col items-center text-center">
                <div className="w-20 h-20 bg-purple-600 rounded-full flex items-center justify-center mb-6 group-hover:bg-purple-700 transition-colors">
                  <Music className="text-white" size={40} />
                </div>
                <h2 className="text-3xl font-bold text-gray-800 mb-4">
                  Album Interface
                </h2>
                <p className="text-gray-600 mb-6">
                  Browse, create, edit, and delete albums with a traditional interface
                </p>
                <div className="text-purple-600 font-semibold group-hover:text-purple-700">
                  Go to Albums →
                </div>
              </div>
            </div>
          </Link>

          {/* Chat Interface Card */}
          <Link href="/chat">
            <div className="bg-white rounded-2xl shadow-2xl p-8 hover:shadow-purple-500/50 transition-all duration-300 transform hover:scale-105 cursor-pointer group">
              <div className="flex flex-col items-center text-center">
                <div className="w-20 h-20 bg-purple-600 rounded-full flex items-center justify-center mb-6 group-hover:bg-purple-700 transition-colors">
                  <MessageSquare className="text-white" size={40} />
                </div>
                <h2 className="text-3xl font-bold text-gray-800 mb-4">
                  Chat Assistant
                </h2>
                <p className="text-gray-600 mb-6">
                  Manage albums using natural language with AI-powered chat
                </p>
                <div className="text-purple-600 font-semibold group-hover:text-purple-700">
                  Go to Chat →
                </div>
              </div>
            </div>
          </Link>
        </div>

        <div className="mt-12 text-center">
          <p className="text-purple-200 text-sm">
            Choose your preferred way to manage your music collection
          </p>
        </div>
      </div>
    </div>
  );
}
