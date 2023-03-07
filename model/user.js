import mongoose from 'mongoose';
import db from '../database/connection.js';

const userSchema = new mongoose.Schema({
    userId: { type: String },
    email: { type: String, required: true },
    password: { type: String, required: true }
}, {
    collection: 'users',
    timestamps: true
});

const userSchemaModel = db.model('users', userSchema);

export default userSchemaModel;