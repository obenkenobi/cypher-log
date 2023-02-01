import * as React from "react"

import {HeadFC, Link, PageProps} from "gatsby"
import Layout from "../components/layout";
import { Button } from "flowbite-react/lib/esm/components/Button";

const CounterPage: React.FC<PageProps> = () => {
  const [count, setCount] = React.useState(0)

  return (
    <Layout>
      <main>
        <Link to={"/"}>
          <Button className="btn btn-green">Main Page</Button>
        </Link>
        <h1>Count</h1>
        <div>{count}</div>
        <div>
          <Button className="btn btn-blue my-2" onClick={() => setCount(count + 1)}>
            Count Up
          </Button>
          <br></br>
          <Button className="btn btn-green my-2" onClick={() => setCount(count - 1)}>
            Count Down
          </Button>
        </div>
      </main>
    </Layout>
  )
}

export default CounterPage

export const Head: HeadFC = () => <title>Count Page</title>