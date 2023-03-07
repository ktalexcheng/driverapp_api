import './database/connection.js';
import express from 'express';
import RideRouter from './route/rides.js';
import ProfileRouter from './route/profile.js';
import AuthRouter from './route/auth.js';

const app = express();

// Middleware: what happens after server gets request but before passing to route
// Use middleware to accept JSON as body
app.use(express.json({ limit: '100mb' }));

// Landing page
app.get('/', (req, res) => {
    res.send('Welcome to the Trailbrake API!');
});

// Use RideRouter whenever URI ends in /rides
app.use('/rides', RideRouter);

// Use ProfileRouter whenever URI ends in /profile
app.use('/profile', ProfileRouter);

// Use AuthRouter whenever URI ends in /auth
app.use('/auth', AuthRouter);

// Start server
app.listen(8080, () => console.log('Server started'));