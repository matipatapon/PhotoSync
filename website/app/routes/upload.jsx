import { useState, useEffect, useRef} from 'react'
import { Link, useNavigate } from 'react-router'
import { uploadPhoto } from '../api/api'
import './upload.css'

function File(file, path){
    this.file = file
    this.path = path
}

function getFilesFromItems(items, output, callback){
    let pendingFileCount = 0
    let pendingDirs = 0

    let setFilesIfAllFilesGathered = () => {
        if(pendingDirs == 0 && pendingFileCount == 0){
            callback()
        }
    }

    let getFilesFromItem = (item, path) => { 
        if(item.isFile){
            pendingFileCount += 1
            item.file((f) => {
                output.push(new File(f, path))
                pendingFileCount -= 1
                setFilesIfAllFilesGathered();
            })
        }
        if(item.isDirectory){
            pendingDirs += 1
            item.createReader().readEntries(function(entries) {
                for (const entry of entries){
                    getFilesFromItem(entry, path + item.name + "/")
                }
                pendingDirs -= 1
                setFilesIfAllFilesGathered()
            })
        }
    }

    for(const item of items){
        const entry = item.webkitGetAsEntry()
        if(entry !== null){
            getFilesFromItem(entry, "")
        }
    }

    setFilesIfAllFilesGathered()
}


export default function Upload(){
    let files = useRef([])
    let nameOfLastUploadedFile = useRef(null)
    let [stage, setStage] = useState("SELECT")
    let [uploadedFileCount, setUploadedFileCount] = useState(0)
    let navigate = useNavigate()

    if(stage !== "UPLOAD" && uploadedFileCount !== 0)
    {
        setUploadedFileCount(0)
    }

    useEffect(
        () => {
            if(stage === "UPLOAD")
            {
                async function upload(){
                    const file = files.current[uploadedFileCount]
                    const status = await uploadPhoto(file.file)
                    if(status === "SUCCESS" || status === "ALREADY_EXISTS")
                    {
                        nameOfLastUploadedFile.current = file.path + file.file.name
                        if(uploadedFileCount + 1 === files.current.length)
                        {
                            setStage("FINISH")
                        }
                        else
                        {
                            setUploadedFileCount(uploadedFileCount + 1)
                        }
                    }
                    else
                    {
                        navigate("/error")
                    }
                }
                upload()
            }
        }, [stage, uploadedFileCount])

    function select(event){
        files.current = []
        for(const file of event.currentTarget.files)
        {
            files.current.push(new File(file, ""))
        }
        setStage("OPTIONS")
    }

    let outlet
    if(stage === "SELECT"){
        outlet = <>
                <h1>Select files</h1>
                <div className='buttons'>
                    <label className="button" htmlFor="file_upload">Select</label>
                    <input id="file_upload" type='file' multiple={true} onChange={select}/>
                </div>
        </>
    }
    if(stage === "LOAD"){
        outlet = <h1>Please wait...</h1>
    }
    if(stage === "OPTIONS"){
        outlet = <>
                <h1>{files.current.length} files selected</h1>
                <div className='buttons'>
                    <div className='button' onClick={() => setStage("UPLOAD")}>Upload</div>
                    <div className='button' onClick={() => setStage("SELECT")}>Clear</div>
                </div>
        </>
    }
    if(stage === "UPLOAD"){
        outlet = <>
                <h1>Uploading {uploadedFileCount}/{files.current.length}</h1>
                <div className='buttons'>
                    <div className='button' onClick={() => setStage("SELECT")}>Cancel</div>
                </div>
        </>
    }
    if(stage === "FINISH"){
        outlet = <>
        <h1>All files uploaded</h1>
        <div className='buttons'>
            <div className='button' onClick={() => setStage("SELECT")}>Ok</div>
        </div>
        </>
    }

    function dragOver(event) {
        event.preventDefault()
    }

    function drop(event){
        event.preventDefault()
        if(stage === "SELECT"){
            files.current = []
            getFilesFromItems(event.dataTransfer.items, files.current, () => setStage("OPTIONS"))
            setStage("LOAD")
        }
    }

    return  <>
                <header><Link className="button" to={"/gallery"}>Gallery</Link><Link className="button" to={"/login"}>Logout</Link></header>
                <div className='window_container'>
                    <div className='window' onDragOver={dragOver} onDrop={drop}>
                        {outlet}
                    </div>
                </div>
            </>
}
