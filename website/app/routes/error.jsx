import "./error.css"
import { Link} from "react-router";

export default function Error(){
    return  <div className="error">
                <div className="content">
                    <h1>Error occured</h1>
                    <Link className="button" to={"/login"}>Ok</Link>
                </div>
            </div>
}
