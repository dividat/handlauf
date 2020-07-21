import mqtt from 'mqtt';

let client = mqtt.connect('ws://0.0.0.0:8080');

client.on('connect', () => {
  client.subscribe('values');
});

client.on('message', (_, data) => {
  console.log(JSON.parse(data));
});
