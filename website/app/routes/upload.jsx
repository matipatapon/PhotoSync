import { useState, useEffect, useRef} from 'react'
import { useNavigate } from 'react-router'
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

export default function Upload({exit}){
    let files = useRef([])
    let [stage, setStage] = useState("SELECT")
    let [processedFileCount, setProcessedFileCount] = useState(0)
    let uploadedFilenames = useRef([])
    let alreadyExistedFilenames = useRef([])
    let unsupportedFilenames = useRef([])
    let failedFilenames = useRef([])

    if(stage !== "UPLOAD" && processedFileCount !== 0)
    {
        setProcessedFileCount(0)
    }

    const filenameToHtml = (filename) => {
        return [<div>&#x25CF; {filename}</div>, <br/>]
    }

    useEffect(
        () => {
            if(stage === "UPLOAD")
            {
                if(processedFileCount === files.current.length){
                    setStage("FINISH")
                    return
                }
                async function upload(){
                    const file = files.current[processedFileCount]
                    const filename = file.path + file.file.name
                    const status = await uploadPhoto(file.file)
                    if(status === "SUCCESS")
                    {
                        uploadedFilenames.current.push(filenameToHtml(filename))
                        setProcessedFileCount(processedFileCount + 1)
                    }
                    else if(status === "ALREADY_EXISTS"){
                        alreadyExistedFilenames.current.push(filenameToHtml(filename))
                        setProcessedFileCount(processedFileCount + 1)
                    }
                    else if(status === "UNSUPPORTED")
                    {
                        unsupportedFilenames.current.push(filenameToHtml(filename))
                        setProcessedFileCount(processedFileCount + 1)
                    }
                    else
                    {
                        failedFilenames.current.push(filenameToHtml(filename))
                        setProcessedFileCount(processedFileCount + 1)
                    }
                }
                upload()
            }
        }, [stage, processedFileCount])

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
                <h1>Upload your files</h1>
                <label className="button" htmlFor="file_upload">Select</label>
                <input id="file_upload" type='file' multiple={true} onChange={select}/>
                <div className='button' onClick={exit}>Cancel</div>
        </>
    }
    if(stage === "LOAD"){
        outlet = <h1>Please wait...</h1>
    }
    if(stage === "OPTIONS"){
        outlet = <>
            <h1>{files.current.length} files selected</h1>
            <div className='button' onClick={() => {setStage("UPLOAD") ; uploadedFilenames.current = []; unsupportedFilenames.current = []; alreadyExistedFilenames.current = [];}}>Upload</div>
            <div className='button' onClick={() => setStage("SELECT")}>Clear</div>
        </>
    }
    if(stage === "FINISH" || stage === "UPLOAD"){
        outlet = <>
            {stage === "UPLOAD" ? <h1>Processed {processedFileCount}/{files.current.length}</h1> : <h1>Finished</h1>}
            <div className='filelist_container' style={{display: uploadedFilenames.current.length != 0 ? "block" : "none"}}>
                <h1>{uploadedFilenames.current.length} files were uploaded</h1>
                <div className='filelist'>
                    {uploadedFilenames.current}
                </div>
            </div>
            <div className='filelist_container' style={{display: alreadyExistedFilenames.current.length != 0 ? "block" : "none"}}>
                <h1>{alreadyExistedFilenames.current.length} files are already uploaded</h1>
                <div className='filelist'>
                    {alreadyExistedFilenames.current}
                </div>
            </div>
            <div className='filelist_container error' style={{display: failedFilenames.current.length != 0 ? "block" : "none"}}>
                <h1>{failedFilenames.current.length} files failed to upload</h1>
                <div className='filelist'>
                    {failedFilenames.current}
                </div>
            </div>
            <div className='filelist_container error' style={{display: unsupportedFilenames.current.length != 0 ? "block" : "none"}}>
                <h1>{unsupportedFilenames.current.length} files are unsupported</h1>
                <div className='filelist'>
                    {unsupportedFilenames.current}
                </div>
            </div>
            {stage === "UPLOAD"
                ? <div className='button' onClick={() => setStage("SELECT")}>Cancel</div>
                : <div className='button' onClick={exit}>Ok</div>}
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

    return <div className='pop_up_window' onDragOver={dragOver} onDrop={drop}>
                {outlet}
            </div>
}
