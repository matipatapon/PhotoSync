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
    let uploadedFiles = useRef([])
    let alreadyExistedFiles = useRef([])
    let unsupportedFiles = useRef([])
    let navigate = useNavigate()

    if(stage !== "UPLOAD" && processedFileCount !== 0)
    {
        setProcessedFileCount(0)
    }

    useEffect(
        () => {
            if(stage === "UPLOAD")
            {
                if(processedFileCount == 0){
                    uploadedFiles.current = [];
                    unsupportedFiles.current = [];
                    alreadyExistedFiles.current = [];
                }
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
                        uploadedFiles.current.push(filename)
                        setProcessedFileCount(processedFileCount + 1)
                    }
                    else if(status === "ALREADY_EXISTS"){
                        alreadyExistedFiles.current.push(filename)
                        setProcessedFileCount(processedFileCount + 1)
                    }
                    else if(status === "UNSUPPORTED")
                    {
                        unsupportedFiles.current.push(filename)
                        setProcessedFileCount(processedFileCount + 1)
                    }
                    else
                    {
                        navigate("/error")
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
            <div className='button' onClick={() => setStage("UPLOAD")}>Upload</div>
            <div className='button' onClick={() => setStage("SELECT")}>Clear</div>
        </>
    }
    if(stage === "UPLOAD"){
        outlet = <>
            <h1>Processed {processedFileCount}/{files.current.length}</h1>
            <div className='button' onClick={() => setStage("SELECT")}>Cancel</div>
        </>
    }
    if(stage === "FINISH"){
        let uploadedFilenames = []
        uploadedFiles.current.forEach(fn => {
            uploadedFilenames.push([<div>&#x25CF; {fn}</div>, <br/>])
        });

        let unsupportedFilenames = []
        unsupportedFiles.current.forEach(fn => {
            unsupportedFilenames.push([<div>&#x25CF; {fn}</div>, <br/>])
        });

        let alreadyExistsFilenames = []
        alreadyExistedFiles.current.forEach(fn => {
            alreadyExistsFilenames.push([<div>&#x25CF; {fn}</div>, <br/>])
        });

        outlet = <>
            <div className='filelist_container' style={{display: uploadedFiles.current.length != 0 ? "block" : "none"}}>
                <h1>{uploadedFiles.current.length} files uploaded</h1>
                <div className='filelist'>
                    {uploadedFilenames}
                </div>
            </div>
            <div className='filelist_container' style={{display: unsupportedFiles.current.length != 0 ? "block" : "none"}}>
                <h1>{unsupportedFiles.current.length} files unsupported</h1>
                <div className='filelist'>
                    {unsupportedFilenames}
                </div>
            </div>
            <div className='filelist_container' style={{display: alreadyExistedFiles.current.length != 0 ? "block" : "none"}}>
                <h1>{alreadyExistedFiles.current.length} files already uploaded</h1>
                <div className='filelist'>
                    {alreadyExistsFilenames}
                </div>
            </div>
            <div className='button' onClick={() => setStage("SELECT")}>Ok</div>
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
