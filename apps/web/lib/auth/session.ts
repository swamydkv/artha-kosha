// Session management utilities for authentication

export interface SessionData {
  session_id: string;
  user_id: string;
  welcome_message: string;
}

export const SessionStorage = {
  // Save session data to localStorage
  save(session: SessionData): void {
    localStorage.setItem('session_id', session.session_id);
    localStorage.setItem('user_id', session.user_id);
    localStorage.setItem('welcome_message', session.welcome_message);
  },

  // Get session data from localStorage
  get(): SessionData | null {
    const sessionId = localStorage.getItem('session_id');
    const userId = localStorage.getItem('user_id');
    const welcomeMessage = localStorage.getItem('welcome_message');

    if (!sessionId || !userId) {
      return null;
    }

    return {
      session_id: sessionId,
      user_id: userId,
      welcome_message: welcomeMessage || 'Welcome!',
    };
  },

  // Clear session data from localStorage
  clear(): void {
    localStorage.removeItem('session_id');
    localStorage.removeItem('user_id');
    localStorage.removeItem('welcome_message');
  },

  // Check if user is authenticated
  isAuthenticated(): boolean {
    const session = this.get();
    return session !== null;
  },

  // Get session ID
  getSessionId(): string | null {
    return localStorage.getItem('session_id');
  },

  // Get user ID
  getUserId(): string | null {
    return localStorage.getItem('user_id');
  },
};