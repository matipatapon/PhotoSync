export async function registerUser(username, password){
    let response = await fetch(
        "http://localhost:8080/v1/register",
        {
            method: "POST",
            body: JSON.stringify({username: username, password: password})
        }
    )
    if(response.status === 200){
        return "SUCCESS"
    }
    else if(response.status === 401){
        return "USER_ALREADY_EXISTS"
    }
    else{
        return "ERROR"
    }
}
