import { useAuth } from "react-oidc-context";
import Dashboard from "./pages/Dashboard";
import { setApiToken } from "./utils/api";

function App() {
  const auth = useAuth();

  const signOutRedirect = () => {
    const clientId = import.meta.env.VITE_COGNITO_CLIENT_ID;
    const logoutUri = window.location.origin; // or whatever URI you prefer
    const cognitoDomain = import.meta.env.VITE_COGNITO_DOMAIN;
    window.location.href = `${cognitoDomain}/logout?client_id=${clientId}&logout_uri=${encodeURIComponent(logoutUri)}`;
  };

  if (auth.isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <p className="text-gray-500">Loading...</p>
      </div>
    );
  }

  if (auth.error) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <p className="text-red-500">Encountering error... {auth.error.message}</p>
      </div>
    );
  }

  if (auth.isAuthenticated) {
    // Set the token for future API calls
    if (auth.user?.access_token) {
      setApiToken(auth.user.access_token);
    }
    
    // Optionally render Dashboard if authenticated
    return (
      <Dashboard auth={auth} onSignOut={() => {
        auth.removeUser();
        signOutRedirect();
      }} />
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <div className="w-full max-w-md p-8 space-y-8 bg-white shadow rounded-lg text-center">
        <h2 className="text-3xl font-extrabold text-gray-900">Welcome to EMS</h2>
        <p className="mt-2 text-sm text-gray-600">Employee Management System</p>
        <button 
          onClick={() => auth.signinRedirect()}
          className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none"
        >
          Sign in
        </button>
      </div>
    </div>
  );
}

export default App;
