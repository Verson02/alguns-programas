const express = require('express');
const request = require('request');
const app = express();
const port = 3000;

app.use(express.json());

app.get('/check-url', (req, res) => {
    const url = req.query.url;
    if (!url) {
        return res.status(400).send('URL query parameter is required');
    }

    request.head(url, (error, response) => {
        if (error || response.statusCode >= 400) {
            res.status(500).send('URL is offline');
        } else {
            res.send('URL is online');
        }
    });
});

app.listen(port, () => {
    console.log(`Proxy server running at http://localhost:${port}`);
});
