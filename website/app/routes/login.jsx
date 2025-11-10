import { useFetcher , redirect, Link} from "react-router";
import { loginUser } from "../api/api";
import './authentication.css'
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
        return <h2>User or Password is invalid</h2>
    }
    return <h2>{`${status}`}</h2>
}

export default function Login(){
    let fetcher = useFetcher()
    let isIdle = fetcher.state
    let status
    if(fetcher.data !== undefined){
        status = fetcher.data.status
    }
    useLayoutEffect(() => {sessionStorage.setItem("Authorization", null)},[])

    return (
        <>
            <header><Link className="button" to={"/register"}>Register</Link></header>
            <div className="window_container">
                <fetcher.Form className="window" method="post" action="">
                    <input type="text" name="username" placeholder="username" disabled={!isIdle}/>
                    <input type="password" name="password" placeholder="password" disabled={!isIdle}/>
                    <div className="buttons">
                        <button className="button" type="submit" disabled={!isIdle}>Login</button>
                    </div>
                </fetcher.Form>
                <Message status={status}/>
            </div>
        </>
    )
}