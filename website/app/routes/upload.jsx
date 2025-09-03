import { useState, useEffect, useRef} from 'react'
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
    let errorMsg = useRef("")
    let nameOfLastUploadedFile = useRef(null)

    let [stage, setStage] = useState("SELECT")
    let [uploadedFileCount, setUploadedFileCount] = useState(0)

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
                        setUploadedFileCount(uploadedFileCount + 1)
                    }
                    else
                    {
                        if(status === "NOT_LOGGED_IN")
                        {
                            errorMsg.current = "You need to login to upload files"
                        }
                        else if(status === "UNSUPPORTED")
                        {
                            errorMsg.current = `${file.path}${file.file.name} has unsupported type`
                        }
                        else if(status === "TOKEN_EXPIRED")
                        {
                            errorMsg.current = `Your session has expired`
                        }
                        else
                        {
                            errorMsg.current = "Error"
                        }
                        setStage("ERROR")
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
                <h2>Drop or select files</h2>
                <label htmlFor="file_upload">Select</label>
                <input id="file_upload" type='file' multiple={true} onChange={select}/>
        </>
    }
    if(stage === "LOAD"){
        outlet = <h2>Please wait...</h2>
    }
    if(stage === "OPTIONS"){
        outlet = <>
                <h2>{files.current.length} files selected</h2>
                <div id='upload' className='button' onClick={() => setStage("UPLOAD")}>Upload</div>
                <div id='clear' className='button' onClick={() => setStage("SELECT")}>Clear</div>
        </>
    }
    if(stage === "UPLOAD"){
        outlet = <>
                <h2>Uploaded {uploadedFileCount}/{files.current.length}</h2>
                <h3>{nameOfLastUploadedFile.current}</h3>
        </>
    }
    if(stage === "FINISH"){
        outlet = <>
        <h2>All files uploaded</h2>
        <div className='button' onClick={() => setStage("SELECT")}>Ok</div>
        </>
    }
    if(stage === "ERROR"){
        outlet = <>
            <h2>{errorMsg.current}</h2>
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

    return  <div className='upload_container'>
                <div className='upload' onDragOver={dragOver} onDrop={drop}>
                    {outlet}
                </div>
            </div>
}
