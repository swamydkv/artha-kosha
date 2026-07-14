'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export default function HomePage() {
  const router = useRouter();
  const [welcomeMessage, setWelcomeMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    // Check if user is authenticated
    const sessionId = localStorage.getItem('session_id');
    const userId = localStorage.getItem('user_id');
    const welcome = localStorage.getItem('welcome_message');

    if (!sessionId || !userId) {
      router.push('/login');
      return;
    }

    setWelcomeMessage(welcome || 'Welcome!');
  }, [router]);

  const handleLogout = async () => {
    setIsLoading(true);

    try {
      const sessionId = localStorage.getItem('session_id');
      
      const response = await fetch('http://localhost:8080/logout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Session-ID': sessionId || '',
        },
      });

      if (response.ok) {
        // Clear session data
        localStorage.removeItem('session_id');
        localStorage.removeItem('user_id');
        localStorage.removeItem('welcome_message');
        
        router.push('/login');
      }
    } catch (error) {
      console.error('Logout failed:', error);
      // Even if logout fails, clear local session
      localStorage.removeItem('session_id');
      localStorage.removeItem('user_id');
      localStorage.removeItem('welcome_message');
      router.push('/login');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-bold text-gray-900">ArthaKosha</h1>
            </div>
            <div className="flex items-center">
              <button
                onClick={handleLogout}
                disabled={isLoading}
                className="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isLoading ? 'Logging out...' : 'Logout'}
              </button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-2xl font-bold text-gray-900 mb-4">
              {welcomeMessage}
            </h2>
            <p className="text-gray-600">
              Welcome to ArthaKosha - your personal finance management platform.
            </p>
            <div className="mt-6">
              <p className="text-sm text-gray-500">
                This is the authenticated home view. Financial features will be added in future releases.
              </p>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}