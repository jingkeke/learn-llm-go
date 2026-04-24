const express = require('express');
const cookieParser = require('cookie-parser');
const cors = require('cors');
const axios = require('axios');

const app = express();
const PORT = process.env.PORT || 4000;
const GO_AUTH_SERVICE_URL = 'http://localhost:8080/user';

app.use(cookieParser());
app.use(cors({
  origin: 'http://localhost:3000', // Allow frontend
  credentials: true
}));

// Middleware to check if user is authenticated by calling the Go Auth Service
const requireAuth = async (req, res, next) => {
  const sessionCookie = req.cookies.session;

  if (!sessionCookie) {
    return res.status(401).json({ error: 'No session cookie found' });
  }

  try {
    // We forward the session cookie to the Go Auth Service to verify it
    const response = await axios.get(`${GO_AUTH_SERVICE_URL}/get`, {
      headers: {
        Cookie: `session=${sessionCookie}`
      }
    });

    // If successful, the Go service returns the user info
    req.user = response.data;
    next();
  } catch (error) {
    // Go Auth service returns 401 if not authenticated
    if (error.response && error.response.status === 401) {
      return res.status(401).json({ error: 'Invalid or expired session' });
    }
    console.error('Error verifying session:', error.message);
    res.status(500).json({ error: 'Internal Server Error' });
  }
};

// Protected endpoint
app.get('/api/protected-data', requireAuth, (req, res) => {
  res.json({
    message: 'This is protected data from the Node.js backend!',
    user: req.user,
    secretData: [1, 2, 3, 4, 5]
  });
});

app.listen(PORT, () => {
  console.log(`Node.js Backend Service running on http://localhost:${PORT}`);
});
