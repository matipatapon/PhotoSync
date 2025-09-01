import { Outlet } from "react-router";
import './layout.css'

export default function Layout(){
    return(
        <html>
            <head>
                <title>PhotoSync</title>
            </head>
            <body>
                <header>PhotoSync</header>
                <Outlet/>
            </body>
        </html>
    )
}
