import { Link } from "react-router";
import './authentication.css'

export default function Index(){
    return <>
            <Link className="link" to={"/registration"}>
                Register
            </Link>
            <Link className="link" to={"/login"}>
                Login
            </Link>
        </>
}
