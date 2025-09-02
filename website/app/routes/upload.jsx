import { useState, useEffect, useRef} from 'react'
import { uploadPhoto } from '../api/api'
import './upload.css'

function File(file, path){
    this.file = file
    this.path = path
}

function getFilesFromItems(items, setFiles){
    let files = []
    let pendingFileCount = 0
    let pendingDirs = 0

    let setFilesIfAllFilesGathered = () => {
        if(pendingDirs == 0 && pendingFileCount == 0){
            setFiles(files)
        }
    }

    let getFilesFromItem = (item, path) => { 
        if(item.isFile){
            pendingFileCount += 1
            item.file((f) => {
                files.push(new File(f, path))
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

function FileSelector({setFiles, setLoading}){
    function drop(event){
        event.preventDefault()
        setLoading(true)
        getFilesFromItems(event.dataTransfer.items, setFiles)
    }

    function dragOver(event) {
        event.preventDefault()
    }

    function filesSelected(event){
        setFiles(event.currentTarget.files)
    }

    return  <div className='upload_container' onDrop={drop} onDragOver={dragOver}>
                <h2>Drop or select files</h2>
                <label htmlFor="file_upload">Select</label>
                <input id="file_upload" type='file' multiple={true} onChange={filesSelected}/>
            </div>
}

function FileUploader({setFiles, files}){
    return <div className='upload_container'>
        <h2>{files.length} files selected</h2>
        <div id='upload' className='button'>Upload</div>
        <div id='clear' className='button' onClick={() => setFiles([])}>Clear</div>
    </div>
}

function FileLoader(){
    return <div className='upload_container'>
        <h2>Please wait...</h2>
    </div>
}

export default function Upload(){
    let [files, setFiles] = useState([])
    let [loading, setLoading] = useState(false)
    let setFilesWrapper = (files) => {
        setFiles(files)
        setLoading(false)
    }
    if(loading){
        return <FileLoader/>
    }
    if(files.length == 0){
        return <FileSelector setFiles={setFilesWrapper} setLoading={setLoading}/>
    }
    return <FileUploader setFiles={setFiles} files={files}/>
}
