import { INVALID_TOKEN, SUCCESS, ERROR} from "./status"
import {getApiUrl} from './api'

class GET_DATES_RESULT
{
    constructor(){
        this.dates = null
        this.status = null
    }
}

export async function getDates()
{
    let result = new GET_DATES_RESULT()
    const token = sessionStorage.getItem("Authorization")
    if(token === null)
    {
        result.status = INVALID_TOKEN
        return result
    }

    try{
        let response = await fetch(
            `${getApiUrl("dates")}`,
            {
                method: "GET",
                headers: {
                    "Authorization": token,
                },
            }
        )

        if(response.status === 200){
            result.status = SUCCESS
            let dates = await response.text()
            result.dates = JSON.parse(dates)
            return result
        }
        if(response.status === 403){
            result.status = INVALID_TOKEN
            return result
        }

} catch(e){}
    result.status = ERROR
    return result
}
