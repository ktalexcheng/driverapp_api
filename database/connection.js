import dotenv from 'dotenv';
import mongoose from 'mongoose';
import path from 'path';
import { fileURLToPath } from 'url';

dotenv.config();
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Configure connection to database
mongoose.connect(process.env.DATABASE_URL, {
    useNewUrlParser: true,
    ssl: true,
    sslValidate: true,
    sslCert: `${__dirname}/../certificate/atlas-admin-X509-cert.pem`,
    sslKey: `${__dirname}/../certificate/atlas-admin-X509-cert.pem`,
    authMechanism: 'MONGODB-X509',
});

// Connect to database
const db = mongoose.connection;
db.on('error', (err) => console.error(err));
db.once('open', () => console.log('Connected to MongoDB'));

// Export to share a single connection
export default db;