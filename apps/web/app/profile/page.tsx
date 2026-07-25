'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function ProfilePage() {
  const router = useRouter();
  const [deleteInput, setDeleteInput] = useState('');
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState('');
  const [sessionId, setSessionId] = useState('');
  const [userId, setUserId] = useState('');

  useEffect(() => {
    // Check if user is authenticated
    const sid = localStorage.getItem('session_id');
    const uid = localStorage.getItem('user_id');

    if (!sid || !uid) {
      router.push('/login');
      return;
    }
    
    setSessionId(sid);
    setUserId(uid);
  }, [router]);

  const handleDeleteAccount = async () => {
    if (deleteInput !== 'DELETE') {
      setError('You must type DELETE to confirm.');
      return;
    }

    setIsDeleting(true);
    setError('');

    try {
      const response = await fetch('http://localhost:8080/user/account', {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'X-Session-ID': sessionId,
          'X-User-ID': userId,
        },
        body: JSON.stringify({ confirmation: deleteInput })
      });

      if (response.ok) {
        // Clear session data
        localStorage.removeItem('session_id');
        localStorage.removeItem('user_id');
        localStorage.removeItem('welcome_message');
        
        router.push('/login?message=account_deleted');
      } else {
        const data = await response.json().catch(() => null);
        setError(data?.error || 'Failed to delete account. Please try again.');
      }
    } catch (err) {
      console.error('Delete failed:', err);
      setError('An error occurred. Please check your connection and try again.');
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center space-x-4">
              <h1 className="text-xl font-bold text-gray-900">ArthaKosha</h1>
              <Link href="/home" className="text-gray-600 hover:text-gray-900">Home</Link>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0 max-w-2xl">
          
          <div className="bg-white rounded-lg shadow p-6 mb-8">
            <h2 className="text-2xl font-bold text-gray-900 mb-4">Profile Settings</h2>
            <p className="text-gray-600 mb-6">Manage your account settings here.</p>
          </div>

          {/* Account Deletion */}
          <div className="bg-white rounded-lg shadow border border-red-200">
            <div className="px-6 py-5 border-b border-gray-200 bg-red-50 rounded-t-lg">
              <h3 className="text-lg leading-6 font-medium text-red-800">
                Account Deletion
              </h3>
            </div>
            <div className="p-6">
              <h4 className="text-md font-bold text-gray-900 mb-2">Delete Account</h4>
              <p className="text-sm text-gray-500 mb-4">
                Once you delete your account, your personal data will be anonymized and archived according to GDPR regulations. This action cannot be undone.
              </p>
              
              {error && (
                <div className="mb-4 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative">
                  <span className="block sm:inline">{error}</span>
                </div>
              )}

              <div className="mt-4">
                <label htmlFor="delete-confirm" className="block text-sm font-medium text-gray-700 mb-2">
                  Please type <strong>DELETE</strong> to confirm.
                </label>
                <div className="flex gap-4">
                  <input
                    type="text"
                    id="delete-confirm"
                    value={deleteInput}
                    onChange={(e) => setDeleteInput(e.target.value)}
                    className="shadow-sm focus:ring-red-500 focus:border-red-500 block w-full sm:text-sm border-gray-300 rounded-md p-2 border"
                    placeholder="DELETE"
                    disabled={isDeleting}
                  />
                  <button
                    onClick={handleDeleteAccount}
                    disabled={deleteInput !== 'DELETE' || isDeleting}
                    className="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap font-medium"
                  >
                    {isDeleting ? 'Deleting...' : 'Delete Account'}
                  </button>
                </div>
              </div>
            </div>
          </div>

        </div>
      </main>
    </div>
  );
}
