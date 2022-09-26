const express = require('express');
const router = express.Router();
const Ride = require('../model/ride');

async function findRideByID(req, res, next) {
    let ride;
    try {
        ride = await Ride.findById(req.params.rideObjID);
        if (ride == null) {
            return res.status(404).json({ message: 'Cannot find ride' });
        } 
    } catch (err) {
        return res.status(500).json({ 
            "method": "findRideByID()", 
            message: err.message 
        });
    }

    res.ride = ride;
    next();
}

// Get all
router.get('/', async (req, res) => {
    try {
        // rideData not included with this call to reduce network IO
        const rides = await Ride.find({}, { rideData: 0 });
        res.json(rides);
    } catch (err) {
        res.status(500).json({ message: err.message });
    }
});

// Get one
router.get('/:rideObjID', findRideByID, (req, res) => {
   res.json(res.ride);
});

// Create one
router.post('/', async (req, res) => {
    let ride = new Ride({
        rideName: req.body.rideName,
        rideDate: req.body.rideDate,
        rideData: req.body.rideData
    });

    try {
        const newRide = await ride.save(function(err) {
            if (err) {
                res.status(400).json(err);
            } else {
                res.status(201).json(newRide);
            }
        });
    } catch (err) {
        res.status(400);
    }
});

// Delete one
router.delete('/:rideObjID', findRideByID, async (req, res) => {
    try {
        await res.ride.remove();
        res.json({ message: 'Successfully deleted ride' });
    } catch (err) {
        res.status(500).json({ 
            "method": "router.delete()", 
            message: err.message 
        });
    }
});

module.exports = router;