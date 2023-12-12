const express = require('express')
const fs = require('fs');
const app = express()
const port = 3000

app.get('/', (req, res) => {
    let rawdata = fs.readFileSync('data.json');
    let data = JSON.parse(rawdata);
    res.send(data)
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})