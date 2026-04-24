'use client';

import { useState, useEffect } from 'react';
import axios from 'axios';

// Configure Axios to send cookies automatically to the Go and Node backends
const goAuthApi = axios.create({
  baseURL: 'http://localhost:8080/user',
  withCredentials: true,
});

const nodeApi = axios.create({
  baseURL: 'http://localhost:4000/api',
  withCredentials: true,
});

export default function Home() {
  const [user, setUser] = useState<any>(null);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [protectedData, setProtectedData] = useState<any>(null);
  const [message, setMessage] = useState('');

  // Check if user is already logged in on mount
  useEffect(() => {
    checkUser();
  }, []);

  const checkUser = async () => {
    try {
      const res = await goAuthApi.get('/get');
      setUser(res.data);
    } catch (err) {
      setUser(null);
    }
  };

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append('email', email);
    formData.append('password', password);
    formData.append('signin', '1'); // Auto sign-in after creation

    try {
      const res = await goAuthApi.post('/create', formData);
      setUser(res.data);
      setMessage('Registration successful');
    } catch (err: any) {
      setMessage(`Registration failed: ${err.response?.data?.error || err.message}`);
    }
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append('email', email);
    formData.append('password', password);

    try {
      const res = await goAuthApi.post('/auth', formData);
      setUser(res.data);
      setMessage('Login successful');
    } catch (err: any) {
      setMessage(`Login failed: ${err.response?.data?.error || err.message}`);
    }
  };

  const handleLogout = async () => {
    try {
      await goAuthApi.get('/signout');
      setUser(null);
      setProtectedData(null);
      setMessage('Logged out');
    } catch (err: any) {
      setMessage('Logout failed');
    }
  };

  const fetchProtectedData = async () => {
    try {
      const res = await nodeApi.get('/protected-data');
      setProtectedData(res.data);
      setMessage('Fetched protected data from Node.js backend');
    } catch (err: any) {
      setMessage(`Failed to fetch protected data: ${err.response?.data?.error || err.message}`);
      setProtectedData(null);
    }
  };

  return (
    <main className="min-h-screen p-8 font-[family-name:var(--font-geist-sans)] max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-8 text-center">Auth Demo: Go + SQLite & Node.js + Next.js</h1>

      {message && (
        <div className="mb-4 p-4 bg-blue-100 text-blue-800 rounded">
          {message}
        </div>
      )}

      {!user ? (
        <div className="flex flex-col md:flex-row gap-8">
          <div className="flex-1 p-6 border rounded shadow-sm">
            <h2 className="text-xl font-semibold mb-4">Login</h2>
            <form onSubmit={handleLogin} className="flex flex-col gap-4 text-black">
              <input
                type="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="p-2 border rounded"
                required
              />
              <input
                type="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="p-2 border rounded"
                required
              />
              <button type="submit" className="p-2 bg-blue-600 text-white rounded hover:bg-blue-700">Login</button>
            </form>
          </div>

          <div className="flex-1 p-6 border rounded shadow-sm">
            <h2 className="text-xl font-semibold mb-4">Register</h2>
            <form onSubmit={handleRegister} className="flex flex-col gap-4 text-black">
              <input
                type="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="p-2 border rounded"
                required
              />
              <input
                type="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="p-2 border rounded"
                required
              />
              <button type="submit" className="p-2 bg-green-600 text-white rounded hover:bg-green-700">Register</button>
            </form>
          </div>
        </div>
      ) : (
        <div className="flex flex-col gap-8">
          <div className="p-6 border rounded bg-gray-50 text-black shadow-sm">
            <h2 className="text-xl font-semibold mb-4">Welcome, {user.email}!</h2>
            <p className="mb-4">You are authenticated via the Go Auth Service.</p>
            <button
              onClick={handleLogout}
              className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
            >
              Logout
            </button>
          </div>

          <div className="p-6 border rounded shadow-sm">
            <h2 className="text-xl font-semibold mb-4">Node.js Protected Resource</h2>
            <p className="mb-4 text-gray-600">
              Clicking this button will send a request to the Node.js backend.
              The Node.js backend will extract your session cookie and verify it with the Go Auth Service before returning data.
            </p>
            <button
              onClick={fetchProtectedData}
              className="px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700 mb-4"
            >
              Fetch Protected Data
            </button>

            {protectedData && (
              <div className="mt-4 p-4 bg-gray-900 text-green-400 rounded overflow-x-auto">
                <pre>{JSON.stringify(protectedData, null, 2)}</pre>
              </div>
            )}
          </div>
        </div>
      )}
    </main>
  );
}
