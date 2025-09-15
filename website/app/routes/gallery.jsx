import { getFileData} from "../api/api"
import Tile from "../components/tile.jsx"
import "./gallery.css"
import { useState, useEffect, useRef} from 'react'

function Files({fileData, offset}){
    let result = []
    for(let i = 0 ; i < fileData.length ; i++){
        result.push(<Tile key={fileData[i].id} fileData={fileData[i]} offset={offset}/>)
    }
    return result
}

export default function Gallery(){
    let [state, setState] = useState("LOADING")

    const rowCount = 10
    const loadRowCount = 2
    let tilesOffset = useRef(0)
    let fileData = useRef([])
    let imageOffsetCount = useRef(0)
    let tilePerRowCount = useRef(0)
    let paddingHeight = useRef(0)
    let tileHeight = useRef(0)
    let gallery = useRef(null)
    let isBottom = useRef(false)
    let isTop = useRef(false)

    useEffect(
        ()=>{
            async function load(){
                if(state === "LOADING")
                {
                    const result = await getFileData(0, 40)
                    fileData.current = result.fileData
                }
                else if(state === "RELOADING")
                {
                    const result = await getFileData(imageOffsetCount.current, tilePerRowCount.current * rowCount)
                    if(result.fileData.length > (rowCount - loadRowCount) * tilePerRowCount.current)
                    {
                        const rowOffset = Math.floor(imageOffsetCount.current / tilePerRowCount.current)
                        paddingHeight.current = tileHeight.current * rowOffset
                        tilesOffset.current = paddingHeight.current
                        fileData.current = result.fileData
                        isBottom.current = false
                        isTop.current = false
                    }
                }
                setState("BROWSING")
            }
            
            if(state !== "BROWSING")
            {
                load()
            }
        },
        [state]
    )

    let onScroll = (element) => {
        if(state === "BROWSING")
        {
            const scrollBottom = element.currentTarget.scrollHeight - element.currentTarget.scrollTop - element.currentTarget.offsetHeight
            const tile = element.currentTarget.querySelector(".tile")
            if(tile === null){
                return
            }
            tileHeight.current = tile.offsetHeight
            const tileWidth = tile.offsetWidth
            tilePerRowCount.current = Math.floor(element.currentTarget.offsetWidth / tileWidth)
            if(!isBottom.current && scrollBottom < loadRowCount * tileHeight.current)
            {
                imageOffsetCount.current += loadRowCount * tilePerRowCount.current
                isBottom.current = true
                setState("RELOADING")
            }
            
            const firstTilePosition = tile.getBoundingClientRect().top - gallery.current.getBoundingClientRect().top
            if(!isTop.current && firstTilePosition * -1 < loadRowCount * tileHeight.current)
            {
                imageOffsetCount.current -= loadRowCount * tilePerRowCount.current
                imageOffsetCount.current = imageOffsetCount.current > 0 ? imageOffsetCount.current : 0
                isTop.current = true
                setState("RELOADING")
            }
        }
    }

    return <div className="gallery" ref={gallery} onScroll={onScroll} onScrollEnd={onScroll}>
                <Files fileData={fileData.current} offset={tilesOffset.current}/>
                <div className="scrollable" style={{height: paddingHeight, display: "none"}}></div>
          </div>
}
