import { Link , redirect, useFetcher} from "react-router";
import { registerUser, loginUser} from "../api/api"
import './authentication.css'

export async function clientAction({request}) {
    let formData = await request.formData();
    let username = formData.get("username");
    let password = formData.get("password");
    let password_repeated = formData.get("password_repeated");
    let status
    if(username === ""){
        status = "USERNAME_EMPTY"
    }
    else if(password === ""){
        status = "PASSWORD_EMPTY"
    }
    else if(password !== password_repeated){
        status = "PASSWORD_MISMATCH"
    }else{
        status = await registerUser(username, password)
    }

    if(status === "SUCCESS"){
        return redirect("/login")
    }

    return {username: username, status: status}
}

function Message({status}){
    if(status === undefined){
        return null
    }

    let msg = "Error occured"
    if(status === "USER_ALREADY_EXISTS"){
        msg = "User with given username already exists!"
    }
    if(status === "PASSWORD_MISMATCH"){
        msg = "Password mismatch!"
    }
    if(status === "EMPTY_PASSWORD"){
        msg = "Password cannot be empty!"
    }
    if(status === "EMPTY_USERNAME"){
        msg = "Username cannot be empty!"
    }
    if(status === "WORKING"){
        msg = "Please wait..."
    }
    return <h2>{msg}</h2>
}

export default function Register(){
    let fetcher = useFetcher()
    let username = undefined
    let status = undefined
    let state = fetcher.state
    let isIdle = state === "idle"
    if(fetcher.data !== undefined){
        username = fetcher.data.username
        status = fetcher.data.status
    }
    return(
        <div className="form_container">
            <fetcher.Form className="form" method="post" action="">
                <input type="text" name="username" disabled={!isIdle}/>
                <input type="password" name="password" disabled={!isIdle}/>
                <input type="password" name="password_repeated" disabled={!isIdle}/>
                <button type="submit" disabled={!isIdle}>Register</button>
            </fetcher.Form>
            <Message status={isIdle ? status : "WORKING"}/>
            <h3>Already have account? You can login <Link to={"/login"}>here</Link></h3>
        </div>
    )
}

