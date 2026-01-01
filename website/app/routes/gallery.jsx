import "./gallery.css"
import { getFileData, getFile, removeFile} from "../api/api"
import { getDates } from "../api/get_dates"
import { SUCCESS } from "../api/status"
import { useRef, useEffect, useState, useLayoutEffect} from "react"
import { useNavigate, Link} from "react-router"
import Upload from "./upload"

const DATE_HEIGHT = 70
const EMPTY_SPACE_AT_THE_END_HEIGHT = 50
const LOAD_MARGIN = 1000
const FILE_DATA_LOAD_DELAY_MS = 100

class DayData
{
    constructor(start, height, date){
        this.start = start
        this.height = height
        this.date = date
    }
}

class FiltrationData
{
    constructor(year, month){
        this.year = year
        this.month = month
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
        elements.push(new DayData(start, height, date))
        lastDayEnd = end
    }
    return elements
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

function Day({day, tileSize, setFocusedFileData, resizeDay, navigate}){
    let [fileData, setFileData] = useState([])
    useEffect(
        ()=>{
            let abort = false
            setTimeout(() => {
                async function fun()
                {
                    if(!abort)
                    {
                        const result = await getFileData(day.date.date)
                        if(result.status != "SUCCESS"){
                            navigate("/error")
                            return
                        }

                        if(!abort)
                        {
                            const fd = result.fileData
                            if(fd.length != day.date.file_count){
                                resizeDay(day.date.date, fd.length,)
                                return
                            }
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
                    {day.date.date}
                </div>
                {tiles}
            </div>
}

function FocusedFile({focusedFileData, focusedFileUrl, setFocusedFileData, removePhoto}){
    let [showInfo, setShowInfo] = useState(false)
    let [showConfirmation, setShowConfirmation] = useState(false)
    if(focusedFileData === null || focusedFileUrl === null){
        if(showInfo)
        {
            setShowInfo(false)
        }
        if(showConfirmation){
            setShowConfirmation(false)
        }
        return null
    }

    function exit(){
        setFocusedFileData(null)
    }

    function info(){
        setShowInfo(!showInfo)
    }

    function del(){
        setShowConfirmation(true)
    }

    const date = focusedFileData.creation_date.substring(0, focusedFileData.creation_date.indexOf(" "))
    let removalConfirmationPopUp = showConfirmation ? <div className="pop_up_container removal_confirmation_container">
            <div className="pop_up_window">
                <h1>Are you sure?</h1>
                <div className="button" onClick={()=>{removePhoto(focusedFileData.id, date)}}>Yes</div>
                <div className="button" onClick={()=>{setShowConfirmation(false)}}>No</div>
            </div>
        </div>
        : null
    let description = showInfo ?
        <div className="description_container">
            <div className="description">
                <h1>{focusedFileData.filename}</h1>
                <h1>{focusedFileData.creation_date}</h1>
                <h1>{focusedFileData.mime_type}</h1>
            </div>
        </div>
        : null
    return <>
                {removalConfirmationPopUp}
                <div className="file">
                    <img src={focusedFileUrl}/>
                </div>
                <div className="exit button" onClick={exit}><img src="exit.png"/></div>
                <div className="info button" onClick={info}><img src="info.png"/></div>
                <div className="del button" onClick={del}><img src="trash.png"/></div>
                <a className="download button" href={focusedFileUrl} download={focusedFileData.filename}><img src="arrow.png"/></a>
                {description}
            </>
}

function FiltrationMenu({setFiltration, setPopUp, setDates}){
    return  <div className="pop_up_window">
                <h1>Set filter</h1>
                <input type="number" max="9999" min="0" placeholder="year" id="year"/>
                <input type="number" min="1" max="12" placeholder="month" id="month"/>
                <button className="button" onClick={()=>{
                    const yearInput = document.getElementById("year")
                    const monthInput = document.getElementById("month")
                    if(yearInput.checkValidity() && month.checkValidity()){
                        setDates(null)
                        setFiltration(new FiltrationData(yearInput.value, monthInput.value))
                        setPopUp("NONE")
                    }
                }}>Apply</button>
            </div>
}

function PopUpContainer({setDates, dates, setFiltration, focusedFileData, setFocusedFileData, focusedFileUrl, setPopUp, popUp}){
    function removePhoto(id, date){
        let fun = async () => {
            let datesCopy = structuredClone(dates)
            for(let i = 0 ; i < dates.length ; i++){
                if(datesCopy[i].date == date){
                    const result = await removeFile(id)
                    if(result != "SUCCESS"){
                        navigate("/error")
                        return
                    }
                    datesCopy[i].file_count -= 1
                    if(datesCopy[i].file_count == 0){
                        datesCopy.splice(i, 1)
                    }
                    setDates(datesCopy)
                    setFocusedFileData(null)
                    break
                }
            }
        }
        fun()
    }

    if(popUp == "NONE" && focusedFileData == null){
        return null;
    }

    return  <div className="pop_up_container">
                {popUp == "FILTRATION" ? <FiltrationMenu setFiltration={setFiltration} setPopUp={setPopUp} setDates={setDates}/> : null}
                {focusedFileData != null ? <FocusedFile focusedFileData={focusedFileData} setFocusedFileData={setFocusedFileData} focusedFileUrl={focusedFileUrl} removePhoto={removePhoto}/> : null}
                {popUp == "UPLOAD" ? <Upload exit={() => setPopUp("NONE")}/> : null}
            </div>
}

export default function Gallery(){
    let tileSize = useRef(null)
    let gallery = useRef(null)
    let content = useRef(null)
    let [filtration, setFiltration] = useState(new FiltrationData("", ""))
    let [containerWidth, setContainerWidth] = useState(null)
    let [dates, setDates] = useState(null)
    let [elements, setElements] = useState(null)
    let [scrollData, setScrollData] = useState(new ScrollData(0, window.innerHeight))
    let lastScrollData = useRef(new ScrollData(0, window.innerHeight))
    let [focusedFileData, setFocusedFileData] = useState(null)
    let [focusedFileUrl, setFocusedFileUrl] = useState(null)
    let [popUp, setPopUp] = useState("NONE")
    let navigate = useNavigate()

    useEffect(
        () => {
            let timeout = null
            const resizeObserver = new ResizeObserver((entries) => {
                    for (const entry of entries) {
                        if (entry.contentBoxSize) {
                            if(timeout !== null){
                                clearTimeout(timeout)
                            }
                            timeout = setTimeout(()=>{setContainerWidth(entry.contentBoxSize[0].inlineSize)}, 200)
                        }
                    }
               }
            )
            resizeObserver.observe(content.current)
            return () => {resizeObserver.disconnect()}
        },[])

        useEffect(
            ()=>{
                let abort = false
                async function fun(){
                    const result = await getDates(filtration)
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
                return () => {abort = true}
            }, [filtration]
        )

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
                setElements(null)
                return
            }
            const elements = createElements(dates, containerWidth, tileSize)
            setElements(elements)
        },[containerWidth, dates]
    )

    function resizeDay(date, newFileCount){
        let fun = async () => {
            let datesCopy = structuredClone(dates)
            for(let i = 0 ; i < dates.length ; i++){
                if(datesCopy[i].date == date){
                    datesCopy[i].file_count = newFileCount
                    if(datesCopy[i].file_count == 0){
                        datesCopy.splice(i, 1)
                    }
                    setDates(datesCopy)
                    break
                }
            }
        }
        fun()
    }

    let outlet = <h2>Loading...</h2>
    if(elements !== null)
    {
        outlet = []
        let totalHeight = 0
        for(let i = 0 ; i < elements.length ; i++)
        {
            const element = elements[i]
            totalHeight += element.height + 1
            if(element.start - LOAD_MARGIN <= scrollData.bottom && element.start + element.height + LOAD_MARGIN >= scrollData.top)
            {
                outlet.push(<Day key={element.date.date + "_" + element.date.file_count} day={element} tileSize={tileSize.current} setFocusedFileData={setFocusedFileData} resizeDay={resizeDay} navigate={navigate}/>)
            }
        }
        outlet.push(<div key={totalHeight} style={{height: `${totalHeight + EMPTY_SPACE_AT_THE_END_HEIGHT}px`}}></div>)
    }

    function scroll(e){
        let gallery = e.currentTarget
        const scrollTop = gallery.scrollTop
        const scrollBottom = scrollTop + gallery.offsetHeight
        const updateThreshold = LOAD_MARGIN / 2
        if(Math.abs(lastScrollData.current.top - scrollTop) > updateThreshold || Math.abs(lastScrollData.current.bottom - scrollBottom) > updateThreshold){
            lastScrollData.current = new ScrollData(scrollTop, scrollBottom)
            setScrollData(lastScrollData.current)
        }
    }

    return <div className="gallery_container">
                <PopUpContainer dates={dates} setFiltration={setFiltration} setDates={setDates} focusedFileData={focusedFileData} setFocusedFileData={setFocusedFileData} focusedFileUrl={focusedFileUrl} popUp={popUp} setPopUp={setPopUp}/>
                <header><div className="button" onClick={()=>{setPopUp("FILTRATION")}}>Filter</div><div className="button" onClick={()=>{setPopUp("UPLOAD")}}>Upload</div><Link className="button" to={"/login"}>Logout</Link></header>
                <div ref={gallery} className="gallery" onScroll={scroll}>
                    <div ref={content} className="content">
                        {outlet}
                    </div>
                </div>
            </div>
}
