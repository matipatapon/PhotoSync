import { useState, useEffect, useRef} from 'react'
import { useNavigate } from 'react-router'
import { uploadPhoto } from '../api/api'
import './upload.css'

function File(file, path){
    this.file = file
    this.path = path
}

function UploadedFile(filename, creationDate){
    this.filename = filename
    this.creationDate = creationDate
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
    let failedFiles = useRef([])

    if(stage !== "UPLOAD" && processedFileCount !== 0)
    {
        setProcessedFileCount(0)
    }

    const addFailedFile = (filename) => {
        failedFiles.current.push([<div key={failedFiles.current.length + "_" + 1}>&#x25CF; {filename}</div>, <br key={failedFiles.current.length + "_" + 2}/>])
    }

    const generateListForUploadedFiles = () => {
        let dateToFilenames = new Map()
        uploadedFiles.current.forEach(uploadedFile => {
            const creationDay = uploadedFile.creationDate.split([" "])[0]
            if(dateToFilenames.has(creationDay)){
                dateToFilenames.get(creationDay).push(uploadedFile.filename)
            } else{
                dateToFilenames.set(creationDay, [uploadedFile.filename])
            }
        });

        let html = []
        const dateToFilenamesSorted = new Map([...dateToFilenames.entries()].sort((a, b) => {
            if(a == b) return 0;
            else if(a > b) return -1;
            else return 1;
        }))
        dateToFilenamesSorted.forEach((value, key) => {
            console.log(value)
            html.push(<h2>{key}</h2>)
            value.forEach((filename) => {html.push(<div>&#x25CF; {filename}</div>) ; html.push(<br/>)})
        })
        return html
    }

    const generateSummaryForFiles = (filenames, postfix, count) => {
        return <div className={'filelist_container'} style={{display: filenames.length != 0 ? "block" : "none"}}>
                <h1>{count} {postfix}</h1>
                <div className='filelist'>
                    {filenames}
                </div>
            </div>
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
                    const result = await uploadPhoto(file.file)
                    if(result.status === "SUCCESS" || result.status === "ALREADY_EXISTS")
                    {
                        uploadedFiles.current.push(new UploadedFile(filename, result.creationDate))
                    }
                    else
                    {
                        addFailedFile(filename)
                    }
                    setProcessedFileCount(processedFileCount + 1)
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
            <div className='button' onClick={() => {setStage("UPLOAD") ; uploadedFiles.current = []; failedFiles.current = [];}}>Upload</div>
            <div className='button' onClick={() => setStage("SELECT")}>Clear</div>
        </>
    }
    if(stage === "UPLOAD"){
        outlet = <>
            <h1>Processing {processedFileCount}/{files.current.length}</h1>
            {generateSummaryForFiles(failedFiles.current, "files failed to be uploaded", failedFiles.current.length)}
            <div className='button' onClick={() => setStage("FINISH")}>Cancel</div>
        </>
    }
    if(stage === "FINISH"){
        outlet = <>
            <h1>Finished</h1>
            {generateSummaryForFiles(generateListForUploadedFiles(), "files were uploaded", uploadedFiles.current.length)}
            {generateSummaryForFiles(failedFiles.current, "files failed to be uploaded", failedFiles.current.length)}
            <div className='button' onClick={() => {window.location.reload(), exit()}}>Ok</div>
        </>
    }

    function dragOver(event) {
        event.preventDefault()
    }

    function drop(event){
        event.preventDefault()
        if(stage === "SELECT"){
            files.current = []
            setStage("LOAD")
            getFilesFromItems(event.dataTransfer.items, files.current, () => setStage("OPTIONS"))
        }
    }

    return <div className='pop_up_window' onDragOver={dragOver} onDrop={drop}>
                {outlet}
            </div>
}
