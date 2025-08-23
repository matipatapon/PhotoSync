import {API_ADDRESS} from "../../config.js"

function getApiUrl(request){
    return `${API_ADDRESS}/v1/${request}`
}

function validateUsername(username){
    if(username === ""){
        return "EMPTY_USERNAME"
    }
    return "OK"
}

function validatePassword(password){
    if(password === ""){
        return "EMPTY_PASSWORD"
    }
    return "OK"
}

export async function registerUser(username, password){
    // await new Promise(r => setTimeout(r, 2000));
    let userStatus = validateUsername(username)
    if(userStatus !== "OK"){
        return userStatus
    }

    let passwordStatus = validatePassword(password)
    if(passwordStatus !== "OK"){
        return passwordStatus
    }

    try{
        let response = await fetch(
            getApiUrl("register"),
            {
                method: "POST",
                body: JSON.stringify({username: username, password: password}),
            }
        )
        if(response.status === 200){
            return "SUCCESS"
        }
        if(response.status === 401){
            return "USER_ALREADY_EXISTS"
        }
    } catch(e){
    }
    return "ERROR"
}

export async function loginUser(username, password){
    try{
        let response = await fetch(
            getApiUrl("login"),
            {
                method: "POST",
                body: JSON.stringify({username: username, password: password}),
            }
        )
        if(response.status === 200){
             let body = await response.text()
             sessionStorage.setItem("Authorization", body)
             return "SUCCESS"
        }
        if(response.status === 401){
            return "INVALID_USER_OR_PASSWORD"
        }
    } catch(e){
    }
    return "ERROR"
}
