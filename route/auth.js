import express from 'express';
import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken';
import User from '../model/user.js';

const router = express.Router();

router.get('/token', async (req, res) => {
    // Parse request body with email and hashed password
    const { email, password } = req.body;

    if (!email || !password) {
        return res.status(400).json({
            message: "Request body is invalid."
        });
    }

    try {
        // Check if user email exists
        const user = await User.findOne({ email });
        if (!user) {
            return res.status(400).json({
                message: "User does not exist."
            });
        }

        // Check if user password match
        const pwMatch = await bcrypt.compare(password, user.password);
        if(!pwMatch) {
            return res.status(400).json({
                message: "Invalid password."
            });
        }

        // Create token
        const token  = jwt.sign({ userId: user._id }, 'secret');
        res.json({ token });
    } catch (err) {
        res.status(500).json({
            message: err.message
        });
    }
});

router.post('/signup', async (req, res) => {
    // Parse request body with email and hashed password
    const { email, password } = req.body;

    if (!email || !password) {
        return res.status(400).json({
            message: "Request body is invalid."
        });
    }

    try {
        // Check if user email exists
        const user = await User.findOne({ email });
        if (user) {
            return res.status(400).json({
                message: "User already exist."
            });
        }

        // Create new user
        const newUser = new User({
            email: email,
            password: password
        });
        newUser.set({
            userId: newUser._id.toString().substring(0, 8)
        });
        console.log(newUser.userId);
        await newUser.save();

        // Create token
        const token  = jwt.sign({ userId: newUser._id }, 'secret');
        res.json({ token });
    } catch (err) {
        res.status(500).json({
            message: err.message
        });
    }
});

export default router;