import "./tile.css"

export default function Tile({fileData}){
    return <div className="tile">
                <img src={`data:image/jpg;base64, ${fileData.thumbnail}`}/>
           </div>
}