import axios from 'axios';
import { fetchAuthSession } from 'aws-amplify/auth';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_GATEWAY_URL,
});

api.interceptors.request.use(async (config) => {
  try {
    const session = await fetchAuthSession();
    if (session.tokens?.accessToken) {
      config.headers.Authorization = `Bearer ${session.tokens.accessToken.toString()}`;
    }
  } catch (err) {
    console.error('Error fetching auth session for api call', err);
  }
  return config;
}, (error) => {
  return Promise.reject(error);
});

api.interceptors.response.use((response) => response, (error) => {
  return Promise.reject(error);
});

export default api;
