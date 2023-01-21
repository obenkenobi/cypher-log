import * as React from "react"

import { HeadFC, PageProps } from "gatsby"

const CounterPage: React.FC<PageProps> = () => {
  const [count, setCount] = React.useState(0)

  return (
    <main>
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
  )
}

export default CounterPage

export const Head: HeadFC = () => <title>Count Page</title>