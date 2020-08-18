import * as d3 from 'd3';
import WebSocket from 'reconnecting-websocket';

/* Data */

const list = [];

/* Visualization */

const width = window.innerWidth - 100;
const height = window.innerHeight - 100;

const svg = d3.select('#viz');

const xAxis = svg.append('g');
const yAxis = svg.append('g');
const paths = svg.append('g');

const colors = d3.scaleOrdinal(d3.schemeCategory10);

function viz() {
  if (list.length === 0) {
    return;
  }

  const x = d3
    .scaleTime()
    .domain(d3.extent(list, (d) => d.time))
    .range([0, width]);

  xAxis.attr('transform', 'translate(50,' + (50 + height) + ')').call(d3.axisBottom(x));

  const y = d3.scaleLinear().domain([0, 1]).range([height, 0]);

  yAxis.attr('transform', 'translate(50,50)').call(d3.axisLeft(y));

  const line = (i) => {
    return d3
      .line()
      .x(function (d) {
        return 50 + x(d.time);
      })
      .y(function (d) {
        return 50 + y(d.sample[i]);
      });
  };

  paths
    .selectAll('path')
    .data(new Array(list[0].sample.length))
    .join('path')
    .attr('fill', 'none')
    .attr('stroke', (_, i) => colors(i))
    .attr('stroke-width', 2)
    .attr('d', (_, i) => line(i)(list));
}

/* WebSocket */

const ws = new WebSocket('ws://0.0.0.0:8080');

ws.onmessage = (event) => {
  // parse sample
  const sample = JSON.parse(event.data.toString());

  // prepare entry
  const entry = {
    sample: sample,
    time: new Date(),
  };

  // add to list
  list.push(entry);

  // trim list
  while (list.length > 1000) {
    list.shift();
  }

  // update
  viz();
};
