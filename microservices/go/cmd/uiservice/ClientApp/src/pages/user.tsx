import * as React from "react"

import { HeadFC, PageProps } from "gatsby"
import {useCookies} from "react-cookie";

const UserPage: React.FC<PageProps> = () => {
  const [profile, setProfile] = React.useState<any>()
  const [ , , removeCookie] = useCookies(["auth-session"]);


  React.useEffect(() => {
    const getDataTask = (async () => {
      const res = await fetch("/auth/me")
      const body = await res.json()
      console.log(body)
      setProfile(body)
    })()

    getDataTask
      .catch(e => console.log(e))
  }, [])

  return (
    <main>
      <div>
        <h1>profile</h1>
        {JSON.stringify(profile)}
      </div>
      <form action="/auth/logout" method="GET" id="logout" onSubmit={e => {
        removeCookie("auth-session", {path: "/", domain: window.location.hostname});
      }}>
        <input type="submit" value="Log out"/>
      </form>
    </main>
  )
}

export default UserPage

export const Head: HeadFC = () => <title>User Page</title>