import * as React from 'react'
import axios from 'axios'
import renderer from 'react-test-renderer'
import { actWait } from '../helper'
import { App } from '../../src/App'
import { mockUser, mockSummary } from '../mock'

jest.mock('axios')

describe('src/App', () => {
  describe('snapshot', () => {
    ;(axios.get as jest.Mock).mockImplementation((url: string) => {
      if (url == '/api/current_user')
        return Promise.resolve({
          data: mockUser(),
        })
      if (url == '/api/current_user/summary')
        return Promise.resolve({
          data: mockSummary(),
        })
      if (url == '/api/strava_auth')
        return Promise.resolve({
          data: { url: 'strava_auth_url' },
        })
      return Promise.reject(new Error('error: ' + url))
    })

    it('loading', () => {
      const tree = renderer.create(<App />).toJSON()
      expect(tree).toMatchSnapshot()
    })
    it('loaded', async () => {
      const root = renderer.create(<App />)
      await actWait()
      expect(root.toJSON()).toMatchSnapshot()
    })
  })
})
