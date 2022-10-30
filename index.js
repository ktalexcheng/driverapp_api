import './database/connection.js';
import express from 'express';
import RideRouter from './route/rides.js';
import ProfileRouter from './route/profile.js';

const app = express();

// Middleware: what happens after server gets request but before passing to route
// Use middleware to accept JSON as body
app.use(express.json({ limit: '10mb' }));

// Landing page
app.get('/', (req, res) => {
    res.send('Welcome to the TrailBrake API!');
});

// Use RideRouter whenever URI ends in /rides
app.use('/rides', RideRouter);

// Use ProfileRouter whenever URI ends in /profile
app.use('/profile', ProfileRouter);

// Start server
app.listen(8080, () => console.log('Server started'));