const express = require('express');
const isEven = require('is-even');

const app = express();

app.get('/:number', (req, res) => {
    res.send(isEven(req?.params?.number ?? 0));
});


const PORT = process.env.PORT || 8080;

app.listen(PORT, () => {
    console.log(`Server is running on port ${PORT}`);
});
