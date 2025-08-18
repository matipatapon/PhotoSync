import { Form } from "react-router";
import {useState} from "react"
import { registerUser } from "../api/register"

export async function clientAction({request}) {
    let formData = await request.formData();
    let username = formData.get("username");
    let password = formData.get("password");
    return await registerUser(username, password)
}

function onClick(){

}

export default function Register(
    {actionData,}
){
    console.log(actionData)
    if(actionData === undefined){
        return <Form method="post" action="?index">
            <input type="text" name="username"/>
            <input type="password" name="password"/>
            <button type="submit">"Submit"</button>
        </Form>
    }
    if(actionData === "ERROR"){
        return <h1>"Something went wrong :("</h1>
    }
    
}

