import { Link , useFetcher} from "react-router";
import { registerUser, loginUser} from "../api/api"
import './registration.css'

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
    if(status === "PASSWORD_EMPTY"){
        msg = "Password cannot be empty!"
    }
    if(status === "USERNAME_EMPTY"){
        msg = "Username cannot be empty!"
    }
    if(status === "WORKING"){
        msg = "Please wait..."
    }
    return <h2>{msg}</h2>
}

function RegistrationForm({fetcher, status, state}){
    let isIdle = state === "idle"
    return (
        <fetcher.Form method="post" action="">
            <label>Username</label>
            <input type="text" name="username" disabled={!isIdle}/>
            <label>Password</label>
            <input type="password" name="password" disabled={!isIdle}/>
            <label>Password</label>
            <input type="password" name="password_repeated" disabled={!isIdle}/>
            <button type="submit" disabled={!isIdle}>Register</button>
            <Message status={isIdle ? status : "WORKING"}/>
        </fetcher.Form>
    )
}

export default function Registration(){
    let fetcher = useFetcher()
    let username = undefined
    let status = undefined
    let state = fetcher.state
    if(fetcher.data !== undefined){
        username = fetcher.data.username
        status = fetcher.data.status
    }
    if(status === "SUCCESS"){
        return(
            <div id="registration_form">
                <h2>{username} was successfully registered!</h2>
                You can now <Link to="/login">login</Link>
            </div>
        )
    }
    return(
        <div id="registration_form">
            <RegistrationForm fetcher={fetcher} status={status} state={state}/>
        </div>
    )

}

