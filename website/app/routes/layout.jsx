import { Outlet, Link} from "react-router";
import './layout.css'

export default function Layout(){
    return(
        <html>
            <head>
                <title>PhotoSync</title>
            </head>
            <body>
                <Outlet/>
            </body>
        </html>
    )
}
