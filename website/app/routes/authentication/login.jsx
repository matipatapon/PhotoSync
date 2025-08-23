import { useFetcher , redirect} from "react-router";
import { loginUser } from "../../api/api";
import './authentication.css'

export async function clientAction({request}) {
    let formData = await request.formData();
    let username = formData.get("username");
    let password = formData.get("password");
    
    let status = await loginUser(username, password)

    if(status === "SUCCESS"){
        return redirect("/")
    }

    return {status: status}
}

function Message({status}){
    if(status === "INVALID_USER_OR_PASSWORD"){
        return <h2>User or Password is invalid</h2>
    }
    if(status === "ERROR"){
        return <h2>Error occured</h2>
    }
}

export default function Login(){
    let fetcher = useFetcher()
    let isIdle = fetcher.state
    let status
    if(fetcher.data !== undefined){
        status = fetcher.data.status
    }
    return (
        <div className="form">
            <fetcher.Form method="post" action="">
                <label>Username</label>
                <input type="text" name="username" disabled={!isIdle}/>
                <label>Password</label>
                <input type="password" name="password" disabled={!isIdle}/>
                <button type="submit" disabled={!isIdle}>Login</button>
            </fetcher.Form>
            <Message status={status}/>
        </div>
    )
}