import Swal from 'sweetalert2'
import { goto } from '@sapper/app'

export const BootstrapButtons = Swal.mixin({
  customClass: { confirmButton: 'swal btn btn-primary' },
  buttonsStyling: false
})

export function valid() {
  var result = true
  Array.from(document.querySelectorAll('input'))
    .forEach(i => { if (!i.checkValidity()) result = false })
  return result
}

export function post(url, data) {
  return fetch(url, {
    method: 'post',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
}

export function checkTime() {
  var date = new Date();
  var hour = date.getUTCHours();
  var day = date.getDay();
  if (hour >= 1 && hour <= 8 && day >= 1 && day <= 5)
    return true
  return false
}

export function color(last, value) {
  if (value == undefined)
    if (last < 0) return { color: 'green' }
    else if (last > 0) return { color: 'red' }
  if (last < value) return { color: 'red' }
  else if (last > value) return { color: 'green' }
}

export function addColor(stock, val) {
  if (stock && stock.name != 'n/a') {
    switch (val) {
      case 'change':
      case 'percent':
        return color(stock.change)
      case 'now':
        return color(stock.last, stock.now)
      case 'high':
        return color(stock.last, stock.high)
      case 'low':
        return color(stock.last, stock.low)
      case 'open':
        return color(stock.last, stock.open)
    }
  }
}

export async function gotoStock(stock) { await goto(`/stock/${stock.index}/${stock.code}`) }

export function timeLabels(start, end) {
  var times = [];
  for (var i = 0; start <= end; i++) {
    times[i] = `${Math.floor(start / 60).toString().padStart(2, '0')}:${(start % 60).toString().padStart(2, '0')}`
    start++
  }
  return times
}

export const intraday = {
  type: 'line',
  data: {
    labels: timeLabels(9 * 60 + 30, 11 * 60 + 30).concat(timeLabels(13 * 60 + 1, 15 * 60)),
    datasets: [
      {
        label: 'Price',
        fill: false,
        lineTension: 0,
        borderWidth: 2,
        borderColor: 'red',
        backgroundColor: 'red',
        pointRadius: 0,
        pointHoverRadius: 3
      }
    ]
  },
  options: {
    scales: {
      xAxes: [{
        gridLines: { drawTicks: false },
        ticks: {
          padding: 10,
          autoSkipPadding: 100,
          maxRotation: 0
        }
      }],
      yAxes: [{
        gridLines: { drawTicks: false },
        ticks: { padding: 12 }
      }]
    },
    annotation: {
      annotations: [
        {
          id: 'PreviousClose',
          type: 'line',
          mode: 'horizontal',
          scaleID: 'y-axis-0',
          borderColor: 'black',
          borderWidth: .75
        }
      ]
    }
  }
}
