import { Link , useFetcher} from "react-router";
import { registerUser } from "../api/register"
import './registration.css'

export async function clientAction({request}) {
    let formData = await request.formData();
    let username = formData.get("username");
    let password = formData.get("password");
    let status = await registerUser(username, password)
    return {username: username, status: status}
}

function Message({status}){
    if(status === undefined){
        return null
    }
    if(status === 401){
        return <h2>User already exists</h2>
    }
    return <h2>Something went wrong</h2>
}

function RegistrationForm({fetcher, status, state}){
    if(state !== "idle"){
        return(
            <fetcher.Form method="post" action="">
                <input type="text" name="username" disabled="true"/>
                <input type="password" name="password" disabled="true"/>
                <button type="submit" disabled="true">Please wait</button>
            </fetcher.Form>
        )
    }
    return (
        <fetcher.Form method="post" action="">
            <label>Username</label>
            <input type="text" name="username"/>
            <label>Password</label>
            <input type="password" name="password"/>
            <button type="submit" disabled={state !== "idle"}>Register</button>
            <Message status={status}/>
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
    if(status === 200){
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

