import jwt from 'jsonwebtoken';

async function verifyToken(req, res, next) {
    const token = req.header('Authorization');

    if (!token) {
        return res.status(401).json({
            message: "Authorization denied."
        });
    }

    try {
        const decoded = jwt.verify(token, 'secret');
        req.user = decoded.user;
        next();
    } catch (err) {
        res.status(401).json({
            message: "Please login."
        });
    }
}

export { verifyToken };