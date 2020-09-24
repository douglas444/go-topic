const express = require("express");
const bodyParser = require("body-parser");
const app = express();
const port = 8081;

app.use(bodyParser.json());

app.post("/notify", (req, res) => {
    console.log(req.body);
});

app.listen(port, () => {
    console.log("running on 8081");
});

