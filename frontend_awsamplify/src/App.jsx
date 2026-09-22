import { useState, useEffect } from "react";
import { signInWithRedirect, signOut, getCurrentUser, fetchAuthSession } from 'aws-amplify/auth';
import { Hub } from 'aws-amplify/utils';
import Dashboard from "./pages/Dashboard";
import api from './utils/api';

function App() {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    checkUser();

    const unsubscribe = Hub.listen('auth', ({ payload }) => {
      switch (payload.event) {
        case 'signInWithRedirect':
          checkUser();
          break;
        case 'signInWithRedirect_failure':
          setError('An error has occurred during the OAuth flow.');
          break;
        case 'customOAuthState':
          break;
        case 'signedOut':
          setUser(null);
          break;
      }
    });

    return unsubscribe;
  }, []);

  async function checkUser() {
    try {
      setIsLoading(true);
      const currentUser = await getCurrentUser(); // Just checks if logged in via Amplify
      
      // Fetch session specifically for the ID token required by /me
      const session = await fetchAuthSession();
      const idToken = session.tokens?.idToken?.toString();

      // Call backend to authenticate and provision the user
      const response = await api.get('/me', {
        headers: idToken ? { 'X-Id-Token': idToken } : {}
      });
      
      const userData = response.data?.data;
      if (!userData) {
        throw new Error('No user data returned from /me');
      }

      setUser(userData);
      setError(null);
    } catch (e) {
      console.error('User check failed', e);
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  }

  const handleSignIn = async () => {
    try {
      await signInWithRedirect();
    } catch (err) {
      console.error(err);
    }
  };

  const handleSignOut = async () => {
    try {
      await signOut();
    } catch (err) {
      console.error(err);
    }
  };

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <p className="text-gray-500">Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <p className="text-red-500">Encountering error... {error}</p>
      </div>
    );
  }

  if (user) {
    // Note: Dashboard expects `auth` prop. Let's mock the structure it needs
    // Dashboard currently uses: `auth.user?.profile?.email`
    // With Amplify, we will get the email from the API payload if needed, or pass the user object.
    const mockAuth = {
      user: {
        profile: {
          email: user.name || user.email || "",
        }
      }
    };
    
    return (
      <Dashboard auth={mockAuth} onSignOut={handleSignOut} />
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <div className="w-full max-w-md p-8 space-y-8 bg-white shadow rounded-lg text-center">
        <h2 className="text-3xl font-extrabold text-gray-900">Welcome to EMS</h2>
        <p className="mt-2 text-sm text-gray-600">Employee Management System</p>
        <button 
          onClick={handleSignIn}
          className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none cursor-pointer"
        >
          Sign in
        </button>
      </div>
    </div>
  );
}

export default App;
