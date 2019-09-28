import * as formatter from '../src/formatter'

describe('formatter', () => {
  describe('formatNatural', () => {
    const data = [
      { arg: 0, exp: '0' },
      { arg: 123, exp: '123' },
      { arg: 1234, exp: '1,234' },
      { arg: 1234.1234, exp: '1,234' },
    ]
    data.forEach(({ arg, exp }) => {
      it(`${arg} => ${exp}`, () => {
        expect(formatter.formatNatural(arg)).toBe(exp)
      })
    })
  })

  describe('formatTime', () => {
    const data = [
      { arg: 0, exp: '00:00' },
      { arg: 123, exp: '02:03' },
      { arg: 59 * 60 + 59, exp: '59:59' },
      { arg: 60 * 60 + 59, exp: '1:00:59' },
      { arg: 100 * 60 * 60 + 59, exp: '100:00:59' },
    ]
    data.forEach(({ arg, exp }) => {
      it(`${arg} => ${exp}`, () => {
        expect(formatter.formatTime(arg)).toBe(exp)
      })
    })
  })

  describe('formatLap', () => {
    const data = [
      { meter: 1000, sec: 0, exp: '0:00/km' },
      { meter: 1000, sec: 123, exp: '2:03/km' },
      { meter: 500, sec: 123, exp: '4:06/km' },
      { meter: 2000, sec: 123, exp: '1:01/km' },
      { meter: 1000, sec: 59 * 60 + 59, exp: '59:59/km' },
      { meter: 1000, sec: 60 * 60 + 59, exp: '60:59/km' },
    ]
    data.forEach(({ meter, sec, exp }) => {
      it(`${meter}, ${sec} => ${exp}`, () => {
        expect(formatter.formatLap(meter, sec)).toBe(exp)
      })
    })
  })
})
