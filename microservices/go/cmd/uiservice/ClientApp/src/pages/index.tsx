import * as React from "react"

import { HeadFC, PageProps } from "gatsby"

const IndexPage: React.FC<PageProps> = () => {
  return (
    <main>
      <div>
        <h3>Auth0 Example</h3>
        <p>Zero friction identity infrastructure, built for developers</p>
        <a href="/auth/login">SignIn</a>
      </div>
    </main>
  )
}

export default IndexPage

export const Head: HeadFC = () => <title>Home Page</title>