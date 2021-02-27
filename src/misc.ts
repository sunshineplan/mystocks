import Swal from 'sweetalert2'
import Chart from 'chart.js'
import annotation from 'chartjs-plugin-annotation'
import type { Stock } from './stores'

Chart.plugins.register(annotation)

const color = (last: number, value?: number) => {
  if (value === undefined) {
    if (last < 0) return 'color:green'
    else if (last > 0) return 'color:red'
  } else {
    if (last < value) return 'color:red'
    else if (last > value) return 'color:green'
  }
  return 'color:initial'
}

const timeLabels = (start: number, end: number) => {
  let times: string[] = []
  for (let i = 0; start <= end; i++) {
    times[i] = `${Math.floor(start / 60)
      .toString()
      .padStart(2, '0')}:${(start % 60).toString().padStart(2, '0')}`
    start++
  }
  return times
}

const labels = timeLabels(9 * 60 + 30, 11 * 60 + 30).concat(timeLabels(13 * 60 + 1, 15 * 60))

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

export const post = (url: string, data?: object, universal?: boolean) => {
  const init: RequestInit = {
    method: 'post',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }
  if (universal) init.credentials = 'include'
  return fetch(url, init)
}

export const checkTime = () => {
  const date = new Date()
  const hour = date.getUTCHours()
  const day = date.getDay()
  if (hour >= 1 && hour <= 8 && day >= 1 && day <= 5)
    return true
  return false
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

export const getColor = (i: number) => {
  const chartColors = [
    'rgb(255, 99, 132)', // red
    'rgb(255, 159, 64)', // orange
    'rgb(255, 205, 86)', // yellow
    'rgb(75, 192, 192)', // green
    'rgb(54, 162, 235)', // blue
    'rgb(153, 102, 255)', // purple
    'rgb(201, 203, 207)' // grey
  ]
  return chartColors[i % chartColors.length]
}

export const intraday = {
  type: 'line',
  data: {
    labels,
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
    maintainAspectRatio: false,
    legend: { display: false },
    hover: {
      mode: 'index',
      intersect: false
    },
    tooltips: {
      mode: 'index',
      intersect: false,
      displayColors: false
    },
    animation: { duration: 0 },
    scales: {
      xAxes: [
        {
          gridLines: { drawTicks: false },
          ticks: {
            padding: 10,
            autoSkipPadding: 100,
            maxRotation: 0,
          }
        }
      ],
      yAxes: [
        {
          gridLines: { drawTicks: false },
          ticks: { padding: 12 }
        }
      ]
    },
    annotation: {
      annotations: [
        {
          id: 'PreviousClose',
          type: 'line',
          mode: 'horizontal',
          scaleID: 'y-axis-0',
          borderColor: 'black',
          borderWidth: 0.75
        }
      ]
    }
  }
} as Chart.ChartConfiguration

export const capitalflows = {
  type: 'line',
  data: { labels },
  options: {
    maintainAspectRatio: false,
    legend: { position: 'right' },
    animation: { duration: 0 },
    tooltips: {
      callbacks: {
        label: (tooltipItem) => {
          const value = tooltipItem.value as string
          return Math.round(parseFloat(value) * 10000) / 10000 + '亿'
        }
      }
    },
    scales: {
      xAxes: [
        {
          gridLines: { drawTicks: false },
          ticks: {
            padding: 10,
            maxTicksLimit: 9,
            maxRotation: 0
          }
        }
      ],
      yAxes: [
        {
          gridLines: { drawTicks: false },
          ticks: {
            padding: 12,
            callback: (value) => {
              if (value) return value + '亿'
              else return value
            }
          }
        }
      ]
    },
    annotation: {
      annotations: [
        {
          id: 'zero',
          type: 'line',
          mode: 'horizontal',
          scaleID: 'y-axis-0',
          value: 0,
          borderColor: 'black',
          borderWidth: 0.75
        }
      ]
    }
  },
  plugins: [
    {
      afterDraw: (chart) => {
        if (!chart.data.datasets || !chart.data.datasets.length) {
          const ctx = chart.ctx as CanvasRenderingContext2D
          const width = chart.width as number
          const height = chart.height as number
          ctx.save()
          ctx.textAlign = 'center'
          ctx.textBaseline = 'middle'
          ctx.font = 'italic bold 48px Arial'
          ctx.fillText('No data', width / 2, height / 2)
          ctx.restore()
        }
      }
    }
  ]
} as Chart.ChartConfiguration
