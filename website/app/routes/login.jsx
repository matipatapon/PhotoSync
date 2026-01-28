import { useFetcher , redirect, Link} from "react-router";
import { loginUser } from "../api/api";
import { useLayoutEffect } from "react";

export async function clientAction({request}) {
    let formData = await request.formData();
    let username = formData.get("username");
    let password = formData.get("password");
    
    let status = await loginUser(username, password)

    if(status === "SUCCESS"){
        return redirect("/gallery")
    }

    return {status: status}
}

function Message({status}){
    if(status === "INVALID_USER_OR_PASSWORD"){
        return <span className="error">User or password is invalid</span>
    }
    if(status === "ERROR"){
        return <span className="error">Error occured</span>
    }
}

export default function Login(){
    let fetcher = useFetcher()
    let isIdle = fetcher.state === "idle"
    let status
    if(fetcher.data !== undefined && isIdle){
        status = fetcher.data.status
    }
    useLayoutEffect(() => {sessionStorage.setItem("Authorization", null)},[])

    return (
        <>
            <header><Link className="button" to={"/register"}>Register</Link></header>
            <div className="window_container">
                <fetcher.Form className="pop_up_window" method="post" action="">
                    <h1>Login</h1>
                    <input type="text" name="username" placeholder="username" disabled={!isIdle}/>
                    <input type="password" name="password" placeholder="password" disabled={!isIdle}/>
                    <Message status={status}/>
                    <button className="button" type="submit" disabled={!isIdle}>Login</button>
                </fetcher.Form>
            </div>
        </>
    )
}