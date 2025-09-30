import "./tile.css"

export default function Tile({fileData, offset}){
    return <div className="tile" style={{transform: `translate(0px, ${offset}px)`}}>
                <img src={`data:image/jpg;base64, ${fileData.thumbnail}`}/>
           </div>
}