const formatter = new Intl.NumberFormat('ja', {
  useGrouping: true,
  minimumFractionDigits: 0,
  maximumFractionDigits: 0,
})

export const formatNatural = (num: number): string => {
  return formatter.format(num)
}

export const formatTime = (sec: number): string => {
  const h = Math.floor(sec / 60 / 60)
  const m = Math.floor(sec / 60 - 60 * h)
  const s = sec % 60
  if (h === 0) {
    return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
  }
  return `${h}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

export const formatLap = (meter: number, sec: number): string => {
  const lap = sec / (meter / 1000)
  const m = Math.floor(lap / 60)
  const s = Math.floor(lap - 60 * m)
  return `${m}:${String(s).padStart(2, '0')}/km`
}
