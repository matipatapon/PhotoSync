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
    let [stage, setStage] = useState("SELECT")
    let [uploadedCount, setUploadedCount] = useState(0)

    useEffect(
        () => {
            async function upload() {
                if(stage !== "UPLOAD"){
                    return
                }
                const status = await uploadPhoto(files.current[uploadedCount].file)
                if(status === "SUCCESS"){
                    if(uploadedCount + 1 < files.current.length){
                        setUploadedCount(uploadedCount + 1)
                        
                    } else{
                        setStage("FINISH")
                    }
                } else {
                    setUploadedCount(uploadedCount + 1)
                }
            }
            upload()
        }
        ,
        [stage, uploadedCount])

    function dragOver(event) {
        event.preventDefault()
        
    }

    function drop(event){
        event.preventDefault()
        if(stage === "SELECT"){
            getFilesFromItems(event.dataTransfer.items, files.current, () => setStage("OPTIONS"))
            setStage("LOAD")
        }
    }

    let outlet
    if(stage === "SELECT"){
        outlet = <>
                <h2>Drop or select files</h2>
                <label htmlFor="file_upload">Select</label>
                <input id="file_upload" type='file' multiple={true} onChange={(event) => {files.current = event.currentTarget.files ; setStage("OPTIONS")}}/>
        </>
    }
    if(stage === "LOAD"){
        outlet = <h2>Please wait...</h2>
    }
    if(stage === "OPTIONS"){
        outlet = <>
                <h2>{files.current.length} files selected</h2>
                <div id='upload' className='button' onClick={() => {uploadedCount = 0 ; setStage("UPLOAD")}}>Upload</div>
                <div id='clear' className='button' onClick={() => {files.current=[] ; setStage("SELECT")}}>Clear</div>
        </>
    }
    if(stage === "UPLOAD"){
        outlet = <>
                    <div className='upload_container'>
                    <h2>Uploaded {uploadedCount}/{files.current.length}</h2>
                    </div>
        </>
    }
    if(stage === "FINISH"){
        outlet = <h3>finished</h3>
    }
    if(stage === "ERROR"){
        outlet = <h3>Something went wrong</h3>
    }
        return <div className='upload_container' onDragOver={dragOver} onDrop={drop}>
            {outlet}
        </div>
}
