import mongoose from 'mongoose';
import db from '../database/connection.js';

const rideRecordSchema = new mongoose.Schema({
    rideName: { type: String, required: true },
    rideDate: { type: Date, required: true },
    rideScore: {
        overall: Number,
        speed: Number,
        acceleration: Number,
        braking: Number,
        cornering: Number,
    }
}, { 
    collection: 'rideRecords',
    timestamps: true 
});

const rideRecordSchemaModel = db.model('rideRecord', rideRecordSchema);

export default rideRecordSchemaModel;