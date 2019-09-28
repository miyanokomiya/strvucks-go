import { User, Summary } from '../src/types'

export const mockUser = (): User => {
  return {
    id: 1,
    athleteId: 10,
    username: 'my name',
    iftttKey: 'my key',
    iftttMessage: 'my message',
  }
}

export const mockSummary = (): Summary => {
  return {
    id: 1,
    athleteId: 2,

    latestDistance: 3,
    latestMovingTime: 4,
    latestTotalElevationGain: 5,
    latestCalories: 6,

    monthBaseDate: '2018/01/08',
    monthlyCount: 14,
    monthlyDistance: 14,
    monthlyMovingTime: 15,
    monthlyTotalElevationGain: 16,
    monthlyCalories: 17,

    weekBaseDate: '2018/01/01',
    weeklyCount: 1004,
    weeklyDistance: 1005,
    weeklyMovingTime: 100006,
    weeklyTotalElevationGain: 1007,
    weeklyCalories: 1008,
  }
}
