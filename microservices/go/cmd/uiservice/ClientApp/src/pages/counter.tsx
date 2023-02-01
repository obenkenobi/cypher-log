import * as React from "react"

import {HeadFC, Link, PageProps} from "gatsby"
import Layout from "../components/layout";

const CounterPage: React.FC<PageProps> = () => {
  const [count, setCount] = React.useState(0)

  return (
    <Layout>
      <main>
        <Link to={"/"}>
          <button className="btn btn-green">Main Page</button>
        </Link>
        <h1>Count</h1>
        <div>{count}</div>
        <div>
          <button className="btn btn-blue my-2" onClick={() => setCount(count + 1)}>
            Count Up
          </button>
          <br></br>
          <button className="btn btn-green my-2" onClick={() => setCount(count - 1)}>
            Count Down
          </button>
        </div>
      </main>
    </Layout>
  )
}

export default CounterPage

export const Head: HeadFC = () => <title>Count Page</title>