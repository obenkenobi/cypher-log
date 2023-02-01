import * as React from "react"

import {HeadFC, Link, PageProps} from "gatsby"
import {useCookies} from "react-cookie";

const IndexPage: React.FC<PageProps> = () => {
  const [profile, setProfile] = React.useState<any>()
  const [userSave, setUserSave] = React.useState<UserSaveDto>({
    userName: "", displayName: ""
  })
  const [ cookies, , removeCookie] = useCookies(["XSRF-TOKEN", "session"]);

  const getDataTask = (async () => {
    try {
      const res = await fetch("/api/userservice/v1/user/me")
      const body = await res.json()
      console.log(body)
      if (res.status == 200) {
        setProfile(body)
      }
    } catch (e) {
      console.log(e)
    }
  })

  React.useEffect(() => {
    getDataTask().then()


    const csrfTask = (async () => {
      const res = await fetch("/csrf")
      const body = await res.json()
      console.log(body)
      console.log(cookies["XSRF-TOKEN"])
    })()

    csrfTask.catch(e => console.log(e))
  }, [])

  const handleSubmit = async (e:  React.FormEvent<HTMLFormElement>): Promise<void> => {
    try {
      e.preventDefault()
      const res = await fetch("/api/userservice/v1/user", {
        method: "PUT",
        body: JSON.stringify(userSave),
        headers: {
          'Content-Type': 'application/json',
          "X-XSRF-TOKEN": cookies["XSRF-TOKEN"]
        }
      })
      const body = await res.json()
      console.log(body)
      await getDataTask()
    } catch (e) {
      console.log(e)
    }
  }

  let authJSX: JSX.Element;
  if (!!profile) {
    authJSX = (
      <>
        <div>
          <h1 className="my-1">profile</h1>
          <pre className="my-1">{JSON.stringify(profile, null, "\t")}</pre>
        </div>
        <h1>Update</h1>
        <form onSubmit={handleSubmit}>
          <label>userName
            <input name="userName"
                   type="text"
                   onChange={e => setUserSave({...userSave, userName: e.target.value})}/>
          </label>
          <br/>
          <label>
            displayName
            <input name="displayName"
                 type="text"
                 onChange={e => setUserSave({...userSave, displayName: e.target.value})}/>
          </label>
          <br/>
          <button type="submit" className="btn btn-blue">Submit</button>
        </form>
        <div>--------------------------------------</div>
        <form className="my-3" action="/auth/logout" method="GET" onSubmit={() => {
          removeCookie("session", {path: "/", domain: window.location.hostname});
        }}>
          <button type="submit" className="btn btn-blue">Log out</button>
        </form>
      </>
    );
  } else {
    authJSX = (
      <>
        <a href="/auth/login">
          <button className="btn btn-green">SignIn</button>
        </a>
      </>
    );
  }

  return (
    <main>
      <div>
        <h3 className="my-1">Auth0 Example</h3>
        <p className="my-1">Zero friction identity infrastructure, built for developers</p>
        <div className="my-4">
          <Link to={"/counter"}>
            <button className="btn btn-green">Counter</button>
          </Link>
        </div>
        <div className="my-4">
          {authJSX}
        </div>
      </div>
    </main>
  )
}

export default IndexPage

export const Head: HeadFC = () => <title>Home Page</title>