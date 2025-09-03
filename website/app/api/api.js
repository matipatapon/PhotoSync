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

export async function uploadPhoto(file){
    let formData = new FormData()
    formData.append("file", file)
    formData.append("filename", file.name)
    formData.append("modification_date", "2025.05.16 16:30:12")
    // TODO FIX DATE

    const token = sessionStorage.getItem("Authorization")
    if(token === null){
        return "NOT_LOGGED_IN"
    }

    try{
        let response = await fetch(
            getApiUrl("upload"),
            {
                method: "POST",
                headers: {
                    "Authorization": token,
                },
                body: formData,
            }
        )
        if(response.status === 200){
             let file_id = await response.text()
             return "SUCCESS"
        }
        if(response.status === 401){
            return "UNSUPPORTED"
        }
        if(response.status === 402){
            return "ALREADY_EXISTS"
        }
        if(response.status === 403){
            return "TOKEN_EXPIRED"
        }
    } catch(e){
    }
    return "ERROR"
}
