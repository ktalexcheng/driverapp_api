import RideData from '../model/ride_data.js';
import RideRecord from '../model/ride_record.js';

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

export { findRideRecordByID, findRideDataByRecordID };