import { Outlet, Link} from "react-router";
import './layout.css'

export default function Layout(){
    return(
        <html>
            <head>
                <title>PhotoSync</title>
            </head>
            <body>
                <header><Link to={"/login"}>Login</Link> <Link to={"/register"}>Register</Link> <Link to={"/upload"}>Upload</Link> <Link to={"/gallery"}>Gallery</Link></header>
                <Outlet/>
            </body>
        </html>
    )
}
