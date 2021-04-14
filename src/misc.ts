import Swal from 'sweetalert2'
import Chart from 'chart.js/auto'
import annotation from 'chartjs-plugin-annotation'
import type { Stock } from './stores'
import type { ChartConfiguration } from 'chart.js'

Chart.register(annotation)

const color = (last: number, value?: number) => {
  if (value === undefined) {
    if (last < 0) return 'color:green'
    else if (last > 0) return 'color:red'
  } else if (value) {
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

const callback = (value: string | number) => {
  value = labels[value as number]
  if (value.includes(':00') || value.includes(':30')) return value
  return ''
}

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
    '#dc3545', // red
    '#fd7e14', // orange
    '#ffc107', // yellow
    '#28a745', // green
    '#007bff', // blue
    '#6f42c1', // purple
    '#6c757d' // gray
  ]
  return chartColors[i % chartColors.length]
}

export const intraday: ChartConfiguration<'line'> = {
  type: 'line',
  data: {
    labels,
    datasets: [
      {
        data: [],
        label: 'Price',
        fill: false,
        tension: 0,
        borderWidth: 2,
        borderColor: '#17a2b8',
        backgroundColor: '#17a2b8',
        pointRadius: 0,
        pointHoverRadius: 3
      }
    ]
  },
  options: {
    maintainAspectRatio: false,
    hover: {
      mode: 'index',
      intersect: false
    },
    animation: false,
    scales: {
      x: {
        grid: { drawTicks: false },
        ticks: {
          padding: 10,
          maxRotation: 0,
          callback
        }
      },
      y: {
        grid: { drawTicks: false },
        ticks: {
          padding: 12,
          format: { useGrouping: false }
        }
      }
    },
    plugins: {
      legend: { display: false },
      tooltip: {
        mode: 'index',
        intersect: false,
        displayColors: false,
        backgroundColor: 'rgba(210, 210, 210, 0.8)',
        titleColor: 'black',
        bodyFont: {
          weight: 'bold',
          size: 15
        },
        callbacks: {}
      },
      annotation: {
        annotations: [
          {
            type: 'line',
            scaleID: 'y',
            borderColor: 'black',
            borderWidth: 0.75
          }
        ]
      }
    }
  }
}

export const capitalflows: ChartConfiguration<'line'> = {
  type: 'line',
  data: {
    labels,
    datasets: []
  },
  options: {
    maintainAspectRatio: false,
    animation: false,
    scales: {
      x: {
        grid: { drawTicks: false },
        ticks: {
          padding: 10,
          maxRotation: 0,
          callback
        }
      },
      y: {
        grid: { drawTicks: false },
        ticks: {
          padding: 12,
          callback: (value) => {
            let suffix = ''
            if (value) suffix = '亿'
            return value + suffix
          }
        }
      }
    },
    plugins: {
      legend: {
        position: 'right',
        fullSize: false,
        labels: {
          boxWidth: 12,
          boxHeight: 12
        }
      },
      tooltip: {
        callbacks: {
          label: (tooltipItem) => {
            const label = tooltipItem.dataset.label
            const value = Math.round(tooltipItem.parsed.y * 10000) / 10000
            return `${label}   ${value}亿`
          }
        }
      },
      annotation: {
        annotations: [
          {
            type: 'line',
            scaleID: 'y',
            value: 0,
            borderColor: 'black',
            borderWidth: 0.75
          }
        ]
      }
    }
  },
  plugins: [
    {
      id: 'nodata',
      afterDraw: (chart) => {
        if (!chart.data.datasets || !chart.data.datasets.length) {
          const ctx = chart.ctx
          const width = chart.width
          const height = chart.height
          ctx.save()
          ctx.textAlign = 'center'
          ctx.textBaseline = 'middle'
          ctx.font = 'italic bold 48px Arial'
          ctx.fillStyle = '#666666'
          ctx.fillText('No data', width / 2, height / 2)
          ctx.restore()
        }
      }
    }
  ]
}
