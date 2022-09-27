const mongoose = require('mongoose');
const db = require('../database/connection');

const rideRecordSchema = new mongoose.Schema({
    rideName: { type: String, required: true },
    rideDate: { type: Date, required: true }
}, { collection: 'rideRecords' });

module.exports = db.model('rideRecord', rideRecordSchema);