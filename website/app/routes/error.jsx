import { Link} from "react-router";

export default function Error(){
    return  <>
                <header></header>
                <div className="window_container">
                    <div className="pop_up_window">
                        <h1>Something went wrong</h1>
                        <Link className="button" to={"/login"}>Ok</Link>
                    </div>
                </div>
            </>
}
