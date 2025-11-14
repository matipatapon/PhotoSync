import { INVALID_TOKEN, SUCCESS, ERROR} from "./status"
import {getApiUrl} from './api'

class GET_DATES_RESULT
{
    constructor(){
        this.dates = null
        this.status = null
    }
}

export async function getDates(filtration)
{
    let result = new GET_DATES_RESULT()
    const token = sessionStorage.getItem("Authorization")
    if(token === null)
    {
        result.status = INVALID_TOKEN
        return result
    }

    let url = getApiUrl("dates")
    const year = filtration.year
    const month = filtration.month
    if(filtration.month !== "" && filtration.year !== ""){
        url += `?${new URLSearchParams({year: year, month: month})}`
    } else if (filtration.year !== ""){
        url += `?${new URLSearchParams({year: year})}`
    }

    try{
        let response = await fetch(
            url,
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
