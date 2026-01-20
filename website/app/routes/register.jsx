import { Link , redirect, useFetcher} from "react-router";
import { registerUser } from "../api/api"
import './authentication.css'

export async function clientAction({request}) {
    let formData = await request.formData();
    let username = formData.get("username");
    let password = formData.get("password");
    let password_repeated = formData.get("password_repeated");
    let errorMsg
    if(username === ""){
        errorMsg = "Username cannot be empty"
    }
    else if(password === ""){
        errorMsg = "Password cannot be empty"
    }
    else if(password !== password_repeated){
        errorMsg = "Password mismatch"
    }
    else{
        const status = await registerUser(username, password)
        if(status === "USER_ALREADY_EXISTS"){
            errorMsg = "User with given username already exists"
        }
        else if(status === "SUCCESS"){
            return redirect("/login")
        }
        else{
            errorMsg = "Error occured"
        }
    }

    return {errorMsg: errorMsg}
}

function ErrorMessage({fetcher}){
    if(fetcher.data === undefined || fetcher.data.errorMsg === null){
        return <></>
    }
    return <span className="error">{fetcher.data.errorMsg}</span>
}

export default function Register(){
    let fetcher = useFetcher()
    let isIdle = fetcher.state === "idle"
    return(
            <>
                <header><Link className="button" to={"/login"}>Login</Link></header>
                <div className="window_container">
                    <fetcher.Form className="pop_up_window" method="post" action="">
                        <h1>Register</h1>
                        <input type="text" name="username" placeholder="username" disabled={!isIdle}/>
                        <input type="password" name="password" placeholder="password" disabled={!isIdle}/>
                        <input type="password" name="password_repeated" placeholder="password" disabled={!isIdle}/>
                        <ErrorMessage fetcher={fetcher}/>
                        <button className="button" type="submit" disabled={!isIdle}>Register</button>
                    </fetcher.Form>
                </div>
            </>
    )
}

