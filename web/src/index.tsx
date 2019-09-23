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
import Snackbar from '@material-ui/core/Snackbar'
import IconButton from '@material-ui/core/IconButton'
import CloseIcon from '@material-ui/icons/Close'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import { User, Summary } from './types'
import { SummaryCard } from './components/Summary'

const parsed = queryString.parse(location.search)
const token = (parsed.token as string) || localStorage.getItem('token') || ''
axios.defaults.headers.common['Authorization'] = token

const App: React.FC = () => {
  const [loading, setLoading] = React.useState(true)
  const [user, setUser] = React.useState(null as User | null)
  const [summary, setSummary] = React.useState(null as Summary | null)
  const [stravaAuth, setStravaAuth] = React.useState('')
  const [draftKey, setDraftKey] = React.useState('')
  const [draftMessage, setDraftMessage] = React.useState('')
  const [snack, setSnack] = React.useState(false)

  const onCloseSnack = React.useCallback(() => {
    setSnack(false)
  }, [])
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
          setSnack(true)
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

  React.useEffect(() => {
    axios.get('/api/current_user/summary').then(res => {
      setSummary(res.data)
    })
  }, [])

  React.useEffect(() => {
    axios.get('/api/strava_auth').then(res => {
      setStravaAuth(res.data.url)
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
              <TextField
                label="IFTTT Key"
                variant="outlined"
                margin="normal"
                value={draftKey}
                onChange={onInputDraftKey}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                label="IFTTT Message"
                variant="outlined"
                margin="normal"
                value={draftMessage}
                onChange={onInputDraftMessage}
              />
            </Grid>
            <Grid item>
              <Button type="submit" variant="contained" color="primary">
                Submit
              </Button>
            </Grid>
          </Grid>
          <Snackbar
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'left',
            }}
            open={snack}
            autoHideDuration={6000}
            onClose={onCloseSnack}
            ContentProps={{
              'aria-describedby': 'message-id',
            }}
            message={<span>Saved</span>}
            action={[
              <IconButton key="close" aria-label="close" color="inherit" onClick={onCloseSnack}>
                <CloseIcon />
              </IconButton>,
            ]}
          />
        </form>
      )
    return (
      <div>
        <Link href={stravaAuth}>
          <p>Login by Strava</p>
          <img src="/assets/strava.jpg" style={{ width: '120px', height: 'auto' }} />
        </Link>
      </div>
    )
  }, [
    draftKey,
    draftMessage,
    onCloseSnack,
    onInputDraftKey,
    onInputDraftMessage,
    onSubmit,
    snack,
    stravaAuth,
    user,
  ])

  const summaryBlock = React.useMemo(() => {
    if (summary) return <SummaryCard summary={summary} />

    return (
      <Card>
        <CardContent>
          <Typography variant="h5" component="h2">
            No Summary
          </Typography>
        </CardContent>
      </Card>
    )
  }, [summary])

  const mainBlock = React.useMemo(() => {
    if (!loading)
      return (
        <div>
          <div>{userBlock}</div>
          <div>{summaryBlock}</div>
        </div>
      )

    return <div>loading...</div>
  }, [loading, summaryBlock, userBlock])

  return <div>{mainBlock}</div>
}

render(<App />, document.getElementById('root') as any)
