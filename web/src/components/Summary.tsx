import * as React from 'react'
import Typography from '@material-ui/core/Typography'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import PropTypes from 'prop-types'
import { Summary } from '../types'

const formatter = new Intl.NumberFormat('ja', {
  useGrouping: true,
  minimumFractionDigits: 0,
  maximumFractionDigits: 0,
})

const formatTime = (sec: number): string => {
  const h = Math.floor(sec / 60 / 60)
  const m = Math.floor(sec / 60 - 60 * h)
  if (h === 0) {
    return `${String(m).padStart(2, '0')}m`
  }
  return `${h}h${String(m).padStart(2, '0')}m`
}

const formatLap = (meter: number, sec: number): string => {
  const lap = sec / (meter / 1000)
  const m = Math.floor(lap / 60)
  const s = Math.floor(lap - 60 * m)
  return `${m}:${String(s).padStart(2, '0')}/km`
}

type Props = {
  summary: Summary
}

export const SummaryCard: React.FC<Props> = props => {
  const summary = props.summary
  return (
    <Card>
      <CardContent>
        <Typography variant="h6">Activity Summary</Typography>
        <Typography variant="subtitle1" color="textSecondary">
          Latest
        </Typography>
        <Typography variant="body1">{formatter.format(summary.latestDistance)} m</Typography>
        <Typography variant="body1">{formatTime(summary.latestMovingTime)}</Typography>
        <Typography variant="body1">
          {formatLap(summary.latestDistance, summary.latestMovingTime)}
        </Typography>
        <Typography variant="body1">{formatter.format(summary.latestCalories)} kcal</Typography>
        <Typography variant="subtitle1" color="textSecondary">
          Weekly ({summary.weeklyCount})
        </Typography>
        <Typography variant="body1">{formatter.format(summary.weeklyDistance)} m</Typography>
        <Typography variant="body1">{formatTime(summary.weeklyMovingTime)}</Typography>
        <Typography variant="body1">
          {formatLap(summary.weeklyDistance, summary.weeklyMovingTime)}
        </Typography>
        <Typography variant="subtitle1" color="textSecondary">
          Monthly ({summary.monthlyCount})
        </Typography>
        <Typography variant="body1">{formatter.format(summary.monthlyDistance)} m</Typography>
        <Typography variant="body1">{formatTime(summary.monthlyMovingTime)}</Typography>
        <Typography variant="body1">
          {formatLap(summary.monthlyDistance, summary.monthlyMovingTime)}
        </Typography>
      </CardContent>
    </Card>
  )
}

SummaryCard.propTypes = {
  summary: PropTypes.shape({
    id: PropTypes.number.isRequired,
    athleteId: PropTypes.number.isRequired,

    latestDistance: PropTypes.number.isRequired,
    latestMovingTime: PropTypes.number.isRequired,
    latestTotalElevationGain: PropTypes.number.isRequired,
    latestCalories: PropTypes.number.isRequired,

    weekBaseDate: PropTypes.string.isRequired,
    weeklyCount: PropTypes.number.isRequired,
    weeklyDistance: PropTypes.number.isRequired,
    weeklyMovingTime: PropTypes.number.isRequired,
    weeklyTotalElevationGain: PropTypes.number.isRequired,
    weeklyCalories: PropTypes.number.isRequired,

    monthBaseDate: PropTypes.string.isRequired,
    monthlyCount: PropTypes.number.isRequired,
    monthlyDistance: PropTypes.number.isRequired,
    monthlyMovingTime: PropTypes.number.isRequired,
    monthlyTotalElevationGain: PropTypes.number.isRequired,
    monthlyCalories: PropTypes.number.isRequired,
  }).isRequired,
}
