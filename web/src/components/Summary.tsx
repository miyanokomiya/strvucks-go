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
        <Typography variant="body1">{formatter.format(summary.latestMovingTime)} min</Typography>
        <Typography variant="body1">{formatter.format(summary.latestCalories)} kcal</Typography>
        <Typography variant="subtitle1" color="textSecondary">
          Weekly ({summary.weeklyCount})
        </Typography>
        <Typography variant="body1">{formatter.format(summary.weeklyDistance)} m</Typography>
        <Typography variant="body1">{formatter.format(summary.weeklyMovingTime)} min</Typography>
        <Typography variant="body1">{formatter.format(summary.weeklyCalories)} kcal</Typography>
        <Typography variant="subtitle1" color="textSecondary">
          Monthly ({summary.monthlyCount})
        </Typography>
        <Typography variant="body1">{formatter.format(summary.monthlyDistance)} m</Typography>
        <Typography variant="body1">{formatter.format(summary.monthlyMovingTime)} min</Typography>
        <Typography variant="body1">{formatter.format(summary.monthlyCalories)} kcal</Typography>
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
