require('dotenv').config();
const mongoose = require('mongoose');

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
db.on('error', (error) => console.error(error));
db.once('open', () => console.log('Connected to MongoDB'));

// Export to share a single connection
module.exports = db;