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
                <div id="container">
                    <Outlet/>
                </div>
                <footer></footer>
            </body>
        </html>
    )
}
