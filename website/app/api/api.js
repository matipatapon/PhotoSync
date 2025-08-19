export async function registerUser(username, password){
    await new Promise(r => setTimeout(r, 2000));
    try{
        let response = await fetch(
            "http://localhost:8080/v1/register",
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
            "http://localhost:8080/v1/login",
            {
                method: "POST",
                body: JSON.stringify({username: username, password: password}),
            }
        )

        if(response.status === 200){
             let body = await response.text()
             console.log(body)
        }
        return response.status
        
    } catch(e){
        return 500
    }
}
