import "./gallery.css"
import { getFileData, getFile } from "../api/api"
import { getDates } from "../api/get_dates"
import { SUCCESS } from "../api/status"
import { useRef, useEffect, useState, useLayoutEffect} from "react"
import { useNavigate, Link} from "react-router"

const DATE_HEIGHT = 70
const EMPTY_SPACE_AT_THE_END_HEIGHT = 50
const LOAD_MARGIN = 1000
const FILE_DATA_LOAD_DELAY_MS = 100

function log(string){
    console.log(`[Gallery]: ` + string)
}

class ElementData
{
    constructor(start, height){
        this.start = start
        this.height = height
    }
}

class DayData extends ElementData
{
    constructor(start, height, date){
        super(start, height)
        this.date = date
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

        const rowCount = Math.ceil(date.file_count / tilesPerRow)
        const height = rowCount * tileSize.current + DATE_HEIGHT
        let end = start + height
        elements.push(new DayData(start, height, date.date))
        lastDayEnd = end
    }

    let start = lastDayEnd + 1
    elements.push(new ElementData(start, EMPTY_SPACE_AT_THE_END_HEIGHT))
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

function Tile({fileData, size, setFocusedFileData}){
    function onClick(){
        setFocusedFileData(fileData)
    }
    return  <div className="tile" style={{width: `${size}px`, height: `${size}px`}} onClick={onClick}>
                <div className="content">
                        <img src={`data:image/jpg;base64, ${fileData.thumbnail}`}/>
                </div>
            </div>
}

function Day({day, tileSize, setFocusedFileData}){
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
        tiles.push(<Tile key={fd.id} fileData={fd} size={tileSize} setFocusedFileData={setFocusedFileData}/>)
    }
    return  <div className="day" style={{height: `${day.height}px`, transform: `translate(0px, ${day.start}px)`}}>
                <div className="text" style={{height: `${DATE_HEIGHT}px`}}>
                    <div className="content">{day.date}</div>
                </div>
                {tiles}
            </div>
}

function FocusedFile({focusedFileData, focusedFileUrl, setFocusedFileData}){
    let [showInfo, setShowInfo] = useState(false)
    let [showConfirmation, setShowConfirmation] = useState(false)
    let toggleShowInfo = () => { setShowInfo(!showInfo) }

    if(focusedFileData === null || focusedFileUrl === null){
        if(showInfo)
        {
            toggleShowInfo()
        }
        return null
    }

    function exit(){
        setFocusedFileData(null)
    }

    function info(){
        toggleShowInfo()
    }

    function del(){
        setShowConfirmation(true)
    }

    let removalConfirmationPopUp = showConfirmation ? <div className="removal_confirmation_container">
            <div className="removal_confirmation"><h2>Are you sure?</h2>
            <div className="button">Yes</div>
            <div className="button" onClick={()=>{setShowConfirmation(false)}}>No</div>
            </div>
        </div>
        : null

    let descriptionClassName = showInfo ? "description" : "description hide"
    return <>
                {removalConfirmationPopUp}
                <div className="focused_file_container">
                    <div className="file">
                        <img src={focusedFileUrl}/>
                    </div>
                    <div className="exit button" onClick={exit}>X</div>
                    <div className="info button" onClick={info}>I</div>
                    <div className="del button" onClick={del}>D</div>
                    <div className="description_container">
                        <div className={descriptionClassName}>
                            <h1>{focusedFileData.filename}</h1>
                            <h1>{focusedFileData.creation_date}</h1>
                            <h1>{focusedFileData.mime_type}</h1>
                        </div>
                    </div>
                </div>
            </>
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
    let [focusedFileData, setFocusedFileData] = useState(null)
    let [focusedFileUrl, setFocusedFileUrl] = useState(null)
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

    useEffect(
        ()=>{
            let abort = false
            async function fun() {
                if(focusedFileData === null)
                {
                    setFocusedFileUrl(null)
                    return
                }

                const result = await getFile(focusedFileData.id)
                if(!abort)
                {
                    if(result.status !== SUCCESS)
                    {
                        navigate("/error")
                        return
                    }
                    setFocusedFileUrl(result.url)
                }
            }
            fun()
            return () => abort = true
        },[focusedFileData]
    )

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
            if(element.start - LOAD_MARGIN <= scrollData.bottom && element.start + element.height + LOAD_MARGIN >= scrollData.top)
            {
                if(element instanceof DayData)
                {
                    outlet.push(<Day key={element.date} day={element} tileSize={tileSize.current} setFocusedFileData={setFocusedFileData}/>)
                }
                else
                {
                    outlet.push(<div key={element.start} style={{height: `${EMPTY_SPACE_AT_THE_END_HEIGHT}px`}}/>)
                }
                if(newScrollAnchor === null && element.start <= scrollData.bottom && element.start + element.height >= scrollData.top)
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
                <FocusedFile focusedFileData={focusedFileData} setFocusedFileData={setFocusedFileData} focusedFileUrl={focusedFileUrl}/>
                <header><Link className="button" to={"/upload"}>Upload</Link><Link className="button" to={"/login"}>Logout</Link></header>
                <div ref={gallery} className="gallery" onScroll={scroll}>
                    <div ref={content} className="content">
                        {outlet}
                    </div>
                </div>
           </div>
}
