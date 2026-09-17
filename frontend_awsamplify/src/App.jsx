import { useState, useEffect } from "react";
import { signInWithRedirect, signOut, getCurrentUser, fetchAuthSession } from 'aws-amplify/auth';
import { Hub } from 'aws-amplify/utils';
import Dashboard from "./pages/Dashboard";

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
      const currentUser = await getCurrentUser();
      
      // Extract user claims from the ID token
      try {
        const session = await fetchAuthSession();
        const idTokenPayload = session.tokens?.idToken?.payload;
        
        currentUser.attributes = {
          email: idTokenPayload?.email?.toString(),
          name: idTokenPayload?.name?.toString()
        };
      } catch (attrErr) {
        console.error('Could not fetch session tokens', attrErr);
      }
      
      setUser(currentUser);
      setError(null);
    } catch (e) {
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
          email: user.attributes?.email || user.attributes?.name || ''
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
