const mongoose = require('mongoose');
const db = require('../database/connection');

// rideData object schema
const rideDataSchema = new mongoose.Schema({
    timestamp: Date,
    locationLat: { type: Number, required: false },
    locationLong: { type: Number, required: false },
    accelerometerX: Number,
    accelerometerY: Number,
    accelerometerZ: Number,
    gyroscopeX: Number,
    gyroscopeY: Number,
    gyroscopeZ: Number
});

// Schema for 'Rides' collection
const rideSchema = new mongoose.Schema({
    rideName: { type: String, required: true },
    rideDate: { type: Date, required: true },
    rideData: { type: [rideDataSchema], required: true }
});

module.exports = db.model('Ride', rideSchema);