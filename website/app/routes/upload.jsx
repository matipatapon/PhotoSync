import { useState, useEffect, useRef} from 'react'

import './upload.css'

function Submit(setUploadState){
    return function(){
        setUploadState(true)
    }
}

function DragOver(){
    return function (event){
        event.preventDefault()
    }
}

function DragOverEnter(setDragOver){
    return function(event){
        setDragOver(true)
    }
}

function DragOverLeave(setDragOver){
    return function(event){
        setDragOver(false)
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
            await new Promise(r => setTimeout(r, 5000));
            setUploadStatus("FINISHED")
            startNext()
        }
        uploadFile()

        return () =>{
            ignore = true
        }

    }, [uploading])
    return <div>{file.name} | {uploadStatus}</div>
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
    const [dragOver, setDragOver] = useState(false)
    const [files, setFiles] = useState([])
    const [uploading, setUploading] = useState(false)
    const [filesKey, setFilesKey] = useState(0)

    function drop(event){
        event.preventDefault()
        setFilesKey(filesKey + 1)
        setUploading(false)
        setFiles(event.dataTransfer.files)
    }

    return <div className={"upload" + " " + (dragOver ? "upload_over": "upload_not_over")} onDrop={drop} onDragOver={DragOver()} onDragEnter={DragOverEnter(setDragOver)} onDragLeave={DragOverLeave(setDragOver)}>
        <ol>
            <Files key={filesKey} files={files} uploading={uploading}/>
            <button onClick={Submit(setUploading)}>Submit</button>
        </ol>
    </div>
}
