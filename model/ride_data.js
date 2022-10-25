import mongoose from 'mongoose';
import db from '../database/connection.js';

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
}, { 
    collection: 'rideData',
    timestamps: true
});

const rideDataSchemaModel = db.model('rideData', rideDataSchema);

export default rideDataSchemaModel;