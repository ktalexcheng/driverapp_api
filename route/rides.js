const express = require('express');
const router = express.Router();
const RideData = require('../model/ride_data');
const RideRecord = require('../model/ride_record');

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
            "middleware": "findRideRecordByID()", 
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
            "middleware": "findRideDataByRecordID()", 
            message: err.message 
        });
    }
}

// Get all
router.get('/', async (req, res) => {
    try {
        const rideRecord = await RideRecord.find().limit(20);

        res.json(rideRecord);
    } catch (err) {
        res.status(500).json({ message: err.message });
    }
});

// Get one
router.get('/:rideObjID', [findRideRecordByID, findRideDataByRecordID], async (req, res) => {
    try{    
        res.json({
            rideName: res.rideRecord.rideName,
            rideDate: res.rideRecord.rideDate,
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
    try {
        // Create new rideRecord document
        const rideRecord = new RideRecord({
            rideName: req.body.rideName,
            rideDate: req.body.rideDate
        });

        rideRecord.save(function(err) {
            if (err) {
                res.status(400).json(err);
            } 
        });

        // Create documents for rideData belonging to rideRecord
        let allRideData = [];
        for (var k in req.body.rideData) {
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

        RideData.collection.insertMany(allRideData, function(err, docs) {
            if (err) {
                res.status(400).json(err);
            } else {
                res.status(201).json(docs);
            }
        });
    } catch (err) {
        res.status(400);
    }
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
            "method": "router.delete()", 
            message: err.message 
        });
    }
});

module.exports = router;