import { redirect } from "react-router"

export function clientLoader(){
    return redirect("/login")
}

export default function(){
    return <div></div>
}