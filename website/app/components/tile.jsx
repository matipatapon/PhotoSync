import { useEffect, useState, useRef} from "react"
import { getFile } from "../api/api"
import "./tile.css"

export default function Tile({fileData}){
    let url = useRef(null)
    let [status, setStatus] = useState("LOADING")

    useEffect(
        ()=>{
            async function loadImage(){
                const result = await getFile(fileData.id)
                if(result.status === "ERROR"){
                    setStatus("ERROR")
                    return
                }
                url.current = result.url
                setStatus("FINISHED")
            }
            loadImage()
        },
        []
    )

    let outlet = null
    if(status == "LOADING"){
        outlet = <h1>Please wait...</h1>
    }
    if(status == "ERROR"){
        outlet = <h1>ERROR</h1>
    }
    if(status == "FINISHED"){
        outlet =  <img src={url.current}/>
    }
    return <div className="tile">
                {outlet}
           </div>
}