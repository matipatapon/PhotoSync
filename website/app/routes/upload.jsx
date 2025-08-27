import { useState, useEffect, useRef} from 'react'
import { uploadPhoto } from '../api/api'

import './upload.css'

function Submit(setUploadState){
    return function(){
        setUploadState(true)
    }
}

function File({file, uploading, startNext}){
    const [uploadStatus, setUploadStatus] = useState("NOT_STARTED")
    useEffect(()=>{
        let ignore = false
        if(uploadStatus !== "NOT_STARTED" || uploading === false){
            return
        }
        setUploadStatus("IN_PROGRESS")
        async function uploadFile(){
            const status = await uploadPhoto(file)
            if(ignore === true){
                return
            }
            setUploadStatus(status)
            startNext()
        }
        uploadFile()

        return () =>{
            ignore = true
        }

    }, [uploading])
    return <div className="file">{file.name} | {uploadStatus}</div>
}

function Files({files, uploading}){
    const [next, setNext] = useState(0)
    function startNext(){
        setNext(next + 1)
    }

    let rows = []
    for(let id = 0 ; id < files.length ; id++){
        let shouldGivenElementUpload = uploading && next === id
        rows.push(<File key={id} file={files[id]} uploading={shouldGivenElementUpload} startNext={startNext}/>)
    }
    return rows
}

export default function Upload(){
    const [areFilesDraggedOver, setAreFilesDraggedOver] = useState(false)
    const [files, setFiles] = useState([])
    const [uploading, setUploading] = useState(false)
    const [filesKey, setFilesKey] = useState(0)
    let dragCounter = useRef(0)

    function drop(event){
        if(event.dataTransfer.files.length == 0){
            return
        }
        event.preventDefault()
        setFilesKey(filesKey + 1)
        setUploading(false)
        setFiles(event.dataTransfer.files)
    }

    function dragEnter(event){
        dragCounter.current += 1
        if(dragCounter.current === 1){
            setAreFilesDraggedOver(true)
        }
    }

    function dragLeave(event){
        dragCounter.current -= 1
        if(dragCounter.current === 0){
            setAreFilesDraggedOver(false)
        }
    }

    function dragOver(event) {
        event.preventDefault()
    }

    const classNames = `upload ${areFilesDraggedOver ? "upload_files_over": ""}`
    return <div className={classNames} onDrop={drop} onDragOver={dragOver} onDragEnter={dragEnter} onDragLeave={dragLeave}>
                <Files key={filesKey} files={files} uploading={uploading}/>
                <div className="submit" onClick={Submit(setUploading)}>Submit</div>
            </div>
}
