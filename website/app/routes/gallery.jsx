import "./gallery.css"
import { getFileData } from "../api/api"
import { getDates } from "../api/get_dates"
import { SUCCESS } from "../api/status"
import { useRef, useEffect, useState, useLayoutEffect} from "react"
import { useNavigate, Link} from "react-router"

const DATE_HEIGHT = 50
const EMPTY_SPACE_AT_THE_END_HEIGHT = 50
const LOAD_MARGIN = 1000
const FILE_DATA_LOAD_DELAY_MS = 100

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

function calculateTilesPerRow(containerWidth){
    if(containerWidth > 2000){
        return 7
    }
    if(containerWidth > 1000){
        return 5
    }
    return 3
}

function createElements(dates, containerWidth, tileSize){
    let elements = []
    let lastDayEnd = -1
    const tilesPerRow = calculateTilesPerRow(containerWidth)
    tileSize.current = Math.floor(containerWidth / tilesPerRow)
    for(const date of dates){
        let start = lastDayEnd + 1
        let end = start + DATE_HEIGHT
        elements.push(new TextData(start, end, DATE_HEIGHT, date.date))

        const rowCount = Math.ceil(date.file_count / tilesPerRow)
        const height = rowCount * tileSize.current
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

function alignScrollTop(elements, gallery, anchor)
{
    for(let i = 0 ; i < elements.length ; i++)
    {
        if(anchor === i)
        {
            gallery.scrollTop = elements[i].start
            return
        }
    }
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
            setTimeout(() => {
                async function fun()
                {
                    if(!abort)
                    {
                        const result = await getFileData(day.date)
                        if(!abort)
                        {
                            const fd = result.fileData
                            setFileData(fd) 
                        }
                    }
                }
                fun()
            }, FILE_DATA_LOAD_DELAY_MS)
            return () => abort = true},[]
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
    let content = useRef(null)
    let anchor = useRef(0)
    let [containerWidth, setContainerWidth] = useState(null)
    let [dates, setDates] = useState(null)
    let [elements, setElements] = useState(null)
    let [scrollData, setScrollData] = useState(new ScrollData(0, window.innerHeight))
    let navigate = useNavigate()
    const resizeObserver = new ResizeObserver((entries) => {
        for (const entry of entries) {
            if (entry.contentBoxSize) {
                setContainerWidth(entry.contentBoxSize[0].inlineSize)
            }
        }
    })

    useEffect(
        () => {
            let abort = false
            resizeObserver.observe(content.current)
            async function fun(){
                const result = await getDates()
                if(!abort)
                {
                    if(result.status !== SUCCESS)
                    {
                        navigate("/error")
                    }
                    else
                    {
                        setDates(result.dates)
                    }
                }
            }
            fun()
            return () => abort = true
        },[])

    useLayoutEffect(
        () => {
            if(containerWidth === null || dates === null || gallery.current === null){
                return
            }
            const elements = createElements(dates, containerWidth, tileSize)
            alignScrollTop(elements, gallery.current, anchor.current)
            setElements(elements)
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
                if(newScrollAnchor === null && element.start <= scrollData.bottom && element.end >= scrollData.top)
                {
                    newScrollAnchor = i
                }
            }
        }
        anchor.current = newScrollAnchor
        outlet.push(<div key={totalHeight} style={{height: `${totalHeight}px`}}></div>)
    }

    function scroll(e){
        let gallery = e.currentTarget
        const scrollTop = gallery.scrollTop
        const scrollBottom = scrollTop + gallery.offsetHeight
        setScrollData(new ScrollData(scrollTop, scrollBottom))
    }

    return <div className="gallery_container">
                <header><Link className="button" to={"/upload"}>Upload</Link><Link className="button" to={"/login"}>Logout</Link></header>
                <div ref={gallery} className="gallery" onScroll={scroll}>
                    <div ref={content} className="content">
                        {outlet}
                    </div>
                </div>
           </div>
}
