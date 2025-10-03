import "./gallery.css"
import { getDates } from "../api/api"
import { useRef, useEffect, useState, useLayoutEffect} from "react"

function log(string){
    console.log(`[Gallery]: ` + string)
}

function DayData(
    start,
    end,
    height,
    date,
    fileCount,
)
{
    this.start = start
    this.end = end
    this.height = height
    this.date = date
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
    log(`headerHeight => ${headerHeight}`)
    const tileSize = 400
    log(`tileSize => ${tileSize}`)
    const tilesPerRow = Math.floor(galleryWidth / tileSize)
    log(`tilesPerRow => ${tilesPerRow}`)
    for(const date of dates){
        log(`date => ${date.date}`)
        const rowCount = Math.ceil(date.file_count / tilesPerRow)
        log(`rowCount => ${rowCount} = ${date.file_count} / ${tilesPerRow}`)
        const bodyHeight = rowCount * tileSize
        log(`bodyHeight => ${bodyHeight} = ${rowCount} * ${tileSize}`)
        const height = bodyHeight + headerHeight
        const start = lastDayEnd + 1
        const end = start + height
        let day = new DayData(start, end, height, date.date, date.file_count)
        days.push(day)
        lastDayEnd = day.end
    }
    return days
}

export default function Gallery(){
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
            if(day.start <= scrollData.bottom && day.end >= scrollData.top)
            {
                outlet.push(
                    <div key={day.date} style={{
                        height: `${day.height}px`,
                        position: "absolute",
                        transform: `translate(0px, ${day.start}px)`,
                        boxSizing: "border-box",
                        border: "1px solid black",
                        display: "block",
                        width: "100%"
                    }} >{day.date} | {day.fileCount}</div>
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
