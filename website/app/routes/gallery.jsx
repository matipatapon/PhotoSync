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
    let [status, setStatus] = useState("SHRINK")
    const imagesPerRow = 6
    const imageHeight = 400;
    const loadedRowsCount = 20;
    const rowsDiffCount = 5;
    let fileData = useRef([])
    let offset = useRef(null)
    let tilesOffset = useRef(0)
    let scrollableHeight = useRef(loadedRowsCount * imageHeight + 100);
    let previousTopScroll = useRef(0)

    let loading = useRef(false)
    useEffect(
        ()=>{
            async function loadFileData() {
                let nextOffset = null
                if(status === "EXPAND"){
                    nextOffset = offset.current + rowsDiffCount
                }
                if(status === "SHRINK"){
                    if(offset.current !== 0){
                        nextOffset = offset.current - rowsDiffCount <= 0 ? 0 : offset.current - rowsDiffCount
                    }
                }

                if(nextOffset !== null)
                {
                    const result = await getFileData((nextOffset) * imagesPerRow, loadedRowsCount * imagesPerRow)
                    const areNewImagesFetched = result.fileData.length >= (loadedRowsCount - rowsDiffCount) * imagesPerRow
                    if(areNewImagesFetched)
                    {
                        offset.current = nextOffset
                        console.log(`Offset ${nextOffset}`)
                        scrollableHeight.current = (loadedRowsCount * imageHeight + offset.current * imageHeight) + 100
                        tilesOffset.current = offset.current * imageHeight
                        fileData.current = result.fileData
                    }
                }
                setStatus("BROWSING")
            }
            if(!loading.current && status !== "BROWSING"){
                loading.current = true
                loadFileData()
                loading.current = false
            }
        },
        [status]
    )
    

    let onScroll = (element) => {
        const direction = previousTopScroll.current > element.currentTarget.scrollTop ? "UP" : "DOWN"
        previousTopScroll.current = element.currentTarget.scrollTop
        const bottomReserve = Math.round(element.currentTarget.scrollHeight - element.currentTarget.offsetHeight - element.currentTarget.scrollTop)
        if(status !== "BROWSING")
        {
            return
        }
        if(direction === "DOWN" && bottomReserve <= rowsDiffCount * imageHeight){
            setStatus("EXPAND")
        }
        if(direction === "UP" && bottomReserve >= rowsDiffCount * imageHeight){
            setStatus("SHRINK")
        }
    }

    return <div className="gallery" onScroll={onScroll} onScrollEnd={onScroll}>
                <div className="scrollable" style={{height: scrollableHeight.current}}>
                    <Files fileData={fileData.current} offset={tilesOffset.current}/>
                </div>
          </div>
}
