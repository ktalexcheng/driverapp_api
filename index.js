require('dotenv').config();
require('./database/connection');
const express = require('express');

const app = express();

// Middleware: what happens after server gets request but before passing to route
// Use middleware to accept JSON as body
app.use(express.json({ limit: '10mb' }));

app.get('/', (req, res) => {
    res.send('Welcome to the DriverApp API!');
});

// Set up middleware for routing
const rideRouter = require('./route/rides');
// Use rideRouter whenever URI ends in /rides
app.use('/rides', rideRouter);

app.listen(8080, () => console.log('Server started'));