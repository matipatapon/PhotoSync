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
    let fileData = useRef([])
    let rowsCount = useRef(6)
    let rowsDiff = useRef(3)
    let offset = useRef(0)
    let tilesOffset = useRef(0)
    let imageHeight = 400;
    let scrollableHeight = useRef(rowsCount.current * imageHeight + 100);
    let topScrollBeforeExpand = useRef(0)

    let loading = useRef(false)
    useEffect(
        ()=>{
            async function loadFileData() {
                if(status === "EXPAND"){
                    const result = await getFileData((offset.current + rowsDiff.current) * 4, rowsCount.current * 4)
                    fileData.current = result.fileData
                    if(result.fileData.length == rowsCount.current * 4){
                        offset.current = offset.current + rowsDiff.current
                        console.log(`Offset ${offset.current}`)
                        scrollableHeight.current = (rowsCount.current * imageHeight + offset.current * imageHeight) + 100
                        tilesOffset.current = offset.current * imageHeight
                    }
                    setStatus("BROWSING")
                }
                if(status === "SHRINK"){
                    const nextOffset = offset.current - rowsDiff.current <= 0 ? 0 : offset.current - rowsDiff.current
                    const result = await getFileData((nextOffset) * 4, rowsCount.current * 4)
                    fileData.current = result.fileData
                    if(result.fileData.length == rowsCount.current * 4){
                        offset.current = nextOffset
                        console.log(`Offset ${nextOffset}`)
                        scrollableHeight.current = (rowsCount.current * imageHeight + offset.current * imageHeight) + 100
                        tilesOffset.current = offset.current * imageHeight
                    }
                    setStatus("BROWSING")

                }
            }
            if(!loading.current){
                loading.current = true
                loadFileData()
                loading.current = false
            }
        },
        [status]
    )
    

    let onScroll = (element) => {
        if(status !== "BROWSING"){
            element.currentTarget.scrollTop = topScrollBeforeExpand.current
        }
        else if(element.currentTarget.offsetHeight + element.currentTarget.scrollTop === element.currentTarget.scrollHeight){
            topScrollBeforeExpand.current = element.currentTarget.scrollTop
            setStatus("EXPAND")
        }
        else if(element.currentTarget.scrollHeight - element.currentTarget.scrollTop - element.currentTarget.offsetHeight > rowsCount.current * imageHeight){
            topScrollBeforeExpand.current = element.currentTarget.scrollTop
            setStatus("SHRINK")
        }
    }

    return <div className="gallery" onScroll={onScroll}>
                <div className="scrollable" style={{height: scrollableHeight.current}}>
                    <Files fileData={fileData.current} offset={tilesOffset.current}/>
                </div>
          </div>
}
