// import dotenv from 'dotenv';
// import got from 'got';
import express from 'express';
// import RideData from '../model/ride_data.js';
import RideRecord from '../model/ride_record.js';

// dotenv.config();
const router = express.Router();

// Get user score summary
// Endpoint sample: /profile/score?useRecentRides=20
router.get('/score', async (req, res) => {
    let limitCount = req.query.useRecentRides ? parseInt(req.query.useRecentRides) : 10;

    let pipeline = [
        { $sort: { 'createdAt': -1 } },
        { $limit: limitCount },
        {
            $group: {
                _id: null,
                '_totalDuration': { $sum: '$rideMeta.duration' },
                '_sumProdOverallScore': {
                    $sum: {
                        $multiply: [ '$rideScore.overall', '$rideMeta.duration' ]
                    }
                },
                '_sumProdAccelerationScore': {
                    $sum: {
                        $multiply: [ '$rideScore.acceleration', '$rideMeta.duration' ]
                    }
                },
                '_sumProdBrakingScore': {
                    $sum: {
                        $multiply: [ '$rideScore.braking', '$rideMeta.duration' ]
                    }
                },
                '_sumProdCorneringScore': {
                    $sum: {
                        $multiply: [ '$rideScore.cornering', '$rideMeta.duration' ]
                    }
                },
                '_sumProdSpeedScore': {
                    $sum: {
                        $multiply: [ '$rideScore.speed', '$rideMeta.duration' ]
                    }
                },
            }
        },
        {
            $project: {
                'overall': { $divide: [ '$_sumProdOverallScore', '$_totalDuration' ] },
                'acceleration': { $divide: [ '$_sumProdAccelerationScore', '$_totalDuration' ] },
                'braking': { $divide: [ '$_sumProdBrakingScore', '$_totalDuration' ] },
                'cornering': { $divide: [ '$_sumProdCorneringScore', '$_totalDuration' ] },
                'speed': { $divide: [ '$_sumProdSpeedScore', '$_totalDuration' ] }
            }
        }
    ];

    try {
        const scoreSummary = await RideRecord.aggregate(pipeline);

        res.json(scoreSummary[0]);
    } catch (err) {
        res.status(500).json({ message: err.message });
    }
});

// Get user lifetime statistics
router.get('/lifetime', async (req, res) => {
    let pipeline = [
        {
            $group: {
                _id: null,
                'totalDistance': { $sum: '$rideMeta.distance' },
                'totalDuration': { $sum: '$rideMeta.duration' },
                'maxAcceleration': { $max: '$rideMeta.maxAcceleration' }
            }
        }
    ];

    try {
        const lifetimeStats = await RideRecord.aggregate(pipeline);

        res.json(lifetimeStats[0]);
    } catch (err) {
        res.status(500).json({ message: err.message });
    }
});

export default router;