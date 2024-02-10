import express from 'express';
import fs from 'fs';
import path from 'path';

const app = express();
const PORT = 6000;

app.use(express.json());

const dbPath = path.join(__dirname, '..', '..', 'db.json');

app.get('/todos', (req, res) => {
  fs.readFile(dbPath, 'utf8', (err, data) => {
    if (err) {
      return res.status(500).send({ error: 'Error reading from db.json' });
    }
    res.status(200).send(JSON.parse(data).todos);
  });
});



app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});
