import "./gallery.css"
import { getDates, getFileData } from "../api/api"
import { useRef, useEffect, useState, useLayoutEffect} from "react"

const SPACE_BETWEEN_DAYS = 10

function log(string){
    console.log(`[Gallery]: ` + string)
}

function DayData(
    start,
    end,
    headerHeight,
    height,
    tileSize,
    date,
    fileCount,
)
{
    this.start = start
    this.end = end
    this.headerHeight = headerHeight
    this.height = height
    this.date = date
    this.tileSize = tileSize
    this.fileCount = fileCount
}

function ScrollData(
    top,
    bottom
)
{
    this.top = top
    this.bottom = bottom
}

function calculateDaysData(dates, galleryWidth){
    let days = []
    let lastDayEnd = -1
    const headerHeight = 50
    const maxTileSize = 400
    const minTileSize = 200
    let tileSize = galleryWidth / 5
    if (tileSize > maxTileSize) {
        tileSize = maxTileSize
    }
    if (tileSize < minTileSize){
        tileSize = minTileSize
    }
    const tilesPerRow = Math.floor(galleryWidth / tileSize)
    for(const date of dates){
        const rowCount = Math.ceil(date.file_count / tilesPerRow)
        const bodyHeight = rowCount * tileSize
        const height = bodyHeight + headerHeight + SPACE_BETWEEN_DAYS
        const start = lastDayEnd + 1
        const end = start + height
        let day = new DayData(
            start,
            end,
            headerHeight,
            height,
            tileSize,
            date.date,
            date.file_count)
        days.push(day)
        lastDayEnd = day.end
    }
    return days
}

function Tile({fileData, size}){
    return  <div className="tile" style={{width: `${size}px`, height: `${size}px`}}>
                <div className="content">
                        <img src={`data:image/jpg;base64, ${fileData.thumbnail}`}/>
                </div>
            </div>
}

function Day({day}){
    let [fileData, setFileData] = useState([])
    useEffect(
        ()=>{
            async function fun() {
                const result = await getFileData(day.date)
                const fd = result.fileData
                setFileData(fd) 
            }
            fun()
        },
        []
    )

    let tiles = []
    for(const fd of fileData){
        tiles.push(<Tile key={fd.id} fileData={fd} size={day.tileSize}/>)
    }
    return  <div className="day" style={{height: `${day.height}px`, transform: `translate(0px, ${day.start}px)`}}>
                <div style={{height: `${day.headerHeight}px`}} className="header">
                    {day.date}
                </div>
                {tiles}
            </div>
}

export default function Gallery(){
    const preloadMargin = 1000
    let gallery = useRef(null)
    let [galleryWidth, setGalleryWidth] = useState(null)
    let [dates, setDates] = useState(null)
    let [days, setDays] = useState(null)
    let [scrollData, setScrollData] = useState(new ScrollData(0, window.innerHeight))
    const resizeObserver = new ResizeObserver((entries) => {
        for (const entry of entries) {
            if (entry.contentBoxSize) {
                log(`galleryWidth => ${entry.contentBoxSize[0].inlineSize}`)
                setGalleryWidth(entry.contentBoxSize[0].inlineSize)
            }
        }
    })

    useEffect(
        () => {
            resizeObserver.observe(gallery.current)
            async function fun(){
                const result = await getDates()
                setDates(result.result)
            }
            fun()
        },[])

    useLayoutEffect(
        () => {
            if(galleryWidth === null){
                log("no galleryWidth, skipping effect")
                return
            }
            if(dates === null){
                log("no dates, skipping effect")
                return
            }
            async function fun(){
                setDays(calculateDaysData(dates, galleryWidth))
            }
            fun()
        },[galleryWidth, dates]
    )

    let outlet = null
    if(days === null){
        outlet = <h2>Loading...</h2>
    } else {
        outlet = []
        let totalHeight = 0
        for(let i = 0 ; i < days.length ; i++)
        {
            const day = days[i]
            totalHeight += day.height
            if(day.start - preloadMargin <= scrollData.bottom && day.end + preloadMargin >= scrollData.top)
            {
                outlet.push(
                    <Day key={day.date} day={day}/>
                )
            }
        }
        outlet.push(<div key={totalHeight} style={{height: `${totalHeight}px`}}></div>)
    }

    function scroll(e){
        let gallery = e.currentTarget
        const scrollTop = gallery.scrollTop
        const scrollBottom = scrollTop + gallery.offsetHeight
        log(`scrollTop => ${gallery.scrollTop} | scrollBottom => ${scrollBottom}`)
        setScrollData(new ScrollData(scrollTop, scrollBottom))
    }

    return <div ref={gallery} className="gallery" onScroll={scroll}>
                {outlet}
           </div>
}
