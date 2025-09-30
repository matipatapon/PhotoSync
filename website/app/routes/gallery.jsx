import "./gallery.css"
import { useRef, useEffect, useState, useLayoutEffect} from "react"

function log(string){
    console.log(`[Gallery]: '${string}'`)
}

export default function Gallery(){
    let gallery = useRef(null)
    let [galleryWidth, setGalleryWidth] = useState(null)
    const resizeObserver = new ResizeObserver((entries) => {
        for (const entry of entries) {
            if (entry.contentBoxSize) {
                log(`galleryWidth => ${entry.contentBoxSize[0].inlineSize}`)
                setGalleryWidth(entry.contentBoxSize[0].inlineSize)
            }
        }
    })

    useEffect(
        () => {
            resizeObserver.observe(gallery.current)
        },[])
    
    useLayoutEffect(
        () => {
            if(galleryWidth === null){
                return
            }
        },[galleryWidth]
    )

    return <div ref={gallery} className="gallery">
        {galleryWidth}
          </div>
}
