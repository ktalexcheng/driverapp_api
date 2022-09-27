const mongoose = require('mongoose');
const db = require('../database/connection');

const rideDataSchema = new mongoose.Schema({
    timestamp: Date,
    metadata: { 
        rideRecordID: { type: mongoose.Schema.Types.ObjectId, required: true } 
    },
    locationLat: { type: Number, required: false },
    locationLong: { type: Number, required: false },
    accelerometerX: Number,
    accelerometerY: Number,
    accelerometerZ: Number,
    gyroscopeX: Number,
    gyroscopeY: Number,
    gyroscopeZ: Number
}, { collection: 'rideData' });

module.exports = db.model('rideData', rideDataSchema);