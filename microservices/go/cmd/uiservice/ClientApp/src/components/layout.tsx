import React from "react"

interface Props {
  children: React.ReactNode
}

// Todo: implement csrf, user sign in flow, & navbar
const Layout: React.FunctionComponent<Props> = (props: Props) => {
  return <>
    {props.children}
  </>
}

export default Layout