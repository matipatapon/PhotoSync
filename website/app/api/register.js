export async function registerUser(username, password){
    try{
        let response = await fetch(
            "http://localhost:8080/v1/register",
            {
                method: "POST",
                body: JSON.stringify({username: username, password: password}),
            }
        )
        return response.status
    } catch(e){
        return 500
    }
}
