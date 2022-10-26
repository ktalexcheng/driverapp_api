import dotenv from 'dotenv';
import got from 'got';
import express from 'express';
import RideData from '../model/ride_data.js';
import RideRecord from '../model/ride_record.js';

dotenv.config();
const router = express.Router();

async function findRideRecordByID(req, res, next) {
    try {
        const rideRecord = await RideRecord.findById(req.params.rideObjID);
        if (rideRecord == null) {
            return res.status(404).json({ message: 'Cannot find ride' });
        } 

        res.rideRecord = rideRecord;
        next();
    } catch (err) {
        return res.status(500).json({ 
            middleware: "findRideRecordByID()", 
            message: err.message 
        });
    }
}

async function findRideDataByRecordID(req, res, next) {
    try {
        const rideData = await RideData.find({ 
            'metadata.rideRecordID': res.rideRecord._id
        });
        if (rideData.length == 0) {
            return res.status(404).json({ message: 'Cannot find ride data' });
        } 

        res.rideData = rideData;
        next();
    } catch (err) {
        return res.status(500).json({ 
            middleware: "findRideDataByRecordID()", 
            message: err.message 
        });
    }
}

// Get all
router.get('/', async (req, res) => {
    try {
        const rideRecord = await RideRecord.find().sort({ createdAt: -1 }).limit(20);

        res.json(rideRecord);
    } catch (err) {
        res.status(500).json({ message: err.message });
    }
});

// Get one
router.get('/:rideObjID', [findRideRecordByID, findRideDataByRecordID], async (req, res) => {
    try{    
        res.json({
            _id: res.rideRecord._id,
            rideName: res.rideRecord.rideName,
            rideDate: res.rideRecord.rideDate,
            rideScore: res.rideRecord.rideScore,
            rideMeta: res.rideRecord.rideMeta,
            rideData: res.rideData
        });
    } catch (err) {
        return res.status(500).json({ 
            message: err.message 
        });
    }
});

// Create one
router.post('/', async (req, res) => {
    let rideRecord;
    let allRideData = [];

    try {
        // Create new rideRecord document
        rideRecord = new RideRecord({
            rideName: req.body.rideName,
            rideDate: req.body.rideDate
        });

        // Create documents for rideData belonging to rideRecord
        for (let k in req.body.rideData) {
            allRideData.push(new RideData({
                timestamp: req.body.rideData[k].timestamp,
                metadata: { rideRecordID: rideRecord._id },
                locationLat: req.body.rideData[k].locationLat,
                locationLong: req.body.rideData[k].locationLong,
                accelerometerZ: req.body.rideData[k].accelerometerX,
                accelerometerX: req.body.rideData[k].accelerometerY,
                accelerometerY: req.body.rideData[k].accelerometerZ,
                gyroscopeX: req.body.rideData[k].gyroscopeX,
                gyroscopeY: req.body.rideData[k].gyroscopeY,
                gyroscopeZ: req.body.rideData[k].gyroscopeZ
            }));
        }

        // Insert data to collection
        RideData.collection.insertMany(allRideData, function(err) {
            if (err) {
                res.status(400).json(err);
            }
        });
    } catch (err) {
        res.status(400).json(err);
    }

    // Evaluate ride and update rideRecord document with ride score using TrailBrake Judge service        
    const judgeApiUrl = `${process.env.JUDGE_URL}/rideScore`;
    const judgeApiOptions = {
        json: {
            rideData: allRideData
        },
        retry: {
            limit: 0
        }
    };
    
    try {
        const { rideScore, rideMeta } = await got.post(judgeApiUrl, judgeApiOptions).json();

        // Update document object properties
        rideRecord.rideScore = rideScore;
        rideRecord.rideMeta = rideMeta;
    } catch (err) {
        res.status(err.response.statusCode).json(JSON.parse(err.response.body));
    }

    rideRecord.save(function(err, doc) {
        if (err) {
            res.status(400).json(err);
        } else {
            res.status(201).json({
                message: 'New ride successfully created',
                rideRecord: doc
            });
        }
    });
});

// Delete one
router.delete('/:rideObjID', findRideRecordByID, async (req, res) => {
    try {
        RideData.deleteMany({ 
            'metadata.rideRecordID': res.rideRecord._id
        }, function(err) {
            if (err) {
                res.status(400).json(err);
            } 
        });
        await res.rideRecord.remove();

        res.json({ message: 'Successfully deleted ride' });
    } catch (err) {
        res.status(500).json({ 
            method: "router.delete()", 
            message: err.message 
        });
    }
});

export default router;