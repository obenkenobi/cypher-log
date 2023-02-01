import { Button } from "flowbite-react/lib/esm/components/Button"
import { Navbar } from "flowbite-react/lib/esm/components/Navbar/Navbar"
import React from "react"

interface Props {
  children: React.ReactNode
}

// Todo: implement csrf, user sign in flow, & navbar
const Layout: React.FunctionComponent<Props> = (props: Props) => {
  // Todo: use your own navbar links
  return <>
    <Navbar
      fluid={true}
      rounded={true}
    >
      <Navbar.Brand href="https://flowbite.com/">
        <img
          src="https://flowbite.com/docs/images/logo.svg"
          className="mr-3 h-6 sm:h-9"
          alt="Flowbite Logo"
        />
        <span className="self-center whitespace-nowrap text-xl font-semibold dark:text-white">
      Flowbite
    </span>
      </Navbar.Brand>
      <div className="flex md:order-2">
        <Button>
          Get started
        </Button>
        <Navbar.Toggle />
      </div>
      <Navbar.Collapse>
        <Navbar.Link
          href="/navbars"
          active={true}
        >
          Home
        </Navbar.Link>
        <Navbar.Link href="/navbars">
          About
        </Navbar.Link>
        <Navbar.Link href="/navbars">
          Services
        </Navbar.Link>
        <Navbar.Link href="/navbars">
          Pricing
        </Navbar.Link>
        <Navbar.Link href="/navbars">
          Contact
        </Navbar.Link>
      </Navbar.Collapse>
    </Navbar>
    {props.children}
  </>
}

export default Layout