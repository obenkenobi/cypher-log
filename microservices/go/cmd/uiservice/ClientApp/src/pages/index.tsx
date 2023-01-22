import * as React from "react"

import { HeadFC, PageProps } from "gatsby"
import {useCookies} from "react-cookie";

const IndexPage: React.FC<PageProps> = () => {
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

  let authJSX: JSX.Element;

  if (!!profile) {
    authJSX = (
      <>
        <div>
          <h1>profile</h1>
          {JSON.stringify(profile)}
        </div>
        <form action="/auth/logout" method="GET" onSubmit={() => {
          removeCookie("auth-session", {path: "/", domain: window.location.hostname});
        }}>
          <button type="submit">Log out</button>
        </form>
      </>
    );
  } else {
    authJSX = (
      <>
        <a href="/auth/login">SignIn</a>
      </>
    );
  }

  return (
    <main>
      <div>
        <h3>Auth0 Example</h3>
        <p>Zero friction identity infrastructure, built for developers</p>
        <div>
          {authJSX}
        </div>
      </div>
    </main>
  )
}

export default IndexPage

export const Head: HeadFC = () => <title>Home Page</title>