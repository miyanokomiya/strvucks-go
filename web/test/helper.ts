import { act } from 'react-test-renderer'

// https://github.com/wesbos/waait/blob/master/index.js
export function wait(amount = 0) {
  return new Promise(resolve => setTimeout(resolve, amount))
}

// Use this in your test after mounting if you need just need to let the query finish without updating the wrapper
export async function actWait(amount = 0) {
  await act(async () => {
    await wait(amount)
  })
}

// Use this in your test after mounting if you want the query to finish and update the wrapper
export async function updateWrapper(wrapper: any, amount = 0) {
  await act(async () => {
    await wait(amount)
    wrapper.update()
  })
}
