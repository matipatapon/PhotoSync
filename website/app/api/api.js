import {API_ADDRESS} from "../../config.js"

export function getApiUrl(request){
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

function addLeadingZero(x){
    if(x < 10){
        return `0${x}`
    }
    return `${x}`
}

export async function uploadPhoto(file){
    let result = {status: null, creationDate: null}
    let formData = new FormData()
    const date = new Date(file.lastModified)
    const day = addLeadingZero(date.getDate())
    const month = addLeadingZero(date.getMonth() + 1)
    const year = `${date.getFullYear()}`
    const hour = addLeadingZero(date.getHours())
    const minutes = addLeadingZero(date.getMinutes())
    const seconds = addLeadingZero(date.getSeconds())
    const dateStr = `${year}.${month}.${day} ${hour}:${minutes}:${seconds}`
    formData.append("file", file)
    formData.append("filename", file.name)
    formData.append("modification_date", dateStr)
    const token = sessionStorage.getItem("Authorization")
    if(token === null){
        result.status = "NOT_LOGGED_IN"
        return result
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
        if(response.status === 200 || response.status === 201){
             let fileDataJSON = await response.text()
             let fileData = JSON.parse(fileDataJSON)
             result.creationDate = fileData.creation_date
        }
        if(response.status === 200){
            result.status = "SUCCESS"
            return result
        }
        if(response.status === 201){
            result.status = "ALREADY_EXISTS"
            return result
        }
        if(response.status === 401){
            result.status = "UNSUPPORTED"
            return result
        }
        if(response.status === 403){
            result.status = "TOKEN_EXPIRED"
            return result
        }
    } catch(e){}
    result.status = "ERROR"
    return result
}

export async function getFileData(date){
    const token = sessionStorage.getItem("Authorization")
    if(token === null){
        return {status: "NOT_LOGGED_IN", fileData: null}
    }
    try{
        let response = await fetch(
            `${getApiUrl("file_data")}?${new URLSearchParams({date: date})}`,
            {
                method: "GET",
                headers: {
                    "Authorization": token,
                },
            }
        )
        if(response.status === 200){
             let fileData = await response.text()
             return {status: "SUCCESS", fileData: JSON.parse(fileData)}
        }
    } catch(e){}
    return {status: "ERROR", fileData: null}
}

export async function getFile(id){
    const token = sessionStorage.getItem("Authorization")
    if(token === null){
        return {status: "NOT_LOGGED_IN", url: null}
    }

    try{
        let response = await fetch(
            `${getApiUrl("file")}?${new URLSearchParams({id:id})}`,
            {
                method: "GET",
                headers: {
                    "Authorization": token,
                },
            }
        )
        if(response.status === 200){
            let file = await response.blob()
            return {status: "SUCCESS", url: URL.createObjectURL(file)}
        }
    } catch(e){}
    return {status: "ERROR", url: null}
}

export async function removeFile(id){
    const token = sessionStorage.getItem("Authorization")
    if(token === null){
        return "NOT_LOGGED_IN"
    }

    try{
        let response = await fetch(
            `${getApiUrl("file")}?${new URLSearchParams({id:id})}`,
            {
                method: "DELETE",
                headers: {
                    "Authorization": token,
                },
            }
        )
        if(response.status === 200){
            return "SUCCESS"
        }
    } catch(e){}
    return "ERROR"
}
