import * as React from 'react'
import renderer from 'react-test-renderer'
import { SummaryCard, Props } from '../../src/components/Summary'
import { mockSummary } from '../mock'

describe('src/components/Summary', () => {
  it('snapshot', () => {
    const props: Props = {
      summary: mockSummary(),
    }
    const tree = renderer.create(<SummaryCard {...props} />).toJSON()
    expect(tree).toMatchSnapshot()
  })
})
