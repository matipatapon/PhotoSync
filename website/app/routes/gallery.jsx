import "./gallery.css"
import { getDates, getFileData } from "../api/api"
import { useRef, useEffect, useState, useLayoutEffect} from "react"

const DATE_HEIGHT = 100
const EMPTY_SPACE_AT_THE_END_HEIGHT = 50
const MAX_TILE_SIZE = 400
const MIN_TILE_SIZE = 150
const LOAD_MARGIN = 1000

function log(string){
    console.log(`[Gallery]: ` + string)
}

class ElementData
{
    constructor(start, end, height){
        this.start = start
        this.end = end
        this.height = height
    }
}

class DayData extends ElementData
{
    constructor(start, end, height, date){
        super(start, end, height)
        this.date = date
    }
}

class TextData extends ElementData
{
    constructor(start, end, height, text){
        super(start, end, height)
        this.text = text
    }
}

function ScrollData(
    top,
    bottom
)
{
    this.top = top
    this.bottom = bottom
}

function calculateTileSize(containerWidth){
    let tileSize = containerWidth / 5
    if (tileSize > MAX_TILE_SIZE) {
        tileSize = MAX_TILE_SIZE
    }
    if (tileSize < MIN_TILE_SIZE){
        tileSize = MIN_TILE_SIZE
    }
    return Math.floor(tileSize)
}

function createElements(dates, containerWidth, tileSize){
    let elements = []
    let lastDayEnd = -1
    const tilesPerRow = Math.floor(containerWidth / tileSize)
    log(`tilesPerRow{${tilesPerRow}} = ${containerWidth} / ${tileSize}`)
    for(const date of dates){
        let start = lastDayEnd + 1
        let end = start + DATE_HEIGHT
        elements.push(new TextData(start, end, DATE_HEIGHT, date.date))
    
        const rowCount = Math.ceil(date.file_count / tilesPerRow)
        const height = rowCount * tileSize
        start = end + 1
        end = start + height
        elements.push(new DayData(start, end, height, date.date))
        lastDayEnd = end
    }

    let start = lastDayEnd + 1
    let end = start + EMPTY_SPACE_AT_THE_END_HEIGHT
    elements.push(new TextData(start, end, EMPTY_SPACE_AT_THE_END_HEIGHT, ""))
    return elements
}

function Tile({fileData, size}){
    return  <div className="tile" style={{width: `${size}px`, height: `${size}px`}}>
                <div className="content">
                        <img src={`data:image/jpg;base64, ${fileData.thumbnail}`}/>
                </div>
            </div>
}

function Day({day, tileSize}){
    let [fileData, setFileData] = useState([])
    useEffect(
        ()=>{
            let abort = false
            async function fun() {
                const result = await getFileData(day.date)
                if(!abort)
                {
                    const fd = result.fileData
                    setFileData(fd) 

                }
            }
            fun()
            return () => abort = true
        },
        []
    )

    let tiles = []
    for(const fd of fileData){
        tiles.push(<Tile key={fd.id} fileData={fd} size={tileSize}/>)
    }
    return  <div className="day" style={{height: `${day.height}px`, transform: `translate(0px, ${day.start}px)`}}>
                {tiles}
            </div>
}

function Text({data}){
    return <div className="text" style={{height: `${data.height}px`, transform: `translate(0px, ${data.start}px)`}}>
                <div className="content">{data.text}</div>
            </div>
}

export default function Gallery(){
    let tileSize = useRef(null)
    let gallery = useRef(null)
    let container = useRef(null)
    let [containerWidth, setContainerWidth] = useState(null)
    let [dates, setDates] = useState(null)
    let [elements, setElements] = useState(null)
    let [scrollData, setScrollData] = useState(new ScrollData(0, window.innerHeight))
    const resizeObserver = new ResizeObserver((entries) => {
        for (const entry of entries) {
            if (entry.contentBoxSize) {
                console.log(`containerWidth => ${entry.contentBoxSize[0].inlineSize}`)
                setContainerWidth(entry.contentBoxSize[0].inlineSize)
            }
        }
    })

    useEffect(
        () => {
            let abort = false
            resizeObserver.observe(container.current)
            async function fun(){
                const result = await getDates()
                if(!abort)
                {
                    setDates(result.result)
                }
                
            }
            fun()
            return () => abort = true
        },[])

    useLayoutEffect(
        () => {
            if(containerWidth === null){
                log("no containerWidth, skipping effect")
                return
            }
            if(dates === null){
                log("no dates, skipping effect")
                return
            }
            tileSize.current = calculateTileSize(containerWidth)
            gallery.current.scrollTop = 0
            setElements(createElements(dates, containerWidth, tileSize.current))
        },[containerWidth, dates]
    )

    let outlet = null
    if(elements === null)
    {
        outlet = <h2>Loading...</h2>
    }
    else
    {
        outlet = []
        let totalHeight = 0
        let newScrollAnchor = null
        for(let i = 0 ; i < elements.length ; i++)
        {
            const element = elements[i]
            totalHeight += element.height
            // TODO sprawdź czy dobrze to porównujesz + dodaj anchora aby ustawić odpowiedni scroll podczas resize ^^
            if(element.start - LOAD_MARGIN <= scrollData.bottom && element.end + LOAD_MARGIN >= scrollData.top)
            {
                if(element instanceof DayData)
                {
                    outlet.push(<Day key={element.date} day={element} tileSize={tileSize.current}/>)
                }
                else if(element instanceof TextData)
                {
                    outlet.push(<Text key={element.start} data={element}/>)
                }
            }
        }
        outlet.push(<div key={totalHeight} style={{height: `${totalHeight}px`}}></div>)
    }

    function scroll(e){
        let gallery = e.currentTarget
        const scrollTop = gallery.scrollTop
        const scrollBottom = scrollTop + gallery.offsetHeight
        setScrollData(new ScrollData(scrollTop, scrollBottom))
    }

    return <div ref={gallery} className="gallery" onScroll={scroll}>
                <div ref={container} className="container">
                    {outlet}
                </div>
           </div>
}
