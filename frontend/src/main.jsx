import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { AuthProvider } from "react-oidc-context";
import './index.css'
import App from './App.jsx'

const cognitoAuthConfig = {
  authority: import.meta.env.VITE_COGNITO_AUTHORITY,
  client_id: import.meta.env.VITE_COGNITO_CLIENT_ID,
  redirect_uri: window.location.origin, // Dynamically grabs localhost
  response_type: "code",
  scope: "email openid phone",
  // Automatically refresh the access token before it expires
  automaticSilentRenew: true,
  onSigninCallback: (_user) => {
    // Remove the ?code=...&state=... from the URL after successful login
    window.history.replaceState(
      {},
      document.title,
      window.location.pathname
    );
  }
};

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <AuthProvider {...cognitoAuthConfig}>
      <App />
    </AuthProvider>
  </StrictMode>,
)
