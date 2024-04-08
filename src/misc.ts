import JSEncrypt from 'jsencrypt'
import Swal from 'sweetalert2'
import Chart from 'chart.js/auto'
import annotation from 'chartjs-plugin-annotation'
import type { ChartConfiguration } from 'chart.js'
import type { AnnotationOptions } from 'chartjs-plugin-annotation'
import { date, trading, info } from './stores'
import { get } from 'svelte/store'

Chart.register(annotation)

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
  const label = labels[Number(value)]
  if (label.includes(':00') || label.includes(':30')) return label
  return null
}

const lunch_break: AnnotationOptions = {
  type: 'line',
  scaleID: 'x',
  drawTime: 'beforeDraw',
  value: 120,
  borderColor: '#e5e5e5',
  borderWidth: 2
}

export const encrypt = (pubkey: string, password: string) => {
  const encrypt = new JSEncrypt()
  encrypt.setPublicKey(pubkey)
  return encrypt.encrypt(password)
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

export const dateStr = (day: Date) => {
  const yyyy = day.getFullYear()
  const mm = String(day.getMonth() + 1).padStart(2, '0')
  const dd = String(day.getDate()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd}`
}

export const checkTradingTime = async (reload?: boolean): Promise<boolean> => {
  const d = new Date()
  const weekday = d.getDay()
  const hour = d.getUTCHours()
  if (weekday < 1 || weekday > 5 || hour < 1 || hour > 7) return false
  if (dateStr(d) === get(date)) return get(trading)
  if (reload) {
    await info()
    return await checkTradingTime(true)
  }
  return true
}

export const color = (last: number, value?: number) => {
  if (value === undefined) {
    if (last < 0) return 'color:green'
    else if (last > 0) return 'color:red'
  } else if (value) {
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
      },
      y2: {
        position: 'right',
        grid: {
          drawTicks: false,
          drawOnChartArea: false
        },
        ticks: { padding: 12 }
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
        annotations: {
          lunch_break,
          last: {
            type: 'line',
            scaleID: 'y',
            drawTime: 'beforeDraw',
            borderWidth: 1
          }
        }
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
          callback: value => {
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
          label: tooltipItem => {
            const label = tooltipItem.dataset.label
            const value = Math.round(tooltipItem.parsed.y * 10000) / 10000
            return `${label}   ${value}亿`
          }
        }
      },
      annotation: {
        annotations: {
          lunch_break,
          zero: {
            type: 'line',
            scaleID: 'y',
            drawTime: 'beforeDraw',
            value: 0,
            borderWidth: 1
          }
        }
      }
    }
  },
  plugins: [
    {
      id: 'nodata',
      afterDraw: chart => {
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
