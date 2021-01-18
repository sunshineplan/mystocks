import Swal from 'sweetalert2'
import type { Stock } from './stores'

export const fire = (
  title?: string | undefined,
  html?: string | undefined,
  icon?: 'success' | 'error' | 'warning' | 'info' | 'question' | undefined
) => {
  const swal = Swal.mixin({
    customClass: { confirmButton: 'swal btn btn-primary' },
    buttonsStyling: false
  })
  return swal.fire(title, html, icon)
}

export const valid = () => {
  let result = true
  Array.from(document.querySelectorAll('input'))
    .forEach(i => { if (!i.checkValidity()) result = false })
  return result
}

export const post = (url: string, data?: object) => {
  return fetch(url, {
    method: 'post',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
}

export const checkTime = () => {
  var date = new Date();
  var hour = date.getUTCHours();
  var day = date.getDay();
  if (hour >= 1 && hour <= 8 && day >= 1 && day <= 5)
    return true
  return false
}

export const color = (last: number, value?: number) => {
  if (value === undefined) {
    if (last < 0) return 'color:green'
    else if (last > 0) return 'color:red'
  } else {
    if (last < value) return 'color:red'
    else if (last > value) return 'color:green'
  }
  return 'color:initial'
}

export const addColor = (stock: Stock, val: string) => {
  if (stock.name != 'n/a') {
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
  return 'color:initial'
}

export const timeLabels = (start: number, end: number) => {
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
