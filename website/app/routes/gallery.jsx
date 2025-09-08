import { getFileData, getFile } from "../api/api"
import { useState , useEffect } from "react"

export async function clientLoader({params}){
    const result = await getFileData(0, 100)
    return result
}

function File({data}){
    let [url, setUrl] = useState("")
    useEffect(
        () => {
            async function loadFile(){
                const result = await getFile(data.id)
                setUrl(result.url)
            }
            loadFile()
        }
        ,[]
    )
    return  <div>
                <h1>{data.filename}</h1>
                <h2>{data.creation_date}</h2>
                <img height="200" src={url}/>
            </div>
}

function Files({fileData}){
    
    let result = []
    console.log(fileData)
    for(let i = 0 ; i < fileData.length ; i++){
        result.push(<File key={i} data={fileData[i]}/>)
    }
    return result
}

export default function Gallery({
    loaderData,
}){
    console.log(loaderData.fileData)
    return <Files fileData={loaderData.fileData}/>
}
