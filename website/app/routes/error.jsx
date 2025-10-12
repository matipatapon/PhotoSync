import { Link} from "react-router";

export default function Error(){
    return  <>
                <header></header>
                <div className="window_container">
                    <div className="window">
                        <h1>Error occured</h1>
                        <div className="buttons">
                            <Link className="button" to={"/login"}>Ok</Link>
                        </div>
                    </div>
                </div>
            </>
}
