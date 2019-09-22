import 'regenerator-runtime/runtime'
import * as React from 'react'
import { render } from 'react-dom'
import Link from '@material-ui/core/Link'
import axios from 'axios'
import queryString from 'query-string'
import Grid from '@material-ui/core/Grid'
import Button from '@material-ui/core/Button'
import TextField from '@material-ui/core/TextField'
import Typography from '@material-ui/core/Typography'

const parsed = queryString.parse(location.search)
const token = (parsed.token as string) || localStorage.getItem('token') || ''
axios.defaults.headers.common['Authorization'] = token
// axios.defaults.headers.common['Authorization'] =
//   'eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE1Njk3NjQ4MzIsImlkIjoyfQ.y3yh8NskL3ONhqHA09hu9e3-zHSK6xVxJBFXOXmUAEmoX0dPFqKAF6k-hrlxeKlkVx3Dafjb_tlGYOiO8KqFs7o6uMeJ0le50pwtak_sWlovr_hV97D0DQutj54Z488vMYu_ulNILYoy-rzLQFhbeHqJ0LeqKniADxVKYx5vAYQ6bdslgaTyqKcsJlXkPSEp9sjgiNC7TNCGmRn1Ta0RfcABt-lpgzo-Ayr0_dAWgSC0aJq6AR-UuWszrw5SampaJqpOF0e0woSEotrGl-hOLsHNk9uIBBJkTHH_bDb4vgqakAxwj09bmRGCY756BC48AiHrNYpENrFCcgCBc2qpRw'

const App: React.FC = () => {
  const [loading, setLoading] = React.useState(true)
  const [user, setUser] = React.useState(null)
  const [draftKey, setDraftKey] = React.useState('')
  const [draftMessage, setDraftMessage] = React.useState('')

  const onInputDraftKey = React.useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setDraftKey(e.currentTarget.value)
  }, [])
  const onInputDraftMessage = React.useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setDraftMessage(e.currentTarget.value)
  }, [])
  const onSubmit = React.useCallback(
    (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault()
      axios
        .post('/api/current_user', {
          iftttKey: draftKey,
          iftttMessage: draftMessage,
        })
        .then(res => {
          setUser(res.data)
        })
    },
    [draftKey, draftMessage],
  )

  React.useEffect(() => {
    axios
      .get('/api/current_user')
      .then(res => {
        setUser(res.data)
        setDraftKey(res.data.iftttKey)
        setDraftMessage(res.data.iftttMessage)
        localStorage.setItem('token', token)
        setLoading(false)
      })
      .catch(() => {
        setLoading(false)
      })
  }, [])

  const userBlock = React.useMemo(() => {
    if (user)
      return (
        <form onSubmit={onSubmit}>
          <Grid container>
            <Grid item xs={12}>
              <Typography>ID</Typography>
              <p>{user.id}</p>
            </Grid>
            <Grid item xs={12}>
              <Typography>Strava ID</Typography>
              <p>{user.athleteId}</p>
            </Grid>
            <Grid item xs={12}>
              <Typography>Strava Username</Typography>
              <p>{user.username}</p>
            </Grid>
            <Grid item xs={12}>
              <Typography>IFTTT Key</Typography>
              <TextField variant="outlined" value={draftKey} onChange={onInputDraftKey} />
              <Typography>IFTTT Message</Typography>
              <TextField variant="outlined" value={draftMessage} onChange={onInputDraftMessage} />
            </Grid>
            <Grid item>
              <Button type="submit" variant="contained" color="primary">
                Submit
              </Button>
            </Grid>
          </Grid>
        </form>
      )
    return (
      <div>
        <p>Not Auth</p>
      </div>
    )
  }, [draftKey, draftMessage, onInputDraftKey, onInputDraftMessage, onSubmit, user])

  return (
    <div>
      <div>{loading ? 'loading...' : userBlock}</div>
    </div>
  )
}

render(<App />, document.getElementById('root') as any)
