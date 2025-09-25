import { getFileData} from "../api/api"
import Tile from "../components/tile.jsx"
import "./gallery.css"
import { useState, useEffect, useRef} from 'react'

function Files({fileData, offset}){
    let result = []
    for(let i = 0 ; i < fileData.length ; i++){
        result.push(<Tile key={fileData[i].id} fileData={fileData[i]} offset={offset}/>)
    }
    return <div className="file_container">{result}</div>
}

export default function Gallery(){
    let [state, setState] = useState("LOADING")
    const maxRowCountAtTheTime = 20
    const commonRowsBetweenLoadings = 5
    const diffRowCount = maxRowCountAtTheTime - commonRowsBetweenLoadings
    let fileData = useRef([])
    let imageOffsetCount = useRef(0)
    let tilePerRowCount = useRef(0)
    let paddingHeight = useRef(0)
    let tileHeight = useRef(0)
    let padding = useRef(null)
    let isBottom = useRef(false)
    let isTop = useRef(false)
    let previousImageOffsetCount = useRef(0)

    let resizeObserver = new ResizeObserver((entries) => {
        for(const entry of entries){
            if(entry.borderBoxSize){
                
            }
        }
    })

    useEffect(
        ()=>{
            async function load(){
                if(state === "LOADING")
                {
                    const result = await getFileData(0, 20) // you need to know how many tiles should be otherwise it brakes apart
                    fileData.current = result.fileData
                }
                else if(state === "RELOADING")
                {
                    // it doesn't work well after
                    if(padding.current != null){
                        const result = await getFileData(imageOffsetCount.current, tilePerRowCount.current * maxRowCountAtTheTime)
                        if(result.fileData.length > commonRowsBetweenLoadings * tilePerRowCount.current)
                        {
                            const rowOffset = Math.floor(imageOffsetCount.current / tilePerRowCount.current)
                            paddingHeight.current = tileHeight.current * rowOffset
                            fileData.current = result.fileData
                            isBottom.current = false
                            isTop.current = false
                            previousImageOffsetCount.current = imageOffsetCount.current
                            padding.current.style.height = `${paddingHeight.current}px`
                        }
                        else
                        {
                            imageOffsetCount.current = previousImageOffsetCount.current
                        }
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
            const tiles = element.currentTarget.querySelectorAll(".tile")
            if(tiles.length == 0){
                return
            }
            const firstTile = tiles[0]
            const lastTile = tiles[tiles.length - 1]
            const tileWidth = firstTile.offsetWidth
            const galleryWidthWithoutScrollbar = element.currentTarget.querySelector(".file_container").offsetWidth
            tileHeight.current = firstTile.offsetHeight
            tilePerRowCount.current = Math.floor(galleryWidthWithoutScrollbar / tileWidth)

            const lastTilePosition = lastTile.getBoundingClientRect().bottom - element.currentTarget.getBoundingClientRect().bottom
            console.log(lastTilePosition)
            if(isBottom.current === false && lastTilePosition <= tileHeight.current)
            {
                console.log(tilePerRowCount.current)
                imageOffsetCount.current += diffRowCount * tilePerRowCount.current
                isBottom.current = true
                setState("RELOADING")
            }
            
            const firstTilePosition = firstTile.getBoundingClientRect().top - element.currentTarget.getBoundingClientRect().top
            if(isTop.current === false && firstTilePosition * -1 < tileHeight.current * 2)
            {
                imageOffsetCount.current -= diffRowCount * tilePerRowCount.current
                imageOffsetCount.current = imageOffsetCount.current > 0 ? imageOffsetCount.current : 0
                if(previousImageOffsetCount.current == 0 && imageOffsetCount.current == 0){
                    return
                }
                isTop.current = true
                setState("RELOADING")
            }
        }
    }

    return <div className="gallery" onScroll={onScroll} onScrollEnd={onScroll}>
                <div ref={padding} className="padding"></div>
                <Files fileData={fileData.current}/>
          </div>
}
