import { useState, useEffect, useRef} from 'react'
import { uploadPhoto } from '../api/api'

import './upload.css'

function File(file, path){
    this.file = file
    this.path = path
}

function gatherFiles(items, setFiles){
    let files = []
    let filesCanBeSetted = false
    let pendingFileCount = 0
    let pendingDirs = 0

    let setFilesIfAllFilesGathered = () => {
        if(filesCanBeSetted && pendingFileCount == 0){
            filesCanBeSetted = false
            setFiles(files)
        }
    }

    function gatherFilesImpl(item, path){
        if(item.isFile){
            pendingFileCount += 1
            item.file((f) => {
                files.push(new File(f, path)) ; pendingFileCount -= 1 ; setFilesIfAllFilesGathered() ; console.log(f.name)
            })
            return
        }
        if (item.isDirectory){
            item.createReader().readEntries(function(entries) {
                for (const entry of entries){
                    pendingDirs += 1
                    gatherFilesImpl(entry, path + item.name + "/")
                    pendingDirs -= 1
                }
            })
        }
        setFilesIfAllFilesGathered()
    }

    for(const item of items){
        const entry = item.webkitGetAsEntry()
        if(entry === null){
            continue
        }
        gatherFilesImpl(entry, "/")
    }
}

function Submit(setUploadState){
    return function(){
        setUploadState(true)
    }
}

// function File({file, uploading, startNext}){
//     const [uploadStatus, setUploadStatus] = useState("NOT_STARTED")
//     useEffect(()=>{
//         let ignore = false
//         if(uploadStatus !== "NOT_STARTED" || uploading === false){
//             return
//         }
//         setUploadStatus("IN_PROGRESS")
//         async function uploadFile(){
//             const status = await uploadPhoto(file)
//             if(ignore === true){
//                 return
//             }
//             setUploadStatus(status)
//             startNext()
//         }
//         uploadFile()

//         return () =>{
//             ignore = true
//         }

//     }, [uploading])
//     return <div className="file">{file.name} | {uploadStatus}</div>
// }

// function Files({files, uploading}){
//     const [next, setNext] = useState(0)
//     function startNext(){
//         setNext(next + 1)
//     }

//     let rows = []
//     for(let id = 0 ; id < files.length ; id++){
//         let shouldGivenElementUpload = uploading && next === id
//         rows.push(<File key={id} file={files[id]} uploading={shouldGivenElementUpload} startNext={startNext}/>)
//     }
//     return rows
// }

function Files({files}){
    let result = []
    for(let i = 0 ; i < files.length ; i++){
        result.push(<li key={i}>{files[i].path + files[i].file.name}</li>)
    }
    return result
}

export default function Upload(){
    const [areFilesDraggedOver, setAreFilesDraggedOver] = useState(false)
    const [uploading, setUploading] = useState(false)
    const [filesKey, setFilesKey] = useState(0)
    let dragCounter = useRef(0)
    let [files, setFiles] = useState([])

    function drop(event){
        event.preventDefault()
        gatherFiles(event.dataTransfer.items, setFiles)

        // 
        // setFilesKey(filesKey + 1)
        // setUploading(false)
        // setFiles(event.dataTransfer.files)
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
    return <div className='upload_container'>
        <div className='upload' onDrop={drop} onDragOver={dragOver}>
            <input type='file' multiple={true}/>
            <Files files={files}/>
        </div>
        
    </div>
    return <div className={classNames} onDrop={drop} onDragOver={dragOver} onDragEnter={dragEnter} onDragLeave={dragLeave}>
                <Files key={filesKey} files={files} uploading={uploading}/>
                <div className="submit" onClick={Submit(setUploading)}>Submit</div>
            </div>
}
