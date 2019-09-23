export interface User {
  id: number
  athleteId: number
  username: string
  iftttKey: string
  iftttMessage: string
}

export interface Summary {
  id: number
  athleteId: number

  latestDistance: number
  latestMovingTime: number
  latestTotalElevationGain: number
  latestCalories: number

  monthBaseDate: string
  monthlyCount: number
  monthlyDistance: number
  monthlyMovingTime: number
  monthlyTotalElevationGain: number
  monthlyCalories: number

  weekBaseDate: string
  weeklyCount: number
  weeklyDistance: number
  weeklyMovingTime: number
  weeklyTotalElevationGain: number
  weeklyCalories: number
}

