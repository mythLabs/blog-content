const express = require("express");
const app = express();
const port = 3000;

const slowDown = () => {
  let result = 0;
  for (var i = Math.pow(9, 7); i >= 0; i--) {
    result += Math.atan(i) * Math.tan(i);
  }
};

app.get("/", (req, res) => {
  slowDown();
  res.send("Hello World");
});

app.get("/products", (req, res) => {
  slowDown();
  res.send({ productId: req.query.id });
});

app.post("/signin", (req, res) => {
  slowDown();
  res.send("OK");
});

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`);
});
