class GET_DATES_RESULT
{
    constructor(){
        this.dates = null
        this.error = null
    }
}

export function get_dates()
{
    let result = new GET_DATES_RESULT()
    const token = sessionStorage.getItem("Authorization")
    if(token === null)
    {
        
    }
}
