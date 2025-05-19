const express = require('express');
const app = express();

app.use(express.json());

app.post('/parse', (req, res) => {
    try {
        const input = req.body.data;
        
        // Simulating vulnerable code
        if (input && input.includes('trigger')) {
            throw new Error('Crash triggered!');
        }
        
        res.json({ success: true, parsed: input });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

const server = app.listen(3000, () => {
    console.log('Server running on port 3000');
});

module.exports = server;