import { getFileData } from "../api/api"
import Tile from "../components/tile.jsx"
import "./gallery.css"

export async function clientLoader({params}){
    const result = await getFileData(0, 100)
    return result
}

function Files({fileData}){
    let result = []
    console.log(fileData)
    for(let i = 0 ; i < fileData.length ; i++){
        result.push(<Tile key={i} fileData={fileData[i]}/>)
    }
    return result
}

export default function Gallery({
    loaderData,
}){
    console.log(loaderData.fileData)
    return <div className="gallery" onScroll={(element) => { console.log(element.currentTarget.scrollTop ) }}>
                <Files fileData={loaderData.fileData}/>
          </div>
}
