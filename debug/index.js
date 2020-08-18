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
const lines = svg.append('g');
const bars = svg.append('g');

const colors = d3.scaleOrdinal(d3.schemeCategory10);

function viz() {
  // check length
  if (list.length === 0) {
    return;
  }

  // prepare x scale
  const xScale = d3
    .scaleTime()
    .domain(d3.extent(list, (d) => d.time))
    .range([0, width]);

  // apply x scale
  xAxis.attr('transform', 'translate(50,' + (50 + height) + ')').call(d3.axisBottom(xScale));

  // prepare y scale
  const yScale = d3.scaleLinear().domain([0, 1]).range([height, 0]);

  // apply y scale
  yAxis.attr('transform', 'translate(50,50)').call(d3.axisLeft(yScale));

  // prepare line generator
  const line = (i) => {
    return d3
      .line()
      .x(function (d) {
        return 50 + xScale(d.time);
      })
      .y(function (d) {
        return 50 + yScale(d.sample[i]);
      });
  };

  // prepare proto
  const proto = new Array(d3.max(list, (d) => d.sample.length));

  // update lines
  lines
    .selectAll('path')
    .data(proto)
    .join('path')
    .attr('fill', 'none')
    .attr('stroke', (_, i) => colors(i))
    .attr('stroke-width', 2)
    .attr('d', (_, i) => line(i)(list));

  // update bars
  bars
    .selectAll('rect')
    .data(proto)
    .join('rect')
    .attr('fill', (_, i) => colors(i))
    .attr('x', (_, i) => 100)
    .attr('y', (_, i) => 50 + i * 50)
    .attr('height', 50)
    .attr('width', (_, i) => 5 + list[list.length - 1].sample[i] * (width - 100));
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
